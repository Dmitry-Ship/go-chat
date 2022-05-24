package postgres

import (
	"GitHub/go-chat/backend/internal/readModel"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type queriesRepository struct {
	db *gorm.DB
}

func NewQueriesRepository(db *gorm.DB) *queriesRepository {
	return &queriesRepository{
		db: db,
	}
}

func (r *queriesRepository) paginate(paginationInfo readModel.PaginationInfo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := paginationInfo.GetPage()

		if page == 0 {
			page = 1
		}

		pageSize := paginationInfo.GetPageSize()

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 50
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (r *queriesRepository) GetContacts(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.ContactDTO, error) {
	users := []*readModel.ContactDTO{}

	err := r.db.Scopes(r.paginate(paginationInfo)).Model(&User{}).Where("id <> ?", userID).Find(&users).Error

	return users, err
}

func (r *queriesRepository) GetParticipants(conversationID uuid.UUID, userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.ContactDTO, error) {
	users := []*readModel.ContactDTO{}

	err := r.db.Scopes(r.paginate(paginationInfo)).
		Model(&User{}).
		Joins("left join participants on participants.user_id = users.id").
		Where("participants.conversation_id", conversationID).
		Find(&users).Error

	return users, err
}

func (r *queriesRepository) GetPotentialInvitees(conversationID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.ContactDTO, error) {
	users := []*readModel.ContactDTO{}

	subQuery := r.db.Select("user_id").Where("conversation_id = ?", conversationID).Table("participants")
	err := r.db.Scopes(r.paginate(paginationInfo)).Model(&User{}).Where("id NOT IN (?)", subQuery).Find(&users).Error

	return users, err
}

func (r *queriesRepository) GetUserByID(id uuid.UUID) (*readModel.UserDTO, error) {
	user := &readModel.UserDTO{}

	err := r.db.Model(&User{}).Where("id = ?", id).First(&user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *queriesRepository) GetConversationMessages(conversationID uuid.UUID, requestUserID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.MessageDTO, error) {
	messages := []*readModel.MessageDTO{}

	err := r.db.Scopes(r.paginate(paginationInfo)).
		Model(&Message{}).
		Select(
			"messages.id", "messages.type as persistence_type", "messages.created_at", "text_messages.text",
			"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
			"conversation_renamed_messages.new_name",
		).
		Joins("LEFT JOIN conversation_renamed_messages ON messages.id = conversation_renamed_messages.message_id").
		Joins("LEFT JOIN users ON messages.user_id = users.id").
		Joins("LEFT JOIN text_messages ON messages.id = text_messages.message_id").
		Where("messages.conversation_id = ?", conversationID).
		Order("messages.created_at asc").
		Find(&messages).Error

	for _, message := range messages {
		message.User = &readModel.UserDTO{
			ID:     message.UserID,
			Avatar: message.UserAvatar,
			Name:   message.UserName,
		}
		message.Type = messageTypesMap[message.PersistenceType]
		message.IsInbound = message.UserID != requestUserID
	}

	return messages, err
}

func (r *queriesRepository) GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (*readModel.MessageDTO, error) {
	message := &readModel.MessageDTO{}

	err := r.db.Model(&Message{}).
		Select(
			"messages.id", "messages.type as persistence_type", "messages.created_at", "text_messages.text",
			"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
			"conversation_renamed_messages.new_name",
		).
		Joins("LEFT JOIN conversation_renamed_messages ON messages.id = conversation_renamed_messages.message_id").
		Joins("LEFT JOIN users ON messages.user_id = users.id").
		Joins("LEFT JOIN text_messages ON messages.id = text_messages.message_id").
		Where("messages.id = ?", messageID).
		Find(&message).Error

	message.User = &readModel.UserDTO{
		ID:     message.UserID,
		Avatar: message.UserAvatar,
		Name:   message.UserName,
	}
	message.Type = messageTypesMap[message.PersistenceType]
	message.IsInbound = message.UserID != requestUserID

	return message, err
}

func (r *queriesRepository) GetUserConversations(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.ConversationDTO, error) {
	conversations := []*readModel.ConversationDTO{}

	err := r.db.Scopes(r.paginate(paginationInfo)).
		Model(&Conversation{}).
		Select(
			"conversations.id", "conversations.created_at", "conversations.type",
			"group_conversations.avatar", "group_conversations.name",
			"direct_conversations.from_user_id", "direct_conversations.to_user_id",
			"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
		).
		Joins("LEFT JOIN group_conversations ON group_conversations.conversation_id = conversations.id").
		Joins("LEFT JOIN direct_conversations ON direct_conversations.conversation_id = conversations.id").
		Joins("LEFT JOIN participants ON participants.conversation_id = conversations.id").
		Joins("LEFT JOIN users ON direct_conversations.from_user_id = users.id OR direct_conversations.to_user_id = users.id").
		Where("users.id IS NULL OR users.id <> ?", userID).
		Where("conversations.is_active = ?", true).
		Where("participants.is_active = ?", true).
		Where("participants.user_id = ?", userID).
		Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	for _, conversation := range conversations {
		if conversation.Type == 1 {
			conversation.Avatar = conversation.UserAvatar
			conversation.Name = conversation.UserName
		}
	}

	return conversations, nil
}

func (r *queriesRepository) GetConversation(id uuid.UUID, userId uuid.UUID) (*readModel.ConversationFullDTO, error) {
	conversation := readModel.ConversationFullDTO{}

	err := r.db.Model(&Conversation{}).
		Select(
			"conversations.id", "conversations.created_at", "conversations.type as persistence_type",
			"group_conversations.avatar", "group_conversations.name",
			"direct_conversations.from_user_id", "direct_conversations.to_user_id",
			"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
		).
		Joins("LEFT JOIN group_conversations ON group_conversations.conversation_id = conversations.id").
		Joins("LEFT JOIN direct_conversations ON direct_conversations.conversation_id = conversations.id").
		Joins("LEFT JOIN users ON direct_conversations.from_user_id = users.id OR direct_conversations.to_user_id = users.id").
		Where("users.id IS NULL OR users.id <> ?", userId).
		Where("conversations.is_active = ?", true).
		Where("conversations.id = ?", id).
		Find(&conversation).Error

	if err != nil {
		return nil, err
	}

	conversation.Type = conversationTypesMap[conversation.PersistenceType]

	if conversation.PersistenceType == 1 {
		conversation.Avatar = conversation.UserAvatar
		conversation.Name = conversation.UserName
	}

	hasUserJoined := true

	participant := Participant{}
	if err := r.db.Where("conversation_id = ?", id).Where("user_id = ?", userId).Where("is_active = ?", true).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hasUserJoined = false
		} else {
			return nil, err
		}
	}

	var participantsCount int64

	err = r.db.Model(&Participant{}).Where("conversation_id = ?", id).Where("is_active = ?", true).Count(&participantsCount).Error

	if err != nil {
		return nil, err
	}

	conversation.HasJoined = hasUserJoined
	conversation.ParticipantsCount = participantsCount

	return &conversation, nil
}
