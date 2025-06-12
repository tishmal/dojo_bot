package main

import (
	"context"
	"dojo_bot/internal/handler"
	"dojo_bot/internal/storage/db"
	"dojo_bot/internal/storage/repository"
	"dojo_bot/internal/svc"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Бот запущен: %s", bot.Self.UserName)

	// Подключение к MongoDB
	mongoClient, err := db.ConnectMongoDB(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Panic("Ошибка MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	database := mongoClient.Database("telegram_bot")

	repo := repository.NewUserRepo(database)
	userSvc := svc.NewUserSvc(repo)
	userHandler := handler.NewUserHandler(userSvc, bot)

	// timeout между запросами
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Канал для обновлений и graceful shutdown
	updates := bot.GetUpdatesChan(u)
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// Канал для обработки сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Worker pool
	workerCount := 5
	jobs := make(chan tgbotapi.Update, 100)

	// Запуск воркеров
	for i := 1; i <= workerCount; i++ {
		go func(id int) {
			for update := range jobs {
				if err := userHandler.HandleUpdate(update); err != nil {
					log.Printf("Worker %d: ошибка обработки: %v", id, err)
				}
			}
		}(i)
	}

	// Главный цикл обработки
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}
			jobs <- update
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("Получен сигнал завершения...")
	//cancel()
	close(jobs)
	log.Println("Бот остановлен")
}
