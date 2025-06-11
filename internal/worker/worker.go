package worker

import (
	"dojo_bot/internal/storage"
	"dojo_bot/model"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func Worker(id int, jobs <-chan tgbotapi.Update, bot *tgbotapi.BotAPI, db *mongo.Database) {
	// Обработка сообщения
	for update := range jobs {
		// сохраняем пользователя в Mongo
		user := model.User{
			Username: update.Message.From.UserName,
			ChatID:   update.Message.Chat.ID,
			JoinedAt: time.Now().Format(time.RFC3339),
		}

		if err := storage.SaveUser(db, user); err != nil {
			log.Println("Ошибка сохранения:", err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Салам! Это сообщение обработано воркером "+strconv.Itoa(id))
		bot.Send(msg)

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую тебя странник, а твой ник: "+update.Message.From.UserName+"!")
		// bot.Send(msg)
	}
}
