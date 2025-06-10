package main

import (
	"context"
	"dojo_bot/model"
	"dojo_bot/storage"
	"log"
	"os"
	"time"

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

	// обработка сообщений
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// сохраняем пользователя в Mongo
		user := model.User{
			Username: update.Message.From.UserName,
			ChatID:   update.Message.Chat.ID,
			JoinedAt: time.Now().Format(time.RFC3339),
		}

		if err := storage.SaveUser(db, user); err != nil {
			log.Println("Ошибка сохранения:", err)
		}

		// отправляем ответ
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую тебя странник, а твой ник: "+update.Message.From.UserName+"!")
		bot.Send(msg)
	}
}
