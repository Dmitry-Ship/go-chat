package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/readModel"
)

var conversationTypesMap = map[uint8]string{
	0: "public",
	1: "private",
}

func toConversationTypePersistence(conversationType string) uint8 {
	for k, v := range conversationTypesMap {
		if v == conversationType {
			return k
		}
	}

	return 0
}

func toConversationDTO(conversation *Conversation) *readModel.ConversationDTO {
	dto := readModel.ConversationDTO{
		ID:        conversation.ID,
		CreatedAt: conversation.CreatedAt,
		Type:      conversationTypesMap[conversation.Type],
	}

	return &dto
}

func toPublicConversationDTO(conversation *Conversation, avatar string, name string) *readModel.ConversationDTO {
	dto := toConversationDTO(conversation)

	dto.Avatar = avatar
	dto.Name = name

	return dto
}

func toPrivateConversationDTO(conversation *Conversation, user *User) *readModel.ConversationDTO {
	dto := toConversationDTO(conversation)

	dto.Avatar = user.Avatar
	dto.Name = user.Name

	return dto
}

func toConversationFullDTO(conversation *Conversation, avatar string, name string, hasJoined bool) *readModel.ConversationFullDTO {
	return &readModel.ConversationFullDTO{
		Conversation: *toPublicConversationDTO(conversation, avatar, name),
		HasJoined:    hasJoined,
	}
}

func toConversationPersistence(conversation domain.BaseConversation) *Conversation {
	conversationBase := conversation.GetBaseData()
	return &Conversation{
		ID:        conversationBase.ID,
		Type:      toConversationTypePersistence(conversationBase.Type),
		CreatedAt: conversationBase.CreatedAt,
	}
}

func toPublicConversationPersistence(conversation *domain.PublicConversation) *PublicConversation {
	return &PublicConversation{
		ID:             conversation.Data.ID,
		ConversationID: conversation.ID,
		Name:           conversation.Data.Name,
		Avatar:         conversation.Data.Avatar,
	}
}

func toPublicConversationDomain(conversation *Conversation, publicConversation *PublicConversation) *domain.PublicConversation {
	return &domain.PublicConversation{
		Data: domain.PublicConversationData{
			ID:     publicConversation.ID,
			Name:   publicConversation.Name,
			Avatar: publicConversation.Avatar,
		},
		Conversation: domain.Conversation{
			ID:        conversation.ID,
			Type:      conversationTypesMap[conversation.Type],
			CreatedAt: conversation.CreatedAt,
		},
	}
}

func toPrivateConversationPersistence(conversation *domain.PrivateConversation) *PrivateConversation {
	return &PrivateConversation{
		ID:             conversation.Data.ID,
		ConversationID: conversation.ID,
		FromUserID:     conversation.Data.FromUser.UserID,
		ToUserID:       conversation.Data.ToUser.UserID,
	}
}
