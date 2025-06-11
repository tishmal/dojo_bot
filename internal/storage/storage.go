package storage

import (
	"context"
	"dojo_bot/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func SaveUser(db *mongo.Database, user model.User) error {
	res, err := db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Ошибка сохранения: %v", err)
		return err
	}

	log.Printf("Пользователь с ID: %v сохранён", res.InsertedID)
	return nil
}
