package postgres

import (
	"GitHub/go-chat/backend/pkg/domain"
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

func (r *conversationRepository) StorePublicConversation(conversation *domain.PublicConversation) error {
	err := r.db.Create(toConversationPersistence(conversation)).Error

	if err != nil {

		return err

	}

	err = r.db.Create(toPublicConversationPersistence(conversation)).Error

	return err
}

func (r *conversationRepository) RenamePublicConversation(id uuid.UUID, name string) error {
	err := r.db.Model(PublicConversation{}).Where("conversation_id = ?", id).Update("name", name).Error

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

	switch conversation.Type {
	case 0:
		publicConversation := PublicConversation{}

		err = r.db.Where("conversation_id = ?", id).First(&publicConversation).Error

		if err != nil {
			return nil, err
		}

		return toConversationFullDTO(&conversation, publicConversation.Avatar, publicConversation.Name, hasUserJoined), nil

	default:
		return nil, errors.New("unsupported conversation type")
	}
}

func (r *conversationRepository) FindMyConversations(userID uuid.UUID) ([]*readModel.ConversationDTO, error) {
	conversations := []*Conversation{}

	err := r.db.Joins("JOIN participants ON participants.conversation_id = conversations.id").Where("participants.user_id = ?", userID).Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	conversationsDTOs := make([]*readModel.ConversationDTO, len(conversations))

	for i, conversation := range conversations {
		switch conversation.Type {
		case 0:
			publicConversation := PublicConversation{}

			err = r.db.Where("conversation_id = ?", conversation.ID).First(&publicConversation).Error

			if err != nil {
				return nil, err
			}

			conversationsDTOs[i] = toPublicConversationDTO(conversation, publicConversation.Avatar, publicConversation.Name)

		default:
			return nil, errors.New("unsupported conversation type")
		}

	}

	return conversationsDTOs, err
}

func (r *conversationRepository) Delete(id uuid.UUID) error {
	err := r.db.Where("id = ?", id).Delete(Conversation{}).Error

	return err
}
