package svc

import (
	"dojo_bot/internal/storage/repository"
	"dojo_bot/model"
	"fmt"
	"time"
)

type UserSvcInterface interface {
	ProcessUser(username string, chatID int64) error
}

type UserService struct {
	repo repository.UserRepoInterface
}

func NewUserSvc(userRepo *repository.UserRepo) UserSvcInterface {
	return &UserService{repo: userRepo}
}

func (s *UserService) ProcessUser(username string, chatID int64) error {
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
