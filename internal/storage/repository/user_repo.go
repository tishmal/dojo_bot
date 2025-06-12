package repository

import (
	"context"
	"dojo_bot/model"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepoInterface interface {
	SaveUser(user model.User) error
	//GetUser(id int) (*User, error)
}

type UserRepo struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) UserRepoInterface {
	return &UserRepo{db: db}
}

// добавить индексы к полям и проверку на существование пользователя в бд, также редис или другой кэш
func (r *UserRepo) SaveUser(user model.User) error {
	res, err := r.db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Ошибка сохранения: %v", err)
		return err
	}

	log.Printf("Пользователь с ID: %v сохранён", res.InsertedID)
	return nil
}
