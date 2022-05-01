package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageRepository struct {
	db     *gorm.DB
	pubsub domain.EventPublisher
}

func NewMessageRepository(db *gorm.DB, pubsub domain.EventPublisher) *messageRepository {
	return &messageRepository{
		db:     db,
		pubsub: pubsub,
	}
}

func (r *messageRepository) StoreTextMessage(message *domain.TextMessage) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toTextMessagePersistence(*message)).Error

	if err != nil {
		return err
	}

	message.Dispatch(r.pubsub)

	return nil
}

func (r *messageRepository) StoreLeftConversationMessage(message *domain.Message) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	message.Dispatch(r.pubsub)

	return nil
}

func (r *messageRepository) StoreJoinedConversationMessage(message *domain.Message) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	message.Dispatch(r.pubsub)

	return nil
}

func (r *messageRepository) StoreRenamedConversationMessage(message *domain.ConversationRenamedMessage) error {
	err := r.db.Create(toMessagePersistence(message)).Error

	if err != nil {
		return err
	}

	err = r.db.Create(toRenameConversationMessagePersistence(*message)).Error

	if err != nil {
		return err
	}

	message.Dispatch(r.pubsub)

	return nil
}

func (r *messageRepository) GetConversationMessages(conversationID uuid.UUID, requestUserID uuid.UUID) ([]*readModel.MessageDTO, error) {
	messages := []*Message{}

	err := r.db.Limit(50).Where("conversation_id = ?", conversationID).Find(&messages).Error

	dtoMessages := make([]*readModel.MessageDTO, len(messages))

	for i, message := range messages {
		msgDTO, err := r.getMessageDTO(message, requestUserID)

		if err != nil {
			return nil, err
		}

		dtoMessages[i] = msgDTO
	}

	return dtoMessages, err
}

func (r *messageRepository) GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (*readModel.MessageDTO, error) {
	message := Message{}

	err := r.db.Where("id = ?", messageID).First(&message).Error

	if err != nil {
		return nil, err
	}

	return r.getMessageDTO(&message, requestUserID)
}

func (r *messageRepository) getMessageDTO(message *Message, requestUserID uuid.UUID) (*readModel.MessageDTO, error) {
	user := User{}

	err := r.db.Where("id = ?", message.UserID).First(&user).Error

	if err != nil {
		return nil, err
	}

	switch message.Type {
	case 0:
		textMessage := TextMessage{}

		err = r.db.Where("message_id = ?", message.ID).First(&textMessage).Error

		if err != nil {
			return nil, err
		}

		return ToTextMessageDTO(message, &user, textMessage.Text, requestUserID), nil
	case 1:
		conversationRenamedMessage := ConversationRenamedMessage{}

		err = r.db.Where("message_id = ?", message.ID).First(&conversationRenamedMessage).Error

		if err != nil {
			return nil, err
		}

		return toConversationRenamedMessageDTO(message, &user, conversationRenamedMessage.NewName), nil
	case 2:
		return toMessageDTO(message, &user), nil
	case 3:
		return toMessageDTO(message, &user), nil

	}

	return nil, nil
}
