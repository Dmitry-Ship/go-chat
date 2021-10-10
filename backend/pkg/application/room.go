package application

import (
	"GitHub/go-chat/backend/domain"
	"fmt"
)

type RoomService interface {
	CreateRoom(name string, userId int32) (*domain.Room, error)
	GetRoom(id int32) (*domain.Room, error)
	GetRooms() ([]*domain.Room, error)
	JoinRoom(userId int32, roomId int32) (*domain.Participant, error)
	LeaveRoom(userId int32, roomId int32) error
}

type roomService struct {
	rooms          domain.RoomRepository
	participants   domain.ParticipantRepository
	userService    UserService
	messageService MessageService
}

func NewRoomService(rooms domain.RoomRepository, participants domain.ParticipantRepository, userService UserService, messageService MessageService) *roomService {
	return &roomService{
		rooms:          rooms,
		userService:    userService,
		participants:   participants,
		messageService: messageService,
	}
}

func (s *roomService) CreateRoom(name string, userId int32) (*domain.Room, error) {
	room := domain.NewRoom(name)
	newRoom, err := s.rooms.Create(room)

	if err != nil {
		return nil, err
	}

	s.JoinRoom(userId, room.Id)

	return newRoom, nil
}

func (s *roomService) GetRoom(id int32) (*domain.Room, error) {
	return s.rooms.FindByID(id)
}

func (s *roomService) GetRooms() ([]*domain.Room, error) {
	return s.rooms.FindAll()
}

func (s *roomService) JoinRoom(userId int32, roomID int32) (*domain.Participant, error) {
	participant := domain.NewParticipant(roomID, userId)
	newParticipant, err := s.participants.Create(participant)

	if err != nil {
		return nil, err
	}

	user, err := s.userService.GetUser(userId)

	if err != nil {
		fmt.Println("err", err)

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

	user, err := s.userService.GetUser(userId)

	if err != nil {
		return err
	}

	s.messageService.SendMessage(fmt.Sprintf("%s %s left", user.Avatar, user.Name), "system", roomID, user.Id)
	return nil
}
