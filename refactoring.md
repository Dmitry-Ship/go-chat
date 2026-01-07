Comprehensive Refactoring Plan
Overview
This plan simplifies the overengineered architecture while maintaining horizontal scaling capabilities and test coverage.
---
Phase 1: Notification System Consolidation (Highest Impact)
Current State
- 4 separate services with channel-based pipeline (100 workers fan-out/fan-in)
- Unnecessary abstraction layers for simple notification flow
Refactor
Files to modify:
- backend/internal/services/notifications.go - consolidate all logic
- backend/internal/services/notificationsPipeline.go - DELETE
- backend/internal/services/notificationsResolver.go - DELETE
- backend/internal/services/notificationsBuilder.go - DELETE
- backend/internal/server/events.go - simplify sendWSNotification
- backend/internal/infra/concurrency.go - DELETE (if only used for notifications)
Changes:
// New consolidated NotificationService
type NotificationService struct {
  ctx            context.Context
  activeClients  ws.ActiveClients
  redisClient    *redis.Client
  participants   domain.ParticipantRepository
  subscriptionSync ws.SubscriptionSync
  queries        readModel.QueriesRepository
}
// Direct event handling without channels
func (s *NotificationService) Notify(event domain.DomainEvent) error {
  // 1. Get recipients directly
  recipients := s.getRecipients(event)
  
  // 2. Build and send messages directly (no pipelines)
  for _, recipient := range recipients {
    message, err := s.buildMessage(recipient.UserID, event)
    if err != nil {
      continue
    }
    s.sendDirectly(recipient, message)
  }
  return nil
}
Tests to update:
- Create new tests for consolidated NotificationService
- Remove tests for deleted services
---
Phase 2: Value Objects Simplification (High Impact)
Current State
- Private structs (userName, userPassword, conversationName)
- Constructor functions with validation
- String() methods
- Password comparison requires injection of compare function
Refactor
Files to modify:
- backend/internal/domain/user.go
- backend/internal/domain/groupConversation.go
- backend/internal/infra/postgres/usersMappers.go
- backend/internal/infra/cache/userCacheDecorator.go
- backend/internal/services/auth.go
Changes:
// Replace value objects with validation functions
func ValidateUsername(username string) error {
  if username == "" {
    return errors.New("username is empty")
  }
  if len(username) > 100 {
    return errors.New("username too long")
  }
  return nil
}
// Direct bcrypt usage
func HashPassword(password string) (string, error) {
  if len(password) < 8 {
    return "", errors.New("password too short")
  }
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}
func ComparePassword(hashed, plain string) error {
  return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
// Updated User struct
type User struct {
  aggregate
  ID           uuid.UUID
  Avatar       string
  Name         string  // Direct string
  PasswordHash string  // Direct string
  RefreshToken string
}
Tests to update:
- backend/internal/domain/user_test.go
- backend/internal/domain/groupConversation_test.go
- backend/internal/services/conversation_test.go
- Update all test helpers that create value objects
---
Phase 3: Simplified Caching Layer (Medium Impact)
Current State
- Separate decorator files per entity
- Repetitive serialization/deserialization logic
- Boilerplate get/set/invalidate patterns
Refactor
Files to modify:
- backend/internal/infra/cache/userCacheDecorator.go - MERGE
- backend/internal/infra/cache/groupConversationCacheDecorator.go - MERGE
- backend/internal/infra/cache/participantCacheDecorator.go - MERGE
- backend/internal/infra/cache/repositoryCache.go - CREATE (new simplified helper)
New simplified approach:
// Single generic caching helper
type CacheHelper struct {
  cache CacheClient
}
// Generic get-or-load pattern
func (h *CacheHelper) GetOrLoad[T any](
  ctx context.Context,
  key string,
  loadFunc func() (*T, error),
  ttl time.Duration,
  serialize func(*T) ([]byte, error),
  deserialize func([]byte) (*T, error),
) (*T, error) {
  data, err := h.cache.Get(ctx, key)
  if err != nil && data != nil {
    return nil, err
  }
  
  if data != nil {
    return deserialize(data)
  }
  
  value, err := loadFunc()
  if err != nil {
    return nil, err
  }
  
  serialized, err := serialize(value)
  if err != nil {
    return nil, err
  }
  
  _ = h.cache.Set(ctx, key, serialized, ttl)
  return value, nil
}
// Usage in repositories - simpler inline caching
func (r *userRepository) GetByID(id uuid.UUID) (*domain.User, error) {
  user, err := r.cache.GetOrLoad(
    context.Background(),
    UserKey(id.String()),
    func() (*domain.User, error) {
      // Load from DB
    },
    TTLUser,
    SerializeUser,
    DeserializeUser,
  )
  return user, err
}
Tests to update:
- backend/internal/infra/cache/userCacheDecorator_test.go - rewrite for new pattern
- Create tests for CacheHelper
---
Phase 4: Domain Events Simplification (Medium Impact)
Current State
- Base structs (domainEvent, conversationEvent) embedded everywhere
- Constructor functions for each event
- Unnecessary interface hierarchy
Refactor
Files to modify:
- backend/internal/domain/events.go
Changes:
// Simplified to plain structs
type DomainEvent interface {
  GetName() string
  GetTopic() string
}
// Remove base structs, use composition or direct fields
type MessageSent struct {
  ConversationID uuid.UUID
  MessageID      uuid.UUID
  UserID         uuid.UUID
}
func (e MessageSent) GetName() string { return MessageSentEventName }
func (e MessageSent) GetTopic() string { return DomainEventTopic }
func (e MessageSent) GetConversationID() uuid.UUID { return e.ConversationID }
// Direct struct literals instead of constructors
event := domain.MessageSent{
  ConversationID: convID,
  MessageID:      msgID,
  UserID:         userID,
}
Tests to update:
- All domain tests that use new*Event() constructors
- Update helper functions
---
Phase 5: Event Bus Simplification (Low Impact)
Current State
- Publisher/Subscriber interfaces
- Only one implementation exists
Refactor
Files to modify:
- backend/internal/infra/eventBus.go
- All files using EventPublisher/EventsSubscriber interfaces
Changes:
// Remove interfaces, use concrete type directly
type EventBus struct {
  mu                  sync.RWMutex
  topicSubscribersMap map[string][]chan Event
  isClosed            bool
}
// Keep same methods, remove interface definitions
// Update all callers to use *EventBus instead of interfaces
Tests to update:
- No significant test changes needed
---
Phase 6: Repository Pattern Simplification (Low Impact)
Current State
- GenericRepository[T] interface
- Base repository with transaction handling
Refactor
Files to modify:
- backend/internal/domain/repository.go
- backend/internal/infra/postgres/repository.go
Changes:
// Remove generic repository interface
// Keep concrete base repository with transaction helpers
type Repository struct {
  db             *gorm.DB
  eventPublisher *EventBus
}
// Simplify transaction handling
func (r *Repository) WithTx(txFunc func(*gorm.DB) error) error {
  tx := r.db.Begin()
  defer func() {
    if r := recover(); r != nil {
      tx.Rollback()
    }
  }()
  
  if err := txFunc(tx); err != nil {
    tx.Rollback()
    return err
  }
  
  return tx.Commit().Error
}
Tests to update:
- Update repository tests
- Update service tests that mock repositories
---
Phase 7: Mappers Simplification (Low Impact)
Current State
- Separate mapper files
- Panic on validation errors
- Boilerplate field mapping
Refactor
Files to modify:
- backend/internal/infra/postgres/usersMappers.go - DELETE, inline
- backend/internal/infra/postgres/participantsMappers.go - DELETE, inline
- backend/internal/infra/postgres/conversationMappers.go - SIMPLIFY
Changes:
// Inline simple mappers in repository files
func (r *userRepository) GetByID(id uuid.UUID) (*domain.User, error) {
  var model User
  err := r.db.Where(&User{ID: id}).First(&model).Error
  if err != nil {
    return nil, err
  }
  
  // Inline conversion with proper error handling
  if err := ValidateUsername(model.Name); err != nil {
    return nil, err
  }
  
  return &domain.User{
    ID:           model.ID,
    Avatar:       model.Avatar,
    Name:         model.Name,
    PasswordHash: model.Password,
    RefreshToken: model.RefreshToken,
  }, nil
}
// Remove panic-based mappers, handle errors properly
Tests to update:
- backend/internal/infra/postgres/conversationMappers_test.go
---
Phase 8: Keep WebSocket Subscription Sync
No changes - Required for horizontal scaling.
---
Test Update Strategy
Test File Modifications:
1. Domain tests - Update to use new validation functions instead of value objects
2. Service tests - Update mock interfaces and test helpers
3. Cache tests - Rewrite for simplified caching pattern
4. Repository tests - Update for inlined mappers
5. Notification tests - Create new tests for consolidated service
Test Helpers to Update:
// Before
userName, _ := domain.NewUserName("test")
userPassword, _ := domain.NewUserPassword("test", func(p []byte) ([]byte, error) { return p, nil })
user := domain.NewUser(userID, userName, userPassword)
// After
_ = domain.ValidateUsername("test") // Just validate
user := &domain.User{
  ID:           userID,
  Name:         "test",
  PasswordHash: "hashed",
  Avatar:       "T",
}
---
Implementation Order
1. Phase 1: Notification System (hardest, highest impact)
2. Phase 2: Value Objects (affects many files)
3. Phase 3: Caching Layer (medium complexity)
4. Phase 4: Domain Events (simplifies tests)
5. Phase 5: Event Bus (quick win)
6. Phase 6: Repository Pattern (affects repositories)
7. Phase 7: Mappers (cleanup)
---
Risk Assessment
| Phase | Risk | Mitigation |
|-------|------|------------|
| Notification | High | Run tests after each change, verify broadcast still works |
| Value Objects | Medium | Update all test helpers before changing domain code |
| Caching | Medium | Keep cache invalidation service unchanged |
| Events | Low | Simple structural change |
| Event Bus | Low | Interface removal is straightforward |
| Repository | Low | Mostly internal refactoring |
| Mappers | Low | Inline changes are safe |
---
Expected Outcomes
- Files reduced: ~7 files deleted
- Lines of code: ~500-700 lines removed
- Abstraction layers: 3-4 layers removed
- Test complexity: Reduced (simpler mocks, fewer interfaces)
- Maintainability: Significantly improved
- Performance: Similar or slightly better (removed channel overhead)
---
Verification Steps
After each phase:
1. Run go test ./... - all tests pass
2. Run golangci-lint run - no new warnings
3. Run go build ./cmd/server - compiles successfully
4. Manual testing of WebSocket connections and notifications