package main

import (
	"context"
	"dojo_bot/internal/storage"
	"dojo_bot/internal/worker"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// инициализация dojo бота
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Бот запущен: %s", bot.Self.UserName)

	// подключение к Mongo, до запуска бота
	mongoClient, err := storage.ConnectMongoDB("mongodb://localhost:27017")
	if err != nil {
		log.Panic("Ошибка MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background()) // закрываем соединение при выходе
	db := mongoClient.Database("telegram_bot")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// запускаем воркеры
	jobs := make(chan tgbotapi.Update, 100)
	for w := 1; w <= 5; w++ {
		go worker.Worker(w, jobs, bot, db)
	}

	// обработка сообщений
	for update := range updates {
		if update.Message == nil {
			continue
		}

		jobs <- update
	}
}
