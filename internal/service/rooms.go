package service

import (
	"github.com/lavatee/camp_backend/internal/model"
	"github.com/lavatee/camp_backend/internal/repository"
)

type RoomsService struct {
	repo *repository.Repository
}

func NewRoomsService(repo *repository.Repository) *RoomsService {
	return &RoomsService{
		repo: repo,
	}
}

func (s *RoomsService) JoinRoom(userId int) (model.Room, error) {
	return s.repo.Rooms.JoinRoom(userId)
}

func (s *RoomsService) LeaveRoom(userId int, roomId int) error {
	return s.repo.Rooms.LeaveRoom(userId, roomId)
}

func (s *RoomsService) NextRoom(userId int, currentRoomId int) (model.Room, error) {
	room, err := s.repo.Rooms.JoinRoom(userId)
	if err != nil {
		return model.Room{}, err
	}
	if err := s.repo.Rooms.LeaveRoom(userId, currentRoomId); err != nil {
		return model.Room{}, err
	}
	return room, nil
}

func (s *RoomsService) GetRoomUser(userId int, roomId int) (model.Room, error) {
	return s.repo.Rooms.GetRoomUser(userId, roomId)
}
