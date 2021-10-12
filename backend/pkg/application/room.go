package application

import (
	"GitHub/go-chat/backend/domain"
	"fmt"

	"github.com/google/uuid"
)

type RoomService interface {
	CreateRoom(name string, userId uuid.UUID) (*domain.Room, error)
	GetRoom(id uuid.UUID) (*domain.Room, error)
	HasJoined(userId uuid.UUID, roomId uuid.UUID) bool
	GetRooms() ([]*domain.Room, error)
	JoinRoom(userId uuid.UUID, roomId uuid.UUID) (*domain.Participant, error)
	LeaveRoom(userId uuid.UUID, roomId uuid.UUID) error
	DeleteRoom(id uuid.UUID) error
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

func (s *roomService) CreateRoom(name string, userId uuid.UUID) (*domain.Room, error) {
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

func (s *roomService) GetRoom(id uuid.UUID) (*domain.Room, error) {
	return s.rooms.FindByID(id)
}

func (s *roomService) GetRooms() ([]*domain.Room, error) {
	return s.rooms.FindAll()
}

func (s *roomService) JoinRoom(userId uuid.UUID, roomID uuid.UUID) (*domain.Participant, error) {
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

func (s *roomService) LeaveRoom(userId uuid.UUID, roomID uuid.UUID) error {
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

func (s *roomService) DeleteRoom(id uuid.UUID) error {
	room, err := s.rooms.FindByID(id)

	if err != nil {
		return err
	}

	s.participants.DeleteByRoomID(id)
	s.rooms.Delete(id)

	s.hub.DeleteRoom(room.Id)

	return nil
}

func (s *roomService) HasJoined(userId uuid.UUID, roomID uuid.UUID) bool {
	_, err := s.participants.FindByRoomIDAndUserID(roomID, userId)

	return err == nil
}
