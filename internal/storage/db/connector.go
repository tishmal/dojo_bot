package db

import (
	"context"
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
		log.Printf("Ошибка подключения к MongoDB: %v", err)
		return nil, err
	}
	log.Printf("Успешно")
	return client, nil
}

func GetDatabase(client *mongo.Client, name string) *mongo.Database {
	return client.Database(name)
}
