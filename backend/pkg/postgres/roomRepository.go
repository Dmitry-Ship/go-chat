package postgres

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roomRepository struct {
	rooms *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *roomRepository {
	return &roomRepository{
		rooms: db,
	}

}

func (r *roomRepository) Store(room *domain.Room) error {
	err := r.rooms.Create(&room).Error

	return err
}

func (r *roomRepository) FindByID(id uuid.UUID) (*domain.Room, error) {
	room := domain.Room{}
	err := r.rooms.Where("id = ?", id).First(&room).Error

	return &room, err
}

func (r *roomRepository) FindAll() ([]*domain.Room, error) {
	rooms := []*domain.Room{}

	err := r.rooms.Limit(50).Find(&rooms).Error

	return rooms, err
}

func (r *roomRepository) Delete(id uuid.UUID) error {
	participant := domain.Participant{}

	err := r.rooms.Where("id = ?", id).Delete(participant).Error

	return err
}
