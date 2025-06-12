package handler

import (
	"dojo_bot/internal/svc"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserHandlerInterface interface {
	HandleUpdate(update tgbotapi.Update) error
}

type UserHandler struct {
	userSvc svc.UserSvcInterface
	bot     *tgbotapi.BotAPI
}

func NewUserHandler(userSvc *svc.UserService, bot *tgbotapi.BotAPI) UserHandlerInterface {
	return &UserHandler{
		userSvc: userSvc,
		bot:     bot,
	}
}

// добавить контекст
// обработка сообщения
func (h *UserHandler) HandleUpdate(update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	if err := h.userSvc.ProcessUser(update.Message.From.UserName, update.Message.Chat.ID); err != nil {
		log.Printf("Ошибка обработки пользователя: %v", err)
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, "+update.Message.From.UserName+"!")
	_, err := h.bot.Send(msg)
	return err
}
