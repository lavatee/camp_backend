package model

import (
	"time"
)

type Chat struct {
	Id              int       `json:"id" db:"id"`
	UserId          int       `json:"user_id" db:"user_id"`
	UserPhotoURL    string    `json:"user_photo_url" db:"user_photo_url"`
	UserName        string    `json:"user_name" db:"user_name"`
	LastMessage     string    `json:"last_message" db:"last_message"`
	LastMessageTime time.Time `json:"last_message_time" db:"last_message_time"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	MessagesAmount  int       `json:"messages_amount" db:"messages_amount"`
	Tree            int       `json:"tree" db:"tree"`
	Messages        []Message `json:"messages" db:"messages"`
}

type UserInChat struct {
	Id     int `json:"id" db:"id"`
	ChatId int `json:"chat_id" db:"chat_id"`
	UserId int `json:"user_id" db:"user_id"`
}

type Message struct {
	Id            int       `json:"id" db:"id"`
	UserId        int       `json:"user_id" db:"user_id"`
	ChatId        int       `json:"chat_id" db:"chat_id"`
	Text          string    `json:"text" db:"text"`
	SentAt        time.Time `json:"sent_at" db:"sent_at"`
	IsUserMessage bool      `json:"is_user_message" db:"is_user_message"`
	IsRead        bool      `json:"is_read" db:"is_read"`
}
