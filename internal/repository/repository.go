package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lavatee/camp_backend/internal/model"
)

type Users interface {
	CreateUser(user model.User) (int, error)
	SignIn(email string, passwordHash string) (model.User, error)
	EditUserInfo(user model.User) error
	CheckTagUnique(tag string) bool
	FindUserByTag(tag string) (model.User, error)
	GetOneUser(userId int) (model.User, error)
}

type Chats interface {
	CreateChat(firstUserId int, secondUserId int) (int, error)
	GetUserChats(userID int, searchQuery string) ([]model.Chat, error)
	GetOneChat(chatId int) (model.Chat, error)
	CheckIsTreeLegit(chatId int, timeZone string) (bool, error)
	CreateMessage(message model.Message, timeZone string) (int, bool, error)
	GetChatMessages(chatId int, userId int, timeZone string) ([]model.Message, error)
	EditMessage(messageId int, userId int, newText string) error
	DeleteMessage(messageId int, userId int) error
	MakeMessageRead(messageId int) error
}

type Rooms interface {
	JoinRoom(userId int) (model.Room, error)
	LeaveRoom(userId int, roomId int) error
	GetRoomUser(userId int, roomId int) (model.Room, error)
}

type Repository struct {
	Users
	Chats
	Rooms
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Users: NewUsersPostgres(db),
		Chats: NewChatsPostgres(db),
		Rooms: NewRoomsPostgres(db),
	}
}
