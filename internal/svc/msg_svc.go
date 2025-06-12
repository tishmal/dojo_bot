package svc

import (
	"dojo_bot/internal/storage"
	"dojo_bot/model"
	"fmt"
	"time"
)



type UserService struct {
	repo storage.UserRepository
}

func (s *UserService) HandleMessage(msg *tgbotapi.Message) (string, error) {
	user := model.User{
		Username: msg.From.UserName,
		ChatID:   msg.Chat.ID,
		JoinedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.repo.SaveUser(user); err != nil {
		return "", fmt.Errorf("не удалось сохранить пользователя: %w", err)
	}

	return "Салам! Твой ник: " + user.Username, nil
}
