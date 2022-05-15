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
	users := []*User{}
	err := r.db.Scopes(r.paginate(paginationInfo)).Where("id <> ?", userID).Find(&users).Error

	dtoContacts := make([]*readModel.ContactDTO, len(users))

	for i, user := range users {
		dtoContacts[i] = toContactDTO(user)
	}

	return dtoContacts, err
}

func (r *queriesRepository) GetConversationMessages(conversationID uuid.UUID, requestUserID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.MessageDTO, error) {
	messages := []*Message{}

	err := r.db.Scopes(r.paginate(paginationInfo)).Order("created_at desc").Where("conversation_id = ?", conversationID).Find(&messages).Error

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

func (r *queriesRepository) GetUserConversations(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.ConversationDTO, error) {
	conversations := []*Conversation{}

	err := r.db.Scopes(r.paginate(paginationInfo)).Joins("JOIN participants ON participants.conversation_id = conversations.id").Where("conversations.is_active = ?", true).Where("participants.is_active = ?", true).Where("participants.user_id = ?", userID).Find(&conversations).Error

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

func (r *queriesRepository) GetPotentialInvitees(conversationID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]*readModel.ContactDTO, error) {
	users := []*User{}
	subQuery := r.db.Select("user_id").Where("conversation_id = ?", conversationID).Table("participants")
	err := r.db.Scopes(r.paginate(paginationInfo)).Where("id NOT IN (?)", subQuery).Find(&users).Error

	dtoContacts := make([]*readModel.ContactDTO, len(users))

	for i, user := range users {
		dtoContacts[i] = toContactDTO(user)
	}

	return dtoContacts, err
}

func (r *queriesRepository) GetUserByID(id uuid.UUID) (*readModel.UserDTO, error) {
	user := User{}
	err := r.db.Where("id = ?", id).First(&user).Error

	return toUserDTO(&user), err
}

func (r *queriesRepository) GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (*readModel.MessageDTO, error) {
	message := Message{}

	err := r.db.Where("id = ?", messageID).First(&message).Error

	if err != nil {
		return nil, err
	}

	return r.getMessageDTO(&message, requestUserID)
}

func (r *queriesRepository) getMessageDTO(message *Message, requestUserID uuid.UUID) (*readModel.MessageDTO, error) {
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
	case 4:
		return toMessageDTO(message, &user), nil
	}

	return nil, nil
}

func (r *queriesRepository) GetConversation(id uuid.UUID, userId uuid.UUID) (*readModel.ConversationFullDTO, error) {
	conversation := Conversation{}

	err := r.db.Where("id = ?", id).Where("is_active = ?", true).First(&conversation).Error

	if err != nil {
		return nil, err
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
