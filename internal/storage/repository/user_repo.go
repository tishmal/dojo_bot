package repository

import (
	"context"
	"dojo_bot/model"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	SaveUser(db *mongo.Database, user model.User) error
}

type UserRepo struct {
	db *mongo.Database
}

func NewMongoDB(db *mongo.Database) UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) SaveUser(db *mongo.Database, user model.User) error {
	res, err := db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Ошибка сохранения: %v", err)
		return err
	}

	log.Printf("Пользователь с ID: %v сохранён", res.InsertedID)
	return nil
}
