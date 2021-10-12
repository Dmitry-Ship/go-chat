package application

import (
	"GitHub/go-chat/backend/domain"
	"fmt"
)

type RoomService interface {
	CreateRoom(name string, userId int32) (*domain.Room, error)
	GetRoom(id int32) (*domain.Room, error)
	HasJoined(userId int32, roomID int32) bool
	GetRooms() ([]*domain.Room, error)
	JoinRoom(userId int32, roomId int32) (*domain.Participant, error)
	LeaveRoom(userId int32, roomId int32) error
	DeleteRoom(id int32) error
}

type roomService struct {
	rooms          domain.RoomRepository
	participants   domain.ParticipantRepository
	users          domain.UserRepository
	hub            Hub
	messageService MessageService
}

func NewRoomService(rooms domain.RoomRepository, participants domain.ParticipantRepository, users domain.UserRepository, messageService MessageService, hub Hub) *roomService {
	return &roomService{
		rooms:          rooms,
		users:          users,
		participants:   participants,
		hub:            hub,
		messageService: messageService,
	}
}

func (s *roomService) CreateRoom(name string, userId int32) (*domain.Room, error) {
	room := domain.NewRoom(name)
	newRoom, err := s.rooms.Create(room)

	if err != nil {
		return nil, err
	}

	_, err = s.JoinRoom(userId, room.Id)

	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (s *roomService) GetRoom(id int32) (*domain.Room, error) {
	return s.rooms.FindByID(id)
}

func (s *roomService) GetRooms() ([]*domain.Room, error) {
	return s.rooms.FindAll()
}

func (s *roomService) JoinRoom(userId int32, roomID int32) (*domain.Participant, error) {
	participant, err := s.participants.FindByRoomIDAndUserID(roomID, userId)

	if err == nil {
		return participant, nil
	}

	newParticipant, err := s.participants.Create(domain.NewParticipant(roomID, userId))

	if err != nil {
		return nil, err
	}

	s.hub.JoinRoom(userId, roomID)

	user, err := s.users.FindByID(userId)
	if err != nil {
		return nil, err
	}

	s.messageService.SendMessage(fmt.Sprintf("%s %s joined", user.Avatar, user.Name), "system", roomID, user.Id)

	return newParticipant, nil
}

func (s *roomService) LeaveRoom(userId int32, roomID int32) error {
	participant, err := s.participants.FindByUserID(userId)

	if err != nil {
		return err
	}

	s.participants.Delete(participant.Id)

	s.hub.LeaveRoom(userId, roomID)

	user, err := s.users.FindByID(userId)

	if err != nil {
		return err
	}

	s.messageService.SendMessage(fmt.Sprintf("%s %s left", user.Avatar, user.Name), "system", roomID, user.Id)
	return nil
}

func (s *roomService) DeleteRoom(id int32) error {
	room, err := s.rooms.FindByID(id)

	if err != nil {
		return err
	}

	s.participants.DeleteByRoomID(id)
	s.rooms.Delete(id)

	s.hub.DeleteRoom(room.Id)

	return nil
}

func (s *roomService) HasJoined(userId int32, roomID int32) bool {
	_, err := s.participants.FindByRoomIDAndUserID(roomID, userId)

	return err == nil
}
