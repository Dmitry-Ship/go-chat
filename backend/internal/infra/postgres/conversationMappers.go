package postgres

import (
	"GitHub/go-chat/backend/internal/domain"
)

var conversationTypesMap = map[uint8]domain.ConversationType{
	0: domain.ConversationTypeGroup,
	1: domain.ConversationTypeDirect,
}

func toConversationTypePersistence(conversationType domain.ConversationType) uint8 {
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
		ID:       conversationBase.ID,
		Type:     toConversationTypePersistence(conversationBase.Type),
		IsActive: conversationBase.IsActive,
	}
}

func toGroupConversationPersistence(conversation *domain.GroupConversation) *GroupConversation {
	return &GroupConversation{
		ID:             conversation.ID,
		ConversationID: conversation.Conversation.ID,
		Name:           conversation.Name.String(),
		Avatar:         conversation.Avatar,
		OwnerID:        conversation.Owner.UserID,
	}
}

func toGroupConversationDomain(conversation *Conversation, groupConversation *GroupConversation, participant *Participant) *domain.GroupConversation {
	name, err := domain.NewConversationName(groupConversation.Name)

	if err != nil {
		name, _ = domain.NewConversationName("Name Corrupted")
	}

	return &domain.GroupConversation{
		ID:     groupConversation.ID,
		Name:   *name,
		Avatar: groupConversation.Avatar,
		Owner: domain.Participant{
			UserID:         participant.UserID,
			ID:             participant.ID,
			ConversationID: groupConversation.ConversationID,
		},
		Conversation: domain.Conversation{
			ID:       conversation.ID,
			Type:     conversationTypesMap[conversation.Type],
			IsActive: conversation.IsActive,
		},
	}
}

func toDirectConversationDomain(conversation *Conversation, directConversation *DirectConversation, toUser *Participant, fromUser *Participant) *domain.DirectConversation {
	return &domain.DirectConversation{
		ID: directConversation.ID,
		ToUser: domain.Participant{
			UserID:         toUser.UserID,
			ID:             toUser.ID,
			ConversationID: directConversation.ConversationID,
		},
		FromUser: domain.Participant{
			UserID:         fromUser.UserID,
			ID:             fromUser.ID,
			ConversationID: directConversation.ConversationID,
		},
		Conversation: domain.Conversation{
			ID:       conversation.ID,
			Type:     conversationTypesMap[conversation.Type],
			IsActive: conversation.IsActive,
		},
	}
}

func toDirectConversationPersistence(conversation *domain.DirectConversation) *DirectConversation {
	return &DirectConversation{
		ID:             conversation.ID,
		ConversationID: conversation.Conversation.ID,
		FromUserID:     conversation.FromUser.UserID,
		ToUserID:       conversation.ToUser.UserID,
	}
}
