package inmemory

import (
	"GitHub/go-chat/backend/domain"
	"errors"
)

type roomRepository struct {
	rooms map[int32]*domain.Room
}

func NewRoomRepository() *roomRepository {
	return &roomRepository{
		rooms: make(map[int32]*domain.Room),
	}

}

func (r *roomRepository) FindByID(id int32) (*domain.Room, error) {
	room, ok := r.rooms[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return room, nil
}

func (r *roomRepository) FindByName(name string) (*domain.Room, error) {
	for _, room := range r.rooms {
		if room.Name == name {
			return room, nil
		}
	}
	return nil, errors.New("not found")
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

func (r *roomRepository) Update(room *domain.Room) error {
	_, ok := r.rooms[room.Id]
	if !ok {
		return errors.New("not found")
	}
	r.rooms[room.Id] = room
	return nil
}

func (r *roomRepository) Delete(id int32) error {
	_, ok := r.rooms[id]
	if !ok {
		return errors.New("not found")
	}
	delete(r.rooms, id)
	return nil
}
