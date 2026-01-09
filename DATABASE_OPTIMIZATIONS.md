# Database Optimization Opportunities

This document outlines operations that could be moved from application level to database level for improved performance and code quality.

---

## High Priority

### 1. Bug Fix: `GetConversationIDsByUserID`

**File:** `internal/infra/postgres/participantsRepository.go:86-93`

**Issue:** Currently uses wrong query (`GetParticipantByID` instead of finding by user ID). This likely causes incorrect behavior in websocket channel subscriptions.

**Current Code:**
```go
func (r *participantRepository) GetConversationIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
    participants, err := r.queries.GetParticipantByID(ctx, uuidToPgtype(userID))  // BUG: Wrong query!
    // ...
    return []uuid.UUID{pgtypeToUUID(participants.ConversationID)}, nil  // Only returns 1 conversation
}
```

**Fix:** Add new query to `queries.sql`:
```sql
-- name: GetConversationIDsByUserID :many
SELECT conversation_id
FROM participants
WHERE user_id = $1 AND is_active = TRUE AND deleted_at IS NULL;
```

---

### 2. DirectConversation.GetByID - 2 queries → 1

**File:** `internal/infra/postgres/directConversationRepository.go:56-88`

**Current:** Makes 2 separate queries:
```go
participants, err := r.queries.GetParticipantsIDsByConversationID(ctx, uuidToPgtype(id))
// ...
conv, err := r.queries.GetConversationByID(ctx, uuidToPgtype(id))
```

**Suggested:** Single query with JOIN:
```sql
-- name: GetDirectConversationWithParticipants :one
SELECT c.*, ARRAY_AGG(p.user_id) as participant_user_ids
FROM conversations c
JOIN participants p ON p.conversation_id = c.id
WHERE c.id = $1 AND c.deleted_at IS NULL AND p.deleted_at IS NULL
GROUP BY c.id;
```

---

### 3. System Message Operations - 3 queries → 1

**File:** `internal/services/systemMessages.go`

**Issue:** All 4 methods (`SaveJoinedMessage`, `SaveLeftMessage`, `SaveInvitedMessage`, `SaveRenamedMessage`) fetch conversation + user just to validate before inserting.

**Current Pattern (repeated 4 times):**
```go
conversation, err := s.groupConversations.GetByID(ctx, conversationID)  // Query 1
user, err := s.users.GetByID(ctx, userID)                               // Query 2
message, err := conversation.SendJoinedConversationMessage(messageID, user)
s.messages.Store(ctx, message)                                          // Query 3
```

**Suggested:** Single INSERT with validation subquery:
```sql
-- name: StoreSystemMessage :exec
INSERT INTO messages (id, conversation_id, user_id, content, type, created_at)
SELECT $1, $2, $3, $4, $5, NOW()
WHERE EXISTS (
    SELECT 1 FROM conversations c
    JOIN group_conversations gc ON gc.conversation_id = c.id
    WHERE c.id = $2 AND c.is_active = TRUE AND c.deleted_at IS NULL
)
AND EXISTS (
    SELECT 1 FROM users WHERE id = $3 AND deleted_at IS NULL
);
```

---

## Medium Priority

### 4. Message Text Formatting - Move to SQL

**Files:** `internal/infra/postgres/queriesRepository.go:152-163, 196-208, 254-266`

**Issue:** Same switch statement duplicated 3 times in Go code:
```go
switch messageTypesMap[uint8(msg.Type)] {
case domain.MessageTypeText:
    text = msg.Content
case domain.MessageTypeRenamedConversation:
    text = msg.UserName.String + " renamed chat to " + msg.Content
case domain.MessageTypeJoinedConversation:
    text = msg.UserName.String + " joined"
case domain.MessageTypeLeftConversation:
    text = msg.UserName.String + " left"
case domain.MessageTypeInvitedToConversation:
    text = msg.UserName.String + " was invited"
}
```

**Suggested:** Move to SQL CASE expression:
```sql
-- name: GetConversationMessagesWithFormattedText :many
SELECT
    m.id,
    m.type,
    m.created_at,
    m.conversation_id,
    CASE m.type
        WHEN 1 THEN m.content
        WHEN 2 THEN u.name || ' renamed chat to ' || m.content
        WHEN 3 THEN u.name || ' joined'
        WHEN 4 THEN u.name || ' left'
        WHEN 5 THEN u.name || ' was invited'
    END as formatted_text,
    u.id as user_id,
    u.name as user_name,
    u.avatar as user_avatar
FROM messages m
LEFT JOIN users u ON u.id = m.user_id
WHERE m.conversation_id = $1 AND m.deleted_at IS NULL
ORDER BY m.created_at ASC
LIMIT $2 OFFSET $3;
```

**Benefits:**
- Eliminates code duplication
- Reduces Go-side processing
- Single source of truth for message formatting

---

### 5. Compound Operations for Chat Actions

These operations currently make many sequential database calls that could be consolidated:

| Operation | File | Current Calls | Reducible To |
|-----------|------|---------------|--------------|
| `Kick()` | `conversation.go:409-460` | 10 | 2-3 |
| `Invite()` | `conversation.go:356-407` | 9 | 2-3 |
| `Join()` | `conversation.go:262-307` | 8 | 2-3 |
| `Leave()` | `conversation.go:309-354` | 8 | 2-3 |
| `Rename()` | `conversation.go:143-186` | 7 | 2-3 |

#### Example: Kick Operation

**Current Flow (10 queries):**
1. `groupConversations.GetByID()`
2. `participants.GetByConversationIDAndUserID()` (kicker)
3. `participants.GetByConversationIDAndUserID()` (target)
4. `participants.Update()`
5-7. `systemMessages.SaveLeftMessage()` (3 queries)
8. `queries.GetConversation()`

**Suggested:** Combined query for validation and update:
```sql
-- name: KickParticipant :one
WITH validation AS (
    SELECT 
        gc.owner_id,
        kicker.id as kicker_id,
        target.id as target_id,
        target.user_id as target_user_id
    FROM group_conversations gc
    JOIN participants kicker ON kicker.conversation_id = gc.conversation_id 
        AND kicker.user_id = $2 AND kicker.is_active = TRUE AND kicker.deleted_at IS NULL
    JOIN participants target ON target.conversation_id = gc.conversation_id 
        AND target.user_id = $3 AND target.is_active = TRUE AND target.deleted_at IS NULL
    WHERE gc.conversation_id = $1 
        AND gc.owner_id = $2  -- kicker must be owner
        AND gc.deleted_at IS NULL
)
UPDATE participants 
SET is_active = FALSE, updated_at = NOW()
FROM validation v
WHERE participants.id = v.target_id
RETURNING v.target_user_id;
```

#### Example: Join Operation

```sql
-- name: JoinConversation :one
WITH valid_conversation AS (
    SELECT gc.conversation_id
    FROM group_conversations gc
    JOIN conversations c ON c.id = gc.conversation_id
    WHERE gc.conversation_id = $1 
        AND c.is_active = TRUE 
        AND c.deleted_at IS NULL
),
valid_user AS (
    SELECT id, name FROM users WHERE id = $2 AND deleted_at IS NULL
),
new_participant AS (
    INSERT INTO participants (id, conversation_id, user_id, is_active, created_at)
    SELECT $3, vc.conversation_id, vu.id, TRUE, NOW()
    FROM valid_conversation vc, valid_user vu
    RETURNING user_id
)
INSERT INTO messages (id, conversation_id, user_id, content, type, created_at)
SELECT $4, $1, np.user_id, '', 3, NOW()  -- type 3 = joined
FROM new_participant np
RETURNING user_id;
```

---

### 6. Batch Insert Participants

**File:** `internal/infra/postgres/directConversationRepository.go:39-50`

**Current:** Loops to insert participants one-by-one:
```go
for _, participant := range conversation.Participants {
    err = r.queries.StoreParticipant(ctx, db.StoreParticipantParams{
        // ...
    })
}
```

**Suggested:** Batch insert using `unnest`:
```sql
-- name: StoreParticipantsBatch :exec
INSERT INTO participants (id, conversation_id, user_id, is_active, created_at)
SELECT unnest($1::uuid[]), $2, unnest($3::uuid[]), TRUE, NOW();
```

**Go usage:**
```go
participantIDs := make([]uuid.UUID, len(conversation.Participants))
userIDs := make([]uuid.UUID, len(conversation.Participants))
for i, p := range conversation.Participants {
    participantIDs[i] = p.ID
    userIDs[i] = p.UserID
}
err = r.queries.StoreParticipantsBatch(ctx, participantIDs, conversationID, userIDs)
```

---

## Summary

| Category | Items | Estimated Improvement |
|----------|-------|----------------------|
| Bug fixes | 1 | Correctness |
| Query consolidation | 2 | 50% fewer round-trips |
| Move logic to SQL | 1 | Removes code duplication |
| Compound operations | 5 | 60-70% fewer round-trips |
| Batch operations | 1 | N queries → 1 query |

---

## Implementation Priority

1. **Immediate:** Fix `GetConversationIDsByUserID` bug - this is likely causing incorrect behavior
2. **High Value:** Create compound queries for `Kick`, `Invite`, `Join`, `Leave`, `Rename` operations
3. **Code Quality:** Move message text formatting to SQL to eliminate duplication
4. **Performance:** Add batch insert for participants when creating direct conversations
5. **Consider:** Using PostgreSQL stored procedures for complex multi-step operations requiring transactional consistency

---

## Notes

- All suggested SQL uses PostgreSQL-specific features (CTEs, `unnest`, `ARRAY_AGG`)
- Compound operations maintain the same business logic validation, just moved to SQL
- Consider adding database indexes if not already present:
  - `participants(conversation_id, user_id)` for lookup operations
  - `participants(user_id)` for `GetConversationIDsByUserID`
  - `messages(conversation_id, created_at)` for message pagination
