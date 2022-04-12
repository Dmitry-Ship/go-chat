package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"
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

func (r *conversationRepository) Store(conversation *domain.ConversationAggregate) error {
	err := r.db.Create(toConversationPersistence(conversation)).Error

	return err
}

func (r *conversationRepository) RenameConversation(id uuid.UUID, name string) error {
	err := r.db.Model(Conversation{}).Where("id = ?", id).Update("name", name).Error

	return err
}

func (r *conversationRepository) GetConversationByID(id uuid.UUID, userId uuid.UUID) (*readModel.ConversationFullDTO, error) {
	conversation := Conversation{}

	err := r.db.Where("id = ?", id).First(&conversation).Error

	if err != nil {
		return nil, err
	}

	hasUserJoined := true

	participant := Participant{}
	if err := r.db.Where("conversation_id = ?", id).Where("user_id = ?", userId).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hasUserJoined = false
		} else {
			return nil, err
		}
	}

	return toConversationFullDTO(&conversation, hasUserJoined), err
}

func (r *conversationRepository) FindAll() ([]*readModel.ConversationDTO, error) {
	conversations := []*Conversation{}

	err := r.db.Limit(50).Find(&conversations).Error

	conversationsDTOs := make([]*readModel.ConversationDTO, len(conversations))

	for i, conversation := range conversations {
		conversationsDTOs[i] = toConversationDTO(conversation)
	}

	return conversationsDTOs, err
}

func (r *conversationRepository) Delete(id uuid.UUID) error {
	err := r.db.Where("id = ?", id).Delete(Conversation{}).Error

	return err
}
