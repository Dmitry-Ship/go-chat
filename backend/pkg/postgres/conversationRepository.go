package postgres

import (
	"GitHub/go-chat/backend/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *conversationRepository {
	return &conversationRepository{
		db: db,
	}

}

func (r *conversationRepository) Store(conversation *domain.Conversation) error {
	err := r.db.Create(domain.ToConversationPersistence(conversation)).Error

	return err
}

func (r *conversationRepository) FindByID(id uuid.UUID, userId uuid.UUID) (*domain.ConversationDTOFull, error) {
	conversation := domain.ConversationPersistence{}

	err := r.db.Where("id = ?", id).First(&conversation).Error

	if err != nil {
		return nil, err
	}

	hasUserJoined := true

	participant := domain.ParticipantPersistence{}
	if err := r.db.Where("conversation_id = ?", id).Where("user_id = ?", userId).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hasUserJoined = false
		} else {
			return nil, err
		}
	}

	return domain.ToConversationDTOFull(&conversation, hasUserJoined), err
}

func (r *conversationRepository) FindAll() ([]*domain.ConversationDTO, error) {
	conversations := []*domain.ConversationPersistence{}

	err := r.db.Limit(50).Find(&conversations).Error

	dtoConversations := make([]*domain.ConversationDTO, len(conversations))

	for i, conversation := range conversations {
		dtoConversations[i] = domain.ToConversationDTO(conversation)
	}

	return dtoConversations, err
}

func (r *conversationRepository) Delete(id uuid.UUID) error {
	err := r.db.Where("id = ?", id).Delete(domain.ConversationPersistence{}).Error

	return err
}
