package service

import (
	"github.com/lavatee/camp_backend/internal/model"
	"github.com/lavatee/camp_backend/internal/repository"
)

type ChatsService struct {
	repo *repository.Repository
}

func NewChatsService(repo *repository.Repository) *ChatsService {
	return &ChatsService{
		repo: repo,
	}
}

func (s *ChatsService) CreateChat(firstUserId int, secondUserId int) (int, error) {
	return s.repo.Chats.CreateChat(firstUserId, secondUserId)
}

func (s *ChatsService) GetUserChats(userID int, searchQuery string) ([]model.Chat, error) {
	return s.repo.Chats.GetUserChats(userID, searchQuery)
}

func (s *ChatsService) GetOneChat(userId int, chatId int, timeZone string) (model.Chat, bool, error) {
	isTreeLegit, err := s.repo.Chats.CheckIsTreeLegit(chatId, timeZone)
	if err != nil {
		return model.Chat{}, false, err
	}
	messages, err := s.repo.Chats.GetChatMessages(chatId, userId, timeZone)
	if err != nil {
		return model.Chat{}, false, err
	}
	chat, err := s.repo.Chats.GetOneChat(chatId)
	if err != nil {
		return model.Chat{}, false, err
	}
	chat.Messages = messages
	return chat, isTreeLegit, nil
}

func (s *ChatsService) CreateMessage(message model.Message, timeZone string) (int, bool, error) {
	return s.repo.Chats.CreateMessage(message, timeZone)
}

func (s *ChatsService) EditMessage(messageId int, userId int, newText string) error {
	return s.repo.Chats.EditMessage(messageId, userId, newText)
}

func (s *ChatsService) DeleteMessage(messageId int, userId int) error {
	return s.repo.Chats.DeleteMessage(messageId, userId)
}

func (s *ChatsService) MakeMessageRead(messageId int) error {
	return s.repo.Chats.MakeMessageRead(messageId)
}
