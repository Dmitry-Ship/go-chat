package cache

import (
	"context"
	"log"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra"
)

type CacheInvalidationService struct {
	cache    CacheClient
	eventBus infra.EventsSubscriber
}

func NewCacheInvalidationService(cache CacheClient, eventBus infra.EventsSubscriber) *CacheInvalidationService {
	return &CacheInvalidationService{
		cache:    cache,
		eventBus: eventBus,
	}
}

func (s *CacheInvalidationService) Run(ctx context.Context) {
	subscription := s.eventBus.Subscribe(domain.DomainEventTopic)

	for {
		select {
		case event := <-subscription:
			if domainEvent, ok := event.Data.(domain.DomainEvent); ok {
				s.handleDomainEvent(ctx, domainEvent)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *CacheInvalidationService) handleDomainEvent(ctx context.Context, event domain.DomainEvent) {
	switch e := event.(type) {
	case domain.GroupConversationRenamed:
		s.invalidateConversation(e.GetConversationID().String())
		log.Printf("🗑️ Cache invalidated for renamed conversation: %s", e.GetConversationID())

	case domain.GroupConversationDeleted:
		s.invalidateConversation(e.GetConversationID().String())
		log.Printf("🗑️ Cache invalidated for deleted conversation: %s", e.GetConversationID())

	case domain.GroupConversationLeft:
		s.invalidateParticipants(e.GetConversationID().String())
		log.Printf("🗑️ Participants cache invalidated for conversation: %s", e.GetConversationID())

	case domain.GroupConversationJoined:
		s.invalidateParticipants(e.GetConversationID().String())
		log.Printf("🗑️ Participants cache invalidated for conversation: %s", e.GetConversationID())

	case domain.GroupConversationInvited:
		s.invalidateParticipants(e.GetConversationID().String())
		log.Printf("🗑️ Participants cache invalidated for conversation: %s", e.GetConversationID())

	case domain.DirectConversationCreated:
		for _, userID := range e.UserIDs {
			s.invalidateUserConversations(userID.String())
		}
		log.Printf("🗑️ User conversation lists invalidated")

	case domain.GroupConversationCreated:
		s.invalidateUserConversations(e.OwnerID.String())
		log.Printf("🗑️ User conversation list invalidated for user: %s", e.OwnerID)
	}
}

func (s *CacheInvalidationService) invalidateConversation(conversationID string) {
	_ = s.cache.Delete(context.Background(), ConversationKey(conversationID))
	_ = s.cache.Delete(context.Background(), ConvMetaKey(conversationID))
	_ = s.cache.DeletePattern(context.Background(), ParticipantsKey(conversationID))
}

func (s *CacheInvalidationService) invalidateParticipants(conversationID string) {
	_ = s.cache.Delete(context.Background(), ParticipantsKey(conversationID))
	_ = s.cache.Delete(context.Background(), ConvMetaKey(conversationID))
}

func (s *CacheInvalidationService) invalidateUserConversations(userID string) {
	_ = s.cache.Delete(context.Background(), UserConvListKey(userID))
}
