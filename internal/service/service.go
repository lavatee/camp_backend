package service

import (
	"mime/multipart"

	"github.com/dgrijalva/jwt-go"
	"github.com/lavatee/camp_backend/internal/model"
	"github.com/lavatee/camp_backend/internal/repository"
	"github.com/minio/minio-go/v7"
)

type Users interface {
	SignUp(user model.User) (int, error)
	SignIn(email string, password string) (string, string, error)
	EditUserInfo(user model.User) error
	CheckTagUnique(tag string) bool
	FindUserByTag(tag string) (model.User, error)
	GetOneUser(userId int) (model.User, error)
	NewProfilePhoto(userId int, file multipart.File) (string, error)
	ParseToken(token string) (jwt.MapClaims, error)
	Refresh(refreshToken string) (string, string, error)
}

type Chats interface {
	CreateChat(firstUserId int, secondUserId int) (int, error)
	GetUserChats(userID int, searchQuery string) ([]model.Chat, error)
	GetOneChat(userId int, chatId int, timeZone string) (model.Chat, bool, error)
	CreateMessage(message model.Message, timeZone string) (int, bool, error)
	EditMessage(messageId int, userId int, newText string) error
	DeleteMessage(messageId int, userId int) error
	MakeMessageRead(messageId int) error
}

type Rooms interface {
	JoinRoom(userId int) (model.Room, error)
	LeaveRoom(userId int, roomId int) error
	GetRoomUser(userId int, roomId int) (model.Room, error)
	NextRoom(userId int, currentRoomId int) (model.Room, error)
}

type Service struct {
	Users
	Chats
	Rooms
}

func NewService(repo *repository.Repository, s3 *minio.Client, bucket string) *Service {
	return &Service{
		Users: NewUsersService(repo, s3, bucket),
		Chats: NewChatsService(repo),
		Rooms: NewRoomsService(repo),
	}
}
