package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
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

func toConversationPersistence(conversation domain.BaseConversation) *Conversation {
	conversationBase := conversation.GetBaseData()
	return &Conversation{
		ID:        conversationBase.ID,
		Type:      toConversationTypePersistence(conversationBase.Type),
		CreatedAt: conversationBase.CreatedAt,
		IsActive:  conversationBase.IsActive,
	}
}

func toPublicConversationPersistence(conversation *domain.PublicConversation) *PublicConversation {
	return &PublicConversation{
		ID:             conversation.Data.ID,
		ConversationID: conversation.ID,
		Name:           conversation.Data.Name,
		Avatar:         conversation.Data.Avatar,
		OwnerID:        conversation.Data.Owner.UserID,
	}
}

func toPublicConversationDomain(conversation *Conversation, publicConversation *PublicConversation, participant *Participant) *domain.PublicConversation {
	return &domain.PublicConversation{
		Data: domain.PublicConversationData{
			ID:     publicConversation.ID,
			Name:   publicConversation.Name,
			Avatar: publicConversation.Avatar,
			Owner: domain.Participant{
				UserID:         participant.UserID,
				ID:             participant.ID,
				ConversationID: publicConversation.ConversationID,
				CreatedAt:      participant.CreatedAt,
			},
		},
		Conversation: domain.Conversation{
			ID:        conversation.ID,
			Type:      conversationTypesMap[conversation.Type],
			CreatedAt: conversation.CreatedAt,
			IsActive:  conversation.IsActive,
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
