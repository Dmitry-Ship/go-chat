package postgres

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type conversationRepository struct {
	conversations *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *conversationRepository {
	return &conversationRepository{
		conversations: db,
	}

}

func (r *conversationRepository) Store(conversation *domain.Conversation) error {
	err := r.conversations.Create(&conversation).Error

	return err
}

func (r *conversationRepository) FindByID(id uuid.UUID) (*domain.Conversation, error) {
	conversation := domain.Conversation{}
	err := r.conversations.Where("id = ?", id).First(&conversation).Error

	return &conversation, err
}

func (r *conversationRepository) FindAll() ([]*domain.Conversation, error) {
	conversations := []*domain.Conversation{}

	err := r.conversations.Limit(50).Find(&conversations).Error

	return conversations, err
}

func (r *conversationRepository) Delete(id uuid.UUID) error {
	participant := domain.Participant{}

	err := r.conversations.Where("id = ?", id).Delete(participant).Error

	return err
}
