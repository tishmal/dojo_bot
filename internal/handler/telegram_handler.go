package handler

import (
	"context"
	"dojo_bot/internal/svc"
	"log"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

type UserHandlerInterface interface {
	SetupMenuButton(ctx context.Context) error
	HandleUpdate(update telego.Update) error
}

type userHandler struct {
	ctx     context.Context
	userSvc svc.UserSvcInterface
	bot     *telego.Bot
	appURL  string
	chanURL string
}

func NewUserHandler(ctx context.Context, userSvc svc.UserSvcInterface, bot *telego.Bot, appURL, chanURL string) UserHandlerInterface {
	handler := &userHandler{
		ctx:     ctx,
		userSvc: userSvc,
		bot:     bot,
		appURL:  appURL,
		chanURL: chanURL,
	}

	// Устанавливаем menu button при создании хендлера
	if err := handler.SetupMenuButton(context.Background()); err != nil {
		log.Printf("Ошибка установки menu button: %v", err)
	}

	return handler
}

// Устанавливаем Web App кнопку в меню бота
func (h *userHandler) SetupMenuButton(ctx context.Context) error {
	menuButton := &telego.MenuButtonWebApp{
		Type: telego.ButtonTypeWebApp,
		Text: "Play",
		WebApp: telego.WebAppInfo{
			URL: h.appURL,
		},
	}

	err := h.bot.SetChatMenuButton(ctx, &telego.SetChatMenuButtonParams{
		MenuButton: menuButton,
	})
	if err != nil {
		return err
	}

	log.Println("Menu button успешно установлена")
	return nil
}

// Обработка апдейтов
func (h *userHandler) HandleUpdate(update telego.Update) error {
	if update.Message != nil && update.Message.WebAppData != nil {
		return h.handleWebAppData(update)
	}

	if update.Message == nil {
		return nil
	}

	if update.Message.Text == "/start" {
		return h.handleStartCommand(update)
	}

	return h.handleRegularMessage(update)
}

// Обработка данных от WebApp
func (h *userHandler) handleWebAppData(update telego.Update) error {
	webAppData := update.Message.WebAppData.Data
	log.Printf("Получены данные из Web App: %s", webAppData)

	msg := telegoutil.Message(update.Message.Chat.ChatID(), "Данные получены из приложения: "+webAppData)
	_, err := h.bot.SendMessage(h.ctx, msg)
	return err
}

// Обработка команды /start
func (h *userHandler) handleStartCommand(update telego.Update) error {
	username := update.Message.From.Username
	firstname := update.Message.From.FirstName
	chatID := update.Message.Chat.ChatID()

	if err := h.userSvc.ProcessUser(username, firstname, chatID); err != nil {
		log.Printf("Ошибка обработки пользователя: %v", err)
		return err
	}

	webAppBtn := telegoutil.KeyboardButton("🎮 Открыть игру").
		WithWebApp(&telego.WebAppInfo{URL: h.appURL})

	regularBtn := telegoutil.KeyboardButton("ℹ️ Информация")
	keyboard := telegoutil.Keyboard(
		telegoutil.KeyboardRow(webAppBtn),
		telegoutil.KeyboardRow(regularBtn))

	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	msg := telegoutil.Message(chatID,
		"🎉 Добро пожаловать в Dojo Bot!\n\n"+
			"✨ Теперь у вас есть доступ к мини-приложению через:\n"+
			"• Кнопку меню (рядом с полем ввода)\n"+
			"• Кнопку ниже\n\n"+
			"Нажмите любую из них, чтобы начать!",
	).WithReplyMarkup(keyboard)

	_, err := h.bot.SendMessage(h.ctx, msg)
	return err
}

// Обработка обычных сообщений
func (h *userHandler) handleRegularMessage(update telego.Update) error {
	chatID := update.Message.Chat.ID
	var response string

	switch update.Message.Text {
	case "ℹ️ Информация":
		response = "📱 Это Dojo Bot с мини-приложением!\n\n" +
			"🎮 Используйте кнопку меню или кнопку 'Открыть игру' для запуска приложения.\n" +
			"💡 Кнопка меню всегда доступна рядом с полем ввода сообщения."
	default:
		response = "Приветствую " + update.Message.From.Username + " 👋\n\n"
	}

	inlineKeyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text: "🔗 Перейти в канал",
					URL:  h.chanURL,
				},
			},
			{
				{
					Text: "🚀 Открыть игру",
					WebApp: &telego.WebAppInfo{
						URL: h.appURL,
					},
				},
			},
		},
	}

	msg := telego.SendMessageParams{
		ChatID:      telego.ChatID{ID: chatID},
		Text:        response,
		ReplyMarkup: inlineKeyboard,
	}
	_, err := h.bot.SendMessage(h.ctx, &msg)
	return err
}
