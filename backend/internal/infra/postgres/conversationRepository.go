package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
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

	if err != nil {
		return err
	}

	err = r.db.Create(toParticipantPersistence(&conversation.Data.Owner)).Error

	return err
}

func (r *conversationRepository) UpdatePublicConversation(conversation *domain.PublicConversation) error {
	err := r.db.Save(toConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Save(toPublicConversationPersistence(conversation)).Error

	return err
}

func (r *conversationRepository) StorePrivateConversation(conversation *domain.PrivateConversation) error {
	err := r.db.Create(toConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toPrivateConversationPersistence(conversation)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toParticipantPersistence(conversation.GetFromUser())).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toParticipantPersistence(conversation.GetToUser())).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *conversationRepository) GetPrivateConversationID(firstUserId uuid.UUID, secondUserID uuid.UUID) (uuid.UUID, error) {
	privateConversation := PrivateConversation{}
	err := r.db.Where("to_user_id = ? AND from_user_id = ?", firstUserId, secondUserID).Or("to_user_id = ? AND from_user_id = ?", secondUserID, firstUserId).First(&privateConversation).Error

	if err != nil {
		return uuid.Nil, err
	}

	return privateConversation.ConversationID, nil
}

func (r *conversationRepository) GetPublicConversation(id uuid.UUID) (*domain.PublicConversation, error) {
	conversation := Conversation{}

	err := r.db.Where("id = ?", id).First(&conversation).Error

	if err != nil {
		return nil, err
	}

	publicConversation := PublicConversation{}

	err = r.db.Where("conversation_id = ?", id).First(&publicConversation).Error

	if err != nil {
		return nil, err
	}

	participant := Participant{}

	err = r.db.Where("conversation_id = ? AND user_id = ?", id, publicConversation.OwnerID).First(&participant).Error

	if err != nil {
		return nil, err
	}

	return toPublicConversationDomain(&conversation, &publicConversation, &participant), nil
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
	case 1:
		privateConversation := PrivateConversation{}

		err = r.db.Where("conversation_id = ?", conversation.ID).First(&privateConversation).Error

		if err != nil {
			return nil, err
		}

		user := User{}

		oppositeUserId := privateConversation.FromUserID

		if privateConversation.FromUserID == userId {
			oppositeUserId = privateConversation.ToUserID
		}

		err = r.db.Where("id = ?", oppositeUserId).First(&user).Error

		if err != nil {
			return nil, err
		}

		return toConversationFullDTO(&conversation, user.Avatar, user.Name, hasUserJoined), nil
	default:
		return nil, errors.New("unsupported conversation type")
	}
}

func (r *conversationRepository) GetUserConversations(userID uuid.UUID) ([]*readModel.ConversationDTO, error) {
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

		case 1:
			privateConversation := PrivateConversation{}

			err = r.db.Where("conversation_id = ?", conversation.ID).First(&privateConversation).Error

			if err != nil {
				return nil, err
			}

			user := User{}

			oppositeUserId := privateConversation.FromUserID
			if privateConversation.FromUserID == userID {
				oppositeUserId = privateConversation.ToUserID
			}

			err = r.db.Where("id = ?", oppositeUserId).First(&user).Error

			if err != nil {
				return nil, err
			}

			conversationsDTOs[i] = toPrivateConversationDTO(conversation, &user)
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
