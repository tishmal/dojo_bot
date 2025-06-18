package svc

import (
	"dojo_bot/internal/storage/repository"
	"dojo_bot/model"
	"fmt"
	"time"

	"github.com/mymmrac/telego"
)

type UserSvcInterface interface {
	ProcessUser(username string, chatID telego.ChatID) error
}

type UserService struct {
	repo repository.UserRepoInterface
}

func NewUserSvc(userRepo repository.UserRepoInterface) UserSvcInterface {
	return &UserService{repo: userRepo}
}

func (s *UserService) ProcessUser(username string, chatID telego.ChatID) error {
	user := model.User{
		Username: username,
		ChatID:   chatID,
		JoinedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.repo.SaveUser(user); err != nil {
		return fmt.Errorf("не удалось сохранить пользователя: %w", err)
	}

	return nil
}
