package cache

import (
	"encoding/json"

	"GitHub/go-chat/backend/internal/domain"
)

type UserCache struct {
	ID     string
	Name   string
	Avatar string
}

func SerializeUser(user *domain.User) ([]byte, error) {
	cacheUser := UserCache{
		ID:     user.ID.String(),
		Name:   user.Name.String(),
		Avatar: user.Avatar,
	}
	return json.Marshal(cacheUser)
}

func DeserializeUser(data []byte) (*UserCache, error) {
	var cacheUser UserCache
	err := json.Unmarshal(data, &cacheUser)
	if err != nil {
		return nil, err
	}
	return &cacheUser, nil
}

type ConversationCache struct {
	ID       string
	Type     string
	Name     string
	Avatar   string
	IsActive bool
}

func SerializeGroupConversation(conv *domain.GroupConversation) ([]byte, error) {
	cacheConv := ConversationCache{
		ID:       conv.ID.String(),
		Type:     "group",
		Name:     conv.Name.String(),
		Avatar:   conv.Avatar,
		IsActive: conv.IsActive,
	}
	return json.Marshal(cacheConv)
}

func DeserializeConversation(data []byte) (*ConversationCache, error) {
	var cacheConv ConversationCache
	err := json.Unmarshal(data, &cacheConv)
	if err != nil {
		return nil, err
	}
	return &cacheConv, nil
}

type ParticipantCache struct {
	ID             string
	ConversationID string
	UserID         string
	IsActive       bool
}

type ParticipantsCache struct {
	Participants []ParticipantCache
	Count        int
}

func SerializeParticipants(participants []*domain.Participant) ([]byte, error) {
	cacheParticipants := make([]ParticipantCache, len(participants))
	for i, p := range participants {
		cacheParticipants[i] = ParticipantCache{
			ID:             p.ID.String(),
			ConversationID: p.ConversationID.String(),
			UserID:         p.UserID.String(),
			IsActive:       p.IsActive,
		}
	}
	cacheData := ParticipantsCache{
		Participants: cacheParticipants,
		Count:        len(participants),
	}
	return json.Marshal(cacheData)
}

func DeserializeParticipants(data []byte) (*ParticipantsCache, error) {
	var cacheData ParticipantsCache
	err := json.Unmarshal(data, &cacheData)
	if err != nil {
		return nil, err
	}
	return &cacheData, nil
}

type ConversationMetaCache struct {
	Name             string
	Avatar           string
	ParticipantCount int
	IsActive         bool
}

func SerializeConversationMeta(name, avatar string, participantCount int, isActive bool) ([]byte, error) {
	meta := ConversationMetaCache{
		Name:             name,
		Avatar:           avatar,
		ParticipantCount: participantCount,
		IsActive:         isActive,
	}
	return json.Marshal(meta)
}

func DeserializeConversationMeta(data []byte) (*ConversationMetaCache, error) {
	var meta ConversationMetaCache
	err := json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

type UserConversationsListCache []ConversationListItemCache

type ConversationListItemCache struct {
	ID               string
	Name             string
	Avatar           string
	Type             string
	ParticipantCount int
	UnreadCount      int
}

func SerializeConversationList(items []ConversationListItemCache) ([]byte, error) {
	return json.Marshal(items)
}

func DeserializeConversationList(data []byte) (UserConversationsListCache, error) {
	var list UserConversationsListCache
	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
