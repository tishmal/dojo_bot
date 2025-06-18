package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(ctx context.Context, uri string) (*mongo.Client, error) {
	// Создаем новый контекст с таймаутом на основе переданного
	connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel() // Важно вызывать cancel для освобождения ресурсов

	client, err := mongo.Connect(connCtx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("MongoDB connection error: %w", err)
	}

	// Проверяем подключение
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}
func GetDatabase(client *mongo.Client, name string) *mongo.Database {
	return client.Database(name)
}
