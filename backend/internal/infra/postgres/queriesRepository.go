package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
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

	err := r.db.Scopes(r.paginate(paginationInfo)).Model(&User{}).Not(&User{ID: userID}).Find(&users).Error

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

func (r *queriesRepository) GetPotentialInvitees(conversationID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.ContactDTO, error) {
	users := []*readModel.ContactDTO{}

	subQuery := r.db.Select("user_id").Where("conversation_id = ?", conversationID).Where("is_active = ?", true).Table("participants")
	err := r.db.Scopes(r.paginate(paginationInfo)).Model(&User{}).Where("id NOT IN (?)", subQuery).Find(&users).Error

	return users, err
}

func (r *queriesRepository) GetUserByID(id uuid.UUID) (*readModel.UserDTO, error) {
	user := User{}

	err := r.db.Where(&User{ID: id}).First(&user).Error

	if err != nil {
		return nil, err
	}

	userDTO := &readModel.UserDTO{
		ID:     user.ID,
		Name:   user.Name,
		Avatar: user.Avatar,
	}

	return userDTO, nil
}

func (r *queriesRepository) GetConversationMessages(conversationID uuid.UUID, requestUserID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.MessageDTO, error) {
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

	messages := make([]*readModel.MessageDTO, len(queryResults))

	for i, result := range queryResults {
		messages[i] = toMessageDTO(result, requestUserID)
	}

	return messages, nil
}

func (r *queriesRepository) GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (*readModel.MessageDTO, error) {
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
		return nil, err
	}

	return toMessageDTO(message, requestUserID), nil
}

func (r *queriesRepository) GetUserConversations(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.ConversationDTO, error) {
	conversations := []*Conversation{}

	err := r.db.Scopes(r.paginate(paginationInfo)).
		Model(&Conversation{}).
		Select("conversations.id", "conversations.created_at", "conversations.type").
		Joins("JOIN participants ON participants.conversation_id = conversations.id").
		Where("participants.user_id = ? ", userID).
		Where(&Conversation{IsActive: true}).
		Where("participants.is_active = ?", true).
		Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	conversationDTOs := make([]*readModel.ConversationDTO, len(conversations))

	for i, conversation := range conversations {
		conversationDTO := &readModel.ConversationDTO{
			ID:   conversation.ID,
			Type: conversationTypesMap[conversation.Type].String(),
		}

		message := &messageQuery{}

		err := r.db.Model(&Message{}).
			Select(
				"messages.id", "messages.type as type", "messages.created_at", "messages.content", "messages.conversation_id",
				"users.id as user_id", "users.name as user_name", "users.avatar as user_avatar",
			).
			Joins("LEFT JOIN users ON messages.user_id = users.id").
			Where(&Message{ConversationID: conversation.ID}).
			Order("messages.created_at desc").
			Limit(1).
			First(&message).Error

		if err == nil {
			conversationDTO.LastMessage = *toMessageDTO(message, userID)
		}

		switch conversationTypesMap[conversation.Type] {
		case domain.ConversationTypeDirect:

			type userQuery struct {
				UserID     uuid.UUID
				UserAvatar string
				UserName   string
			}

			userQueryResult := &userQuery{}

			err := r.db.Model(&Participant{}).
				Select("users.id as user_id", "users.name as user_name", "users.avatar as user_avatar").
				Joins("JOIN users ON participants.user_id = users.id").
				Where("users.id <> ?", userID).
				Where(&Participant{ConversationID: conversation.ID, IsActive: true}).
				First(&userQueryResult).Error

			if err != nil {
				return nil, err
			}

			conversationDTO.Avatar = userQueryResult.UserAvatar
			conversationDTO.Name = userQueryResult.UserName
		case domain.ConversationTypeGroup:
			groupConversationQueryResult := &GroupConversation{}

			err := r.db.Where(&GroupConversation{ConversationID: conversation.ID}).First(&groupConversationQueryResult).Error

			if err != nil {
				return nil, err
			}

			conversationDTO.Avatar = groupConversationQueryResult.Avatar
			conversationDTO.Name = groupConversationQueryResult.Name
		}

		conversationDTOs[i] = conversationDTO
	}

	return conversationDTOs, nil
}

func (r *queriesRepository) GetConversation(id uuid.UUID, userID uuid.UUID) (*readModel.ConversationFullDTO, error) {
	conversation := Conversation{}

	err := r.db.Where(&Conversation{ID: id, IsActive: true}).First(&conversation).Error

	if err != nil {
		return nil, err
	}

	conversationDTO := &readModel.ConversationFullDTO{
		ID:        conversation.ID,
		CreatedAt: conversation.CreatedAt,
		Type:      conversationTypesMap[conversation.Type].String(),
	}

	switch conversationTypesMap[conversation.Type] {
	case domain.ConversationTypeDirect:
		type userQuery struct {
			UserID     uuid.UUID
			UserAvatar string
			UserName   string
		}

		userQueryResult := &userQuery{}

		err := r.db.Model(&Participant{}).
			Select("users.id as user_id", "users.name as user_name", "users.avatar as user_avatar").
			Joins("JOIN users ON participants.user_id = users.id").
			Where("users.id <> ?", userID).
			Where(&Participant{ConversationID: conversation.ID, IsActive: true}).
			First(&userQueryResult).Error

		if err != nil {
			return nil, err
		}

		conversationDTO.Avatar = userQueryResult.UserAvatar
		conversationDTO.Name = userQueryResult.UserName
	case domain.ConversationTypeGroup:
		groupConversationQueryResult := &GroupConversation{}

		err := r.db.Where(&GroupConversation{ConversationID: groupConversationQueryResult.ConversationID}).First(&groupConversationQueryResult).Error

		if err != nil {
			return nil, err
		}

		conversationDTO.Avatar = groupConversationQueryResult.Avatar
		conversationDTO.Name = groupConversationQueryResult.Name
		conversationDTO.IsOwner = groupConversationQueryResult.OwnerID == userID

		var participantsCount int64

		err = r.db.Model(&Participant{}).Where(&Participant{ConversationID: conversation.ID, IsActive: true}).Count(&participantsCount).Error

		if err != nil {
			return nil, err
		}
		conversationDTO.ParticipantsCount = participantsCount

		hasUserJoined := true

		participant := Participant{}
		if err := r.db.Where(&Participant{ConversationID: groupConversationQueryResult.ConversationID, IsActive: true, UserID: userID}).First(&participant).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				hasUserJoined = false
			} else {
				return nil, err
			}
		}

		conversationDTO.HasJoined = hasUserJoined
	}

	return conversationDTO, nil
}
