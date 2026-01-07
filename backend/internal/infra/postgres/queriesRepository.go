package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
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

func (r *queriesRepository) GetContacts(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	users := []*User{}

	err := r.db.Scopes(r.paginate(paginationInfo)).Model(&User{}).Not(&User{ID: userID}).Find(&users).Error

	if err != nil {
		return nil, err
	}

	usersDTO := make([]readModel.ContactDTO, len(users))

	for i, user := range users {
		usersDTO[i] = readModel.ContactDTO{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		}
	}

	return usersDTO, nil
}

func (r *queriesRepository) GetParticipants(conversationID uuid.UUID, userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	users := []*User{}

	err := r.db.Scopes(r.paginate(paginationInfo)).
		Model(&User{}).
		Joins("LEFT JOIN participants on participants.user_id = users.id").
		Where("participants.conversation_id = ?", conversationID).
		Where("participants.is_active = ?", true).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	usersDTO := make([]readModel.ContactDTO, len(users))

	for i, user := range users {
		usersDTO[i] = readModel.ContactDTO{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		}
	}

	return usersDTO, nil
}

func (r *queriesRepository) GetPotentialInvitees(conversationID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	users := []*User{}

	subQuery := r.db.Select("user_id").Where("conversation_id = ?", conversationID).Where("is_active = ?", true).Table("participants")
	err := r.db.Scopes(r.paginate(paginationInfo)).Model(&User{}).Where("id NOT IN (?)", subQuery).Find(&users).Error

	if err != nil {
		return nil, err
	}

	usersDTO := make([]readModel.ContactDTO, len(users))

	for i, user := range users {
		usersDTO[i] = readModel.ContactDTO{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		}
	}

	return usersDTO, nil
}

func (r *queriesRepository) GetUserByID(id uuid.UUID) (readModel.UserDTO, error) {
	user := User{}

	err := r.db.Where(&User{ID: id}).First(&user).Error

	if err != nil {
		return readModel.UserDTO{}, err
	}

	userDTO := readModel.UserDTO{
		ID:     user.ID,
		Name:   user.Name,
		Avatar: user.Avatar,
	}

	return userDTO, nil
}

func (r *queriesRepository) GetConversationMessages(conversationID uuid.UUID, requestUserID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.MessageDTO, error) {
	queryResults := []*messageQuery{}

	err := r.db.Scopes(r.paginate(paginationInfo)).
		Model(&Message{}).
		Select(
			"messages.id", "messages.type as type", "messages.created_at", "messages.conversation_id", "messages.content",
			"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
		).
		Joins("LEFT JOIN users ON messages.user_id = users.id").
		Where(&Message{ConversationID: conversationID}).
		Order("messages.created_at asc").
		Find(&queryResults).Error

	if err != nil {
		return nil, err
	}

	messages := make([]readModel.MessageDTO, len(queryResults))

	for i, result := range queryResults {
		messages[i] = toMessageDTO(*result, requestUserID)
	}

	return messages, nil
}

func (r *queriesRepository) GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (readModel.MessageDTO, error) {
	message := &messageQuery{}

	err := r.db.Model(&Message{}).
		Select(
			"messages.id", "messages.type as type", "messages.created_at", "messages.content", "messages.conversation_id",
			"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
		).
		Joins("LEFT JOIN users ON messages.user_id = users.id").
		Where(&Message{ID: messageID}).
		Find(&message).Error

	if err != nil {
		return readModel.MessageDTO{}, err
	}

	return toMessageDTO(*message, requestUserID), nil
}

func (r *queriesRepository) GetUserConversations(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ConversationDTO, error) {
	type conversationQuery struct {
		ConversationID    uuid.UUID
		CreatedAt         string
		Type              uint8
		MessageID         *uuid.UUID
		MessageType       *uint8
		MessageContent    *string
		MessageCreatedAt  *string
		MessageUserID     *uuid.UUID
		MessageUserName   *string
		MessageUserAvatar *string
		GroupAvatar       *string
		GroupName         *string
		OtherUserID       *uuid.UUID
		OtherUserName     *string
		OtherUserAvatar   *string
	}

	queryResults := []*conversationQuery{}

	lastMessagesSubquery := r.db.Table("messages").
		Select("conversation_id, MAX(created_at) as max_created_at").
		Group("conversation_id")

	err := r.db.Scopes(r.paginate(paginationInfo)).
		Table("conversations").
		Select(`
			conversations.id as conversation_id,
			conversations.created_at,
			conversations.type,
			messages.id as message_id,
			messages.type as message_type,
			messages.content as message_content,
			messages.created_at as message_created_at,
			messages.user_id as message_user_id,
			msg_users.name as message_user_name,
			msg_users.avatar as message_user_avatar,
			group_conversations.avatar as group_avatar,
			group_conversations.name as group_name,
			other_users.id as other_user_id,
			other_users.name as other_user_name,
			other_users.avatar as other_user_avatar
		`).
		Joins("JOIN participants ON participants.conversation_id = conversations.id").
		Joins("LEFT JOIN (?) AS last_messages ON last_messages.conversation_id = conversations.id", lastMessagesSubquery).
		Joins("LEFT JOIN messages ON messages.conversation_id = conversations.id AND messages.created_at = last_messages.max_created_at").
		Joins("LEFT JOIN users msg_users ON msg_users.id = messages.user_id").
		Joins("LEFT JOIN group_conversations ON group_conversations.conversation_id = conversations.id").
		Joins(`LEFT JOIN participants other_participants 
			ON other_participants.conversation_id = conversations.id 
			AND other_participants.user_id <> ? 
			AND other_participants.is_active = ?`, userID, true).
		Joins("LEFT JOIN users other_users ON other_users.id = other_participants.user_id").
		Where("participants.user_id = ? ", userID).
		Where(&Conversation{IsActive: true}).
		Where("participants.is_active = ?", true).
		Find(&queryResults).Error

	if err != nil {
		return nil, err
	}

	conversationDTOs := make([]readModel.ConversationDTO, len(queryResults))

	for i, result := range queryResults {
		conversationDTO := readModel.ConversationDTO{
			ID:   result.ConversationID,
			Type: conversationTypesMap[result.Type].String(),
		}

		if result.MessageID != nil {
			createdAt, _ := time.Parse(time.RFC3339, *result.MessageCreatedAt)

			text := ""
			switch messageTypesMap[*result.MessageType] {
			case domain.MessageTypeText:
				text = *result.MessageContent
			case domain.MessageTypeRenamedConversation:
				text = *result.MessageUserName + " renamed chat to " + *result.MessageContent
			case domain.MessageTypeJoinedConversation:
				text = *result.MessageUserName + " joined"
			case domain.MessageTypeLeftConversation:
				text = *result.MessageUserName + " left"
			case domain.MessageTypeInvitedConversation:
				text = *result.MessageUserName + " was invited"
			}

			messageDTO := readModel.MessageDTO{
				ID:             *result.MessageID,
				CreatedAt:      createdAt,
				Text:           text,
				Type:           messageTypesMap[*result.MessageType].String(),
				ConversationId: result.ConversationID,
				User: readModel.UserDTO{
					ID:     *result.MessageUserID,
					Avatar: *result.MessageUserAvatar,
					Name:   *result.MessageUserName,
				},
			}

			if messageTypesMap[*result.MessageType] == domain.MessageTypeText {
				messageDTO.IsInbound = *result.MessageUserID != userID
			}

			conversationDTO.LastMessage = messageDTO
		}

		switch conversationTypesMap[result.Type] {
		case domain.ConversationTypeDirect:
			if result.OtherUserID != nil {
				conversationDTO.Avatar = *result.OtherUserAvatar
				conversationDTO.Name = *result.OtherUserName
			}
		case domain.ConversationTypeGroup:
			if result.GroupAvatar != nil && result.GroupName != nil {
				conversationDTO.Avatar = *result.GroupAvatar
				conversationDTO.Name = *result.GroupName
			}
		}

		conversationDTOs[i] = conversationDTO
	}

	return conversationDTOs, nil
}

func (r *queriesRepository) GetConversation(id uuid.UUID, userID uuid.UUID) (readModel.ConversationFullDTO, error) {
	type conversationQuery struct {
		ConversationID    uuid.UUID
		CreatedAt         time.Time
		Type              uint8
		OtherUserID       *uuid.UUID
		OtherUserName     *string
		OtherUserAvatar   *string
		GroupAvatar       *string
		GroupName         *string
		GroupOwnerID      *uuid.UUID
		ParticipantsCount *int64
		UserParticipantID *uuid.UUID
	}

	queryResult := &conversationQuery{}

	participantsCountSubquery := r.db.Table("participants").
		Select("COUNT(*)").
		Where("participants.conversation_id = conversations.id").
		Where("participants.is_active = ?", true)

	err := r.db.Table("conversations").
		Select(`
			conversations.id as conversation_id,
			conversations.created_at,
			conversations.type,
			other_users.id as other_user_id,
			other_users.name as other_user_name,
			other_users.avatar as other_user_avatar,
			group_conversations.avatar as group_avatar,
			group_conversations.name as group_name,
			group_conversations.owner_id as group_owner_id,
			(?) as participants_count,
			user_participants.id as user_participant_id
		`, participantsCountSubquery).
		Joins(`LEFT JOIN participants other_participants 
			ON other_participants.conversation_id = conversations.id 
			AND other_participants.user_id <> ? 
			AND other_participants.is_active = ?`, userID, true).
		Joins("LEFT JOIN users other_users ON other_users.id = other_participants.user_id").
		Joins("LEFT JOIN group_conversations ON group_conversations.conversation_id = conversations.id").
		Joins(`LEFT JOIN participants user_participants 
			ON user_participants.conversation_id = conversations.id 
			AND user_participants.user_id = ? 
			AND user_participants.is_active = ?`, userID, true).
		Where(&Conversation{ID: id, IsActive: true}).
		First(&queryResult).Error

	if err != nil {
		return readModel.ConversationFullDTO{}, err
	}

	conversationDTO := readModel.ConversationFullDTO{
		ID:        queryResult.ConversationID,
		CreatedAt: queryResult.CreatedAt,
		Type:      conversationTypesMap[queryResult.Type].String(),
	}

	switch conversationTypesMap[queryResult.Type] {
	case domain.ConversationTypeDirect:
		if queryResult.OtherUserID != nil {
			conversationDTO.Avatar = *queryResult.OtherUserAvatar
			conversationDTO.Name = *queryResult.OtherUserName
		}
	case domain.ConversationTypeGroup:
		if queryResult.GroupAvatar != nil && queryResult.GroupName != nil {
			conversationDTO.Avatar = *queryResult.GroupAvatar
			conversationDTO.Name = *queryResult.GroupName
		}
		if queryResult.GroupOwnerID != nil {
			conversationDTO.IsOwner = *queryResult.GroupOwnerID == userID
		}
		if queryResult.ParticipantsCount != nil {
			conversationDTO.ParticipantsCount = *queryResult.ParticipantsCount
		}
		conversationDTO.HasJoined = queryResult.UserParticipantID != nil
	}

	return conversationDTO, nil
}
