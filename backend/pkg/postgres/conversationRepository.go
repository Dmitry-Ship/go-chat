package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/mappers"

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

	err := r.conversations.Create(mappers.ToConversationPersistence(conversation)).Error

	return err
}

func (r *conversationRepository) FindByID(id uuid.UUID) (*domain.Conversation, error) {
	conversation := mappers.ConversationPersistence{}

	err := r.conversations.Where("id = ?", id).First(&conversation).Error

	return mappers.ToConversationDomain(&conversation), err
}

func (r *conversationRepository) FindAll() ([]*domain.Conversation, error) {
	conversations := []*mappers.ConversationPersistence{}

	err := r.conversations.Limit(50).Find(&conversations).Error

	domainConversations := make([]*domain.Conversation, len(conversations))

	for i, conversation := range conversations {
		domainConversations[i] = mappers.ToConversationDomain(conversation)
	}

	return domainConversations, err
}

func (r *conversationRepository) Delete(id uuid.UUID) error {
	conversation := mappers.ConversationPersistence{}

	err := r.conversations.Where("id = ?", id).Delete(conversation).Error

	return err
}
