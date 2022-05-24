package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
)

var conversationTypesMap = map[uint8]string{
	0: "group",
	1: "direct",
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

func toGroupConversationPersistence(conversation *domain.GroupConversation) *GroupConversation {
	return &GroupConversation{
		ID:             conversation.Data.ID,
		ConversationID: conversation.ID,
		Name:           conversation.Data.Name,
		Avatar:         conversation.Data.Avatar,
		OwnerID:        conversation.Data.Owner.UserID,
	}
}

func toGroupConversationDomain(conversation *Conversation, groupConversation *GroupConversation, participant *Participant) *domain.GroupConversation {
	return &domain.GroupConversation{
		Data: domain.GroupConversationData{
			ID:     groupConversation.ID,
			Name:   groupConversation.Name,
			Avatar: groupConversation.Avatar,
			Owner: domain.Participant{
				UserID:         participant.UserID,
				ID:             participant.ID,
				ConversationID: groupConversation.ConversationID,
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

func toDirectConversationDomain(conversation *Conversation, directConversation *DirectConversation, toUser *Participant, fromUser *Participant) *domain.DirectConversation {
	return &domain.DirectConversation{
		Data: domain.DirectConversationData{
			ID: directConversation.ID,
			ToUser: domain.Participant{
				UserID:         toUser.UserID,
				ID:             toUser.ID,
				ConversationID: directConversation.ConversationID,
				CreatedAt:      toUser.CreatedAt,
			},
			FromUser: domain.Participant{
				UserID:         fromUser.UserID,
				ID:             fromUser.ID,
				ConversationID: directConversation.ConversationID,
				CreatedAt:      fromUser.CreatedAt,
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

func toDirectConversationPersistence(conversation *domain.DirectConversation) *DirectConversation {
	return &DirectConversation{
		ID:             conversation.Data.ID,
		ConversationID: conversation.ID,
		FromUserID:     conversation.Data.FromUser.UserID,
		ToUserID:       conversation.Data.ToUser.UserID,
	}
}
