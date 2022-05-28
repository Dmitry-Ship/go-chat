package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
	"errors"
	"time"

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
		Joins("LEFT JOIN participants on participants.user_id = users.id").
		Where("participants.conversation_id = ?", conversationID).
		Where("participants.is_active = ?", true).
		Find(&users).Error

	return users, err
}

func (r *queriesRepository) GetParticipantsCount(conversationID uuid.UUID) (int64, error) {
	var count int64

	err := r.db.Model(&Participant{}).
		Where("participants.conversation_id = ?", conversationID).
		Where("participants.is_active = ?", true).
		Count(&count).Error

	return count, err
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
	type queryResult struct {
		ID             uuid.UUID
		CreatedAt      time.Time
		Type           uint8
		UserID         uuid.UUID
		ConversationID uuid.UUID
		UserName       string
		UserAvatar     string
		Content        string
	}

	queryResults := []*queryResult{}

	err := r.db.Scopes(r.paginate(paginationInfo)).
		Model(&Message{}).
		Select(
			"messages.id", "messages.type as type", "messages.created_at", "messages.conversation_id", "messages.content",
			"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
		).
		Joins("LEFT JOIN users ON messages.user_id = users.id").
		Where("messages.conversation_id = ?", conversationID).
		Order("messages.created_at asc").
		Find(&queryResults).Error

	messages := make([]*readModel.MessageDTO, len(queryResults))

	for i, result := range queryResults {
		text := ""

		switch messageTypesMap[result.Type] {
		case domain.MessageTypeText:
			text = result.Content
		case domain.MessageTypeRenamedConversation:
			text = result.UserName + " renamed chat to " + result.Content
		case domain.MessageTypeJoinedConversation:
			text = result.UserName + " joined"
		case domain.MessageTypeLeftConversation:
			text = result.UserName + " left"
		case domain.MessageTypeInvitedConversation:
			text = result.UserName + " was invited"
		default:
			text = "Unknown message type"
		}

		messageDTO := &readModel.MessageDTO{
			ID:             result.ID,
			CreatedAt:      result.CreatedAt,
			Text:           text,
			Type:           messageTypesMap[result.Type].String(),
			ConversationId: result.ConversationID,
			User: &readModel.UserDTO{
				ID:     result.UserID,
				Avatar: result.UserAvatar,
				Name:   result.UserName,
			},
		}

		if messageTypesMap[result.Type] == domain.MessageTypeText {
			messageDTO.IsInbound = result.UserID != requestUserID
		}

		messages[i] = messageDTO
	}

	return messages, err
}

func (r *queriesRepository) GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (*readModel.MessageDTO, error) {
	type queryResult struct {
		ID             uuid.UUID
		CreatedAt      time.Time
		Type           uint8
		UserID         uuid.UUID
		ConversationID uuid.UUID
		UserName       string
		UserAvatar     string
		Content        string
	}

	message := &queryResult{}

	err := r.db.Model(&Message{}).
		Select(
			"messages.id", "messages.type as type", "messages.created_at", "messages.content", "messages.conversation_id",
			"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
		).
		Joins("LEFT JOIN users ON messages.user_id = users.id").
		Where("messages.id = ?", messageID).
		Find(&message).Error

	text := ""

	switch messageTypesMap[message.Type] {
	case domain.MessageTypeText:
		text = message.Content
	case domain.MessageTypeRenamedConversation:
		text = message.UserName + " renamed chat to " + message.Content
	case domain.MessageTypeJoinedConversation:
		text = message.UserName + " joined"
	case domain.MessageTypeLeftConversation:
		text = message.UserName + " left"
	case domain.MessageTypeInvitedConversation:
		text = message.UserName + " was invited"
	default:
		text = "Unknown message type"
	}

	massageDTO := &readModel.MessageDTO{
		ID:             message.ID,
		CreatedAt:      message.CreatedAt,
		Text:           text,
		Type:           messageTypesMap[message.Type].String(),
		ConversationId: message.ConversationID,
		User: &readModel.UserDTO{
			ID:     message.UserID,
			Avatar: message.UserAvatar,
			Name:   message.UserName,
		},
	}

	if messageTypesMap[message.Type] == domain.MessageTypeText {
		massageDTO.IsInbound = message.UserID != requestUserID
	}

	return massageDTO, err
}

func (r *queriesRepository) GetUserConversations(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.ConversationDTO, error) {
	type queryResult struct {
		ID         uuid.UUID
		Name       string
		Avatar     string
		CreatedAt  time.Time
		Type       uint8
		UserID     uuid.UUID
		UserAvatar string
		UserName   string
	}

	queryResults := []*queryResult{}

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
		Find(&queryResults).Error

	if err != nil {
		return nil, err
	}

	conversationDTOs := make([]*readModel.ConversationDTO, len(queryResults))

	for i, result := range queryResults {
		conversationDTO := &readModel.ConversationDTO{
			ID:        result.ID,
			CreatedAt: result.CreatedAt,
		}

		if result.Type == 1 {
			conversationDTO.Avatar = result.UserAvatar
			conversationDTO.Name = result.UserName
		} else {
			conversationDTO.Avatar = result.Avatar
			conversationDTO.Name = result.Name
		}

		conversationDTOs[i] = conversationDTO
	}

	return conversationDTOs, nil
}

func (r *queriesRepository) GetConversation(id uuid.UUID, userID uuid.UUID) (*readModel.ConversationFullDTO, error) {
	type queryResult struct {
		ID         uuid.UUID
		Name       string
		Avatar     string
		CreatedAt  time.Time
		Type       uint8
		OwnerID    uuid.UUID
		UserID     uuid.UUID
		UserAvatar string
		UserName   string
	}

	conversation := queryResult{}

	err := r.db.Model(&Conversation{}).
		Select(
			"conversations.id", "conversations.created_at", "conversations.type as type",
			"group_conversations.avatar", "group_conversations.name", "group_conversations.owner_id",
			"direct_conversations.from_user_id", "direct_conversations.to_user_id",
			"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
		).
		Joins("LEFT JOIN group_conversations ON group_conversations.conversation_id = conversations.id").
		Joins("LEFT JOIN direct_conversations ON direct_conversations.conversation_id = conversations.id").
		Joins("LEFT JOIN users ON direct_conversations.from_user_id = users.id OR direct_conversations.to_user_id = users.id").
		Where("users.id IS NULL OR users.id <> ?", userID).
		Where("conversations.is_active = ?", true).
		Where("conversations.id = ?", id).
		Find(&conversation).Error

	if err != nil {
		return nil, err
	}

	conversationDTO := &readModel.ConversationFullDTO{
		ID:        conversation.ID,
		CreatedAt: conversation.CreatedAt,
		Type:      conversationTypesMap[conversation.Type].String(),
	}

	if conversation.Type == 1 {
		conversationDTO.Avatar = conversation.UserAvatar
		conversationDTO.Name = conversation.UserName
	} else {
		conversationDTO.Avatar = conversation.Avatar
		conversationDTO.Name = conversation.Name
		conversationDTO.IsOwner = conversation.OwnerID == userID
	}

	hasUserJoined := true

	participant := Participant{}
	if err := r.db.Where("conversation_id = ?", id).Where("user_id = ?", userID).Where("is_active = ?", true).First(&participant).Error; err != nil {
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

	conversationDTO.HasJoined = hasUserJoined
	conversationDTO.ParticipantsCount = participantsCount

	return conversationDTO, nil
}
