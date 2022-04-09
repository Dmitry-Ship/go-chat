package postgres

import (
	"GitHub/go-chat/backend/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *messageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) Store(chatMessage *domain.Message) error {

	err := r.db.Create(domain.ToMessagePersistence(chatMessage)).Error

	return err
}

func (r *messageRepository) FindAllByConversationID(conversationID uuid.UUID, requestUserID uuid.UUID) ([]*domain.MessageDTO, error) {
	messages := []*domain.MessagePersistence{}

	err := r.db.Limit(50).Where("conversation_id = ?", conversationID).Find(&messages).Error

	dtoMessages := make([]*domain.MessageDTO, len(messages))

	for i, message := range messages {
		user := domain.UserPersistence{}

		if message.Type == 0 {
			err := r.db.Where("id = ?", message.UserID).First(&user).Error

			if err != nil {
				return nil, err
			}
		}

		dtoMessages[i] = domain.ToMessageDTO(message, &user, requestUserID)
	}

	return dtoMessages, err
}

func (r *messageRepository) FindByID(messageID uuid.UUID, requestUserID uuid.UUID) (*domain.MessageDTO, error) {
	message := domain.MessagePersistence{}

	err := r.db.Where("id = ?", messageID).Find(message).Error

	user := domain.UserPersistence{}

	if message.Type == 0 {
		err := r.db.Where("id = ?", message.UserID).Find(&user).Error

		if err != nil {
			return nil, err
		}
	}

	dtoMessage := domain.ToMessageDTO(&message, &user, requestUserID)

	return dtoMessage, err
}
