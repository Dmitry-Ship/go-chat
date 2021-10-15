package inmemory

import (
	"GitHub/go-chat/backend/domain"
	"errors"

	"github.com/google/uuid"
)

type roomRepository struct {
	rooms map[uuid.UUID]*domain.Room
}

func NewRoomRepository() *roomRepository {
	return &roomRepository{
		rooms: make(map[uuid.UUID]*domain.Room),
	}

}

func (r *roomRepository) FindByID(id uuid.UUID) (*domain.Room, error) {
	room, ok := r.rooms[id]
	if !ok {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (r *roomRepository) FindAll() ([]*domain.Room, error) {
	rooms := make([]*domain.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *roomRepository) Create(room *domain.Room) (*domain.Room, error) {
	r.rooms[room.Id] = room

	return room, nil
}

func (r *roomRepository) Delete(id uuid.UUID) error {
	_, ok := r.rooms[id]
	if !ok {
		return errors.New("room not found")
	}
	delete(r.rooms, id)
	return nil
}
