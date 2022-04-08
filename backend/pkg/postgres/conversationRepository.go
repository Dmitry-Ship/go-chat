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

	err := r.conversations.Create(domain.ToConversationPersistence(conversation)).Error

	return err
}

func (r *conversationRepository) FindByID(id uuid.UUID) (*domain.ConversationDTO, error) {
	conversation := domain.ConversationPersistence{}

	err := r.conversations.Where("id = ?", id).First(&conversation).Error

	return domain.ToConversationDTO(&conversation), err
}

func (r *conversationRepository) FindAll() ([]*domain.ConversationDTO, error) {
	conversations := []*domain.ConversationPersistence{}

	err := r.conversations.Limit(50).Find(&conversations).Error

	dtoConversations := make([]*domain.ConversationDTO, len(conversations))

	for i, conversation := range conversations {
		dtoConversations[i] = domain.ToConversationDTO(conversation)
	}

	return dtoConversations, err
}

func (r *conversationRepository) Delete(id uuid.UUID) error {
	conversation := domain.ConversationPersistence{}

	err := r.conversations.Where("id = ?", id).Delete(conversation).Error

	return err
}
