package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id"`                // Авто-генерация
	TelegramID int64              `bson:"telegram_id,unique"` // Индекс!
	Profile    struct {
		Username  string    `bson:"username"`
		FirstName string    `bson:"first_name"`
		LastSeen  time.Time `bson:"last_seen"`
	} `bson:"profile"`
	Settings  map[string]interface{} `bson:"settings"`   // Гибкие настройки
	CreatedAt time.Time              `bson:"created_at"` // TTL индекс
}
