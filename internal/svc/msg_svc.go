package svc

import (
	"dojo_bot/internal/storage/repository"
	"dojo_bot/model"
	"fmt"
	"time"

	"github.com/mymmrac/telego"
)

type UserSvcInterface interface {
	ProcessUser(user_name string, first_name string, chatID telego.ChatID) error
}

type UserService struct {
	repo repository.UserRepoInterface
}

func NewUserSvc(userRepo repository.UserRepoInterface) UserSvcInterface {
	return &UserService{repo: userRepo}
}

func (s *UserService) ProcessUser(user_name string, first_name string, chatID telego.ChatID) error {
	user := model.User{
		TelegramID: chatID.ID,
		Profile: struct {
			Username  string    `bson:"username"`
			FirstName string    `bson:"first_name"`
			LastSeen  time.Time `bson:"last_seen"`
		}{
			Username:  user_name,
			FirstName: first_name,
			LastSeen:  time.Now(),
		},
		Settings:  make(map[string]interface{}), // пустые настройки
		CreatedAt: time.Now(),
	}

	if err := s.repo.SaveUser(user); err != nil {
		return fmt.Errorf("не удалось сохранить пользователя: %w", err)
	}

	return nil
}
