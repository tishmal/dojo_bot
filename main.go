package main

import (
	"context"
	"dojo_bot/internal/handler"
	"dojo_bot/internal/storage/db"
	"dojo_bot/internal/storage/repository"
	"dojo_bot/internal/svc"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
)

func main() {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}

	// Инициализация бота
	botToken := os.Getenv("TELEGRAM_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	// Контекст с graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	// Создаем Gin роутер
	router := gin.Default()

	// Инициализация бота
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatalf("Bot creation error: %v", err)
	}

	// Подключение к MongoDB
	mongoClient, err := db.ConnectMongoDB(ctx, os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	// Инициализация сервисов
	userRepo := repository.NewUserRepo(mongoClient.Database("telegram_bot"))
	userSvc := svc.NewUserSvc(userRepo)
	userHandler := handler.NewUserHandler(ctx, userSvc, bot)

	// Настройка вебхука
	if webhookURL := os.Getenv("WEBHOOK_URL"); webhookURL != "" {
		setupWebhook(ctx, bot, webhookURL, router, userHandler)
	} else {
		setupLongPolling(ctx, bot, userHandler)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Запуск HTTP сервера
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}

func setupWebhook(ctx context.Context, bot *telego.Bot, url string, router *gin.Engine, handler handler.UserHandlerInterface) {
	// Удаляем предыдущий вебхук
	if err := bot.DeleteWebhook(ctx, &telego.DeleteWebhookParams{}); err != nil {
		log.Printf("Error deleting webhook: %v", err)
	}

	// Настраиваем новый вебхук
	if err := bot.SetWebhook(ctx, &telego.SetWebhookParams{
		URL:            url,
		MaxConnections: 50,
	}); err != nil {
		log.Fatalf("Webhook setup error: %v", err)
	}

	// Получаем канал обновлений
	updates, err := bot.UpdatesViaWebhook(
		ctx,
		func(wh telego.WebhookHandler) error {
			router.POST("/webhook", func(c *gin.Context) {
				// Читаем тело запроса
				body, err := io.ReadAll(c.Request.Body)
				if err != nil {
					log.Printf("Error reading request body: %v", err)
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}

				// Вызываем обработчик telego
				if err := wh(ctx, body); err != nil {
					log.Printf("Webhook handler error: %v", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				c.Status(http.StatusOK)
			})
			return nil
		},
	)
	if err != nil {
		log.Fatalf("Failed to setup webhook updates: %v", err)
	}

	// Обработка обновлений
	go processUpdates(ctx, updates, handler)
}

func setupLongPolling(ctx context.Context, bot *telego.Bot, handler handler.UserHandlerInterface) {
	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

	go processUpdates(ctx, updates, handler)
}

func processUpdates(ctx context.Context, updates <-chan telego.Update, handler handler.UserHandlerInterface) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			if err := handler.HandleUpdate(update); err != nil {
				log.Printf("Error handling update: %v", err)
			}
		}
	}
}
