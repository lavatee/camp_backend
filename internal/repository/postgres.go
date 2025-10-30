package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable       = "users"
	chatsTable       = "chats"
	usersInChatTable = "users_in_chat"
	messagesTable    = "messages"
	roomsTable       = "rooms"
	usersInRoomTable = "users_in_room"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(config PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
