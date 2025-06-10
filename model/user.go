package model

type User struct {
	Username string `bson:"username"`
	ChatID   int64  `bson:"chat_id"`
	JoinedAt string `bson:"joined_at"`
}
