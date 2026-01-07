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
		Name:           conversation.Name,
		Avatar:         conversation.Avatar,
		OwnerID:        conversation.Owner.UserID,
	}
}

func toGroupConversationDomain(conversation Conversation, groupConversation GroupConversation, participant Participant) *domain.GroupConversation {
	return &domain.GroupConversation{
		ID:     groupConversation.ID,
		Name:   groupConversation.Name,
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

func toDirectConversationDomain(conversation *Conversation, participants []*Participant) *domain.DirectConversation {
	participantsDomain := make([]domain.Participant, len(participants))

	for i, participant := range participants {
		participantsDomain[i] = domain.Participant{
			UserID:         participant.UserID,
			ID:             participant.ID,
			ConversationID: participant.ConversationID,
		}
	}

	return &domain.DirectConversation{
		Participants: participantsDomain,
		Conversation: domain.Conversation{
			ID:       conversation.ID,
			Type:     conversationTypesMap[conversation.Type],
			IsActive: conversation.IsActive,
		},
	}
}
