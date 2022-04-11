package postgres

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/readModel"

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

	err := r.db.Create(ToMessagePersistence(chatMessage)).Error

	return err
}

func (r *messageRepository) FindAllByConversationID(conversationID uuid.UUID, requestUserID uuid.UUID) ([]*readModel.MessageDTO, error) {
	messages := []*Message{}

	err := r.db.Limit(50).Where("conversation_id = ?", conversationID).Find(&messages).Error

	dtoMessages := make([]*readModel.MessageDTO, len(messages))

	for i, message := range messages {
		user := User{}

		if message.Type == 0 {
			err := r.db.Where("id = ?", message.UserID).First(&user).Error

			if err != nil {
				return nil, err
			}
		}

		dtoMessages[i] = ToMessageDTO(message, &user, requestUserID)
	}

	return dtoMessages, err
}

func (r *messageRepository) GetMessageByID(messageID uuid.UUID, requestUserID uuid.UUID) (*readModel.MessageDTO, error) {
	message := Message{}

	err := r.db.Where("id = ?", messageID).First(&message).Error

	if err != nil {
		return nil, err
	}

	user := User{}

	if message.Type == 0 {
		err = r.db.Where("id = ?", message.UserID).First(&user).Error

		if err != nil {
			return nil, err
		}
	}

	dtoMessage := ToMessageDTO(&message, &user, requestUserID)

	return dtoMessage, err
}
