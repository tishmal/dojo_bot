package pkg

import (
	"context"

	"github.com/mymmrac/telego"
)

func SetMenuButton(ctx context.Context, bot *telego.Bot) error {
	// Создаем WebApp кнопку для списка чатов
	menuButton := &telego.MenuButtonWebApp{
		Type: "web_app",
		Text: "🚀 Открыть приложение", // Текст в списке чатов
		WebApp: telego.WebAppInfo{
			URL: "https://tishmal.github.io/dojo-app/",
		},
	}

	// Устанавливаем кнопку
	return bot.SetChatMenuButton(ctx, &telego.SetChatMenuButtonParams{
		MenuButton: menuButton,
	})
}
