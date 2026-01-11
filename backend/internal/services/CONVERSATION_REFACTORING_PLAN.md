# Conversation Service Refactoring Plan

## Overview

Refactor the monolithic `conversationService` into focused subservices using the Facade Pattern.

## Current State

- **File:** `internal/services/conversation.go` (358 lines)
- **Problem:** Single service handles 9 methods across 4 concerns with 8 dependencies
- **Issue:** Not all methods use all dependencies - poor separation of concerns

## Concerns Identified

1. **Group Conversation CRUD** - Create, Rename, Delete
2. **Direct Conversation** - StartDirectConversation
3. **Membership** - Join, Leave, Invite, Kick
4. **Messaging** - SendTextMessage

## Solution: Facade Pattern

Create focused subservices, compose them into a facade implementing the existing `ConversationService` interface.

## File Structure After

```
internal/services/
├── conversation.go           # Facade (~50 lines)
├── group_conversation.go     # GroupConversationService (~70 lines)
├── group_conversation_test.go
├── direct_conversation.go    # DirectConversationService (~30 lines)
├── direct_conversation_test.go
├── membership.go             # MembershipService (~100 lines)
├── membership_test.go
├── message.go                # MessageService (~30 lines)
├── message_test.go
├── auth.go
├── cache.go
├── notifications.go
└── redis_broadcaster.go
```

## New Services

### GroupConversationService

```go
type GroupConversationService interface {
    CreateGroupConversation(ctx context.Context, conversationID uuid.UUID, name string, userID uuid.UUID) error
    DeleteGroupConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
    Rename(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, name string) error
}

type groupConversationService struct {
    groupConversations domain.GroupConversationRepository
    queries            readModel.QueriesRepository
    messages           domain.MessageRepository
    notifications      NotificationService
    cache              CacheService
}
```

### DirectConversationService

```go
type DirectConversationService interface {
    StartDirectConversation(ctx context.Context, fromUserID uuid.UUID, toUserID uuid.UUID) (uuid.UUID, error)
}

type directConversationService struct {
    directConversations domain.DirectConversationRepository
    notifications       NotificationService
    cache               CacheService
}
```

### MembershipService

```go
type MembershipService interface {
    Join(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
    Leave(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
    Invite(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, inviteeID uuid.UUID) error
    Kick(ctx context.Context, conversationID uuid.UUID, kickerID uuid.UUID, targetID uuid.UUID) error
}

type membershipService struct {
    participants   domain.ParticipantRepository
    queries        readModel.QueriesRepository
    messages       domain.MessageRepository
    notifications  NotificationService
    cache          CacheService
}
```

### MessageService

```go
type MessageService interface {
    SendTextMessage(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, messageText string) error
}

type messageService struct {
    queries       readModel.QueriesRepository
    notifications NotificationService
}
```

## Facade (conversation.go)

```go
type conversationService struct {
    groupConversation GroupConversationService
    directConversation DirectConversationService
    membership MembershipService
    message MessageService
}

func (s *conversationService) CreateGroupConversation(...) error {
    return s.groupConversation.CreateGroupConversation(...)
}
// ... delegate all methods
```

## Dependency Mapping

| Method | Current Dependencies | New Service |
|--------|---------------------|-------------|
| CreateGroupConversation | gc, cache | GroupConversationService |
| StartDirectConversation | dc, notif, cache | DirectConversationService |
| DeleteGroupConversation | gc, notif, cache | GroupConversationService |
| Rename | gc, msg, queries, notif, cache | GroupConversationService |
| SendTextMessage | queries, notif | MessageService |
| Join | participants, msg, notif, cache | MembershipService |
| Leave | queries, msg, notif, cache | MembershipService |
| Invite | queries, msg, notif, cache | MembershipService |
| Kick | queries, msg, notif, cache | MembershipService |

## Implementation Steps

1. Create `group_conversation.go` with `GroupConversationService`
2. Create `group_conversation_test.go`
3. Create `direct_conversation.go` with `DirectConversationService`
4. Create `direct_conversation_test.go`
5. Create `membership.go` with `MembershipService`
6. Create `membership_test.go`
7. Create `message.go` with `MessageService`
8. Create `message_test.go`
9. Refactor `conversation.go` into facade
10. Update `cmd/server/main.go` wireup
11. Run tests and lint

## Benefits

- **Single Responsibility:** Each service has one reason to change
- **Explicit Dependencies:** Services only inject what they use
- **Testability:** Smaller units, focused mocks
- **Zero Handler Changes:** Facade preserves interface contract
- **Extensibility:** Easy to add new features to specific services

## Naming Convention

- Exported interfaces: `GroupConversationService`, `DirectConversationService`, etc.
- Unexported implementations: `groupConversationService`, `directConversationService`, etc.
- Constructor: `NewGroupConversationService(...)`
