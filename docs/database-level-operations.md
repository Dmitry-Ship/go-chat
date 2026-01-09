# Database-Level Operations Recommendations

## Overview

This document identifies operations that could benefit from being moved from application-level (Go) to database-level (PostgreSQL) for improved performance, consistency, and reduced complexity.

---

## 1. System Messages (High Priority)

### Current Implementation
- **File**: `internal/services/systemMessages.go`
- **Methods**: `SaveJoinedMessage`, `SaveLeftMessage`, `SaveInvitedMessage`, `SaveRenamedMessage`

Currently involves multiple database fetches followed by manual message creation:
```go
conversation := s.groupConversations.GetByID(ctx, conversationID)
user := s.users.GetByID(ctx, userID)
message := conversation.SendJoinedConversationMessage(messageID, user)
s.messages.Store(ctx, message)
```

### Proposed Database Solution
Use database triggers to auto-create system messages on participant changes:

```sql
-- Trigger for joined messages
CREATE OR REPLACE FUNCTION create_joined_message()
RETURNS TRIGGER AS $$
DECLARE
    message_id UUID;
BEGIN
    message_id := gen_random_uuid();
    INSERT INTO messages (id, conversation_id, user_id, content, type)
    VALUES (
        message_id,
        NEW.conversation_id,
        NEW.user_id,
        'joined the conversation',
        3 -- Assuming 3 is the joined message type
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER participant_joined_trigger
AFTER INSERT ON participants
FOR EACH ROW
WHEN (NEW.is_active = TRUE)
EXECUTE FUNCTION create_joined_message();
```

### Benefits
- Eliminates N+1 queries (conversation + user fetches)
- Ensures message consistency
- Reduces application code complexity
- Automatic message creation even if application crashes mid-operation

---

## 2. Direct Conversation Upsert (High Priority)

### Current Implementation
- **File**: `internal/services/conversation.go:85-111`
- **Method**: `StartDirectConversation`

Currently checks for existence, then creates if not found:
```go
existingConversationID, err := s.directConversations.GetID(ctx, fromUserID, toUserID)
if err == nil {
    return existingConversationID, nil
}
newConversationID := uuid.New()
conversation, err := domain.NewDirectConversation(newConversationID, toUserID, fromUserID)
s.directConversations.Store(ctx, conversation)
```

### Proposed Database Solution
Use PostgreSQL `INSERT ... ON CONFLICT` or a stored procedure:

```sql
-- name: CreateOrGetDirectConversation :one
WITH new_conv AS (
    INSERT INTO conversations (id, type, is_active)
    SELECT gen_random_uuid(), 1, TRUE
    WHERE NOT EXISTS (
        SELECT 1 FROM conversations c
        JOIN participants p1 ON p1.conversation_id = c.id AND p1.user_id = $1 AND p1.is_active = TRUE
        JOIN participants p2 ON p2.conversation_id = c.id AND p2.user_id = $2 AND p2.is_active = TRUE
        WHERE c.type = 1 AND c.is_active = TRUE AND c.deleted_at IS NULL
    )
    RETURNING id
)
SELECT COALESCE(
    (SELECT id FROM new_conv),
    (SELECT c.id FROM conversations c
     JOIN participants p1 ON p1.conversation_id = c.id AND p1.user_id = $1 AND p1.is_active = TRUE
     JOIN participants p2 ON p2.conversation_id = c.id AND p2.user_id = $2 AND p2.is_active = TRUE
     WHERE c.type = 1 AND c.is_active = TRUE AND c.deleted_at IS NULL
     LIMIT 1)
) as conversation_id;
```

### Benefits
- Single query instead of multiple
- Atomic operation prevents race conditions
- Eliminates potential duplicate conversations
- Simpler application logic

---

## 3. Permission Checks (Medium Priority)

### Current Implementation
- **File**: `internal/domain/groupConversation.go`
- **Methods**: `Delete`, `Rename`, `Invite`, `Kick`

Currently fetches full objects to check permissions:
```go
func (groupConversation *GroupConversation) Delete(participant *Participant) error {
    if !groupConversation.isJoined(participant) {
        return ErrorUserNotInConversation
    }
    if groupConversation.Owner.UserID != participant.UserID {
        return ErrorUserNotOwner
    }
    ...
}
```

### Proposed Database Solution
Move permission checks into stored procedures:

```sql
-- name: RenameGroupConversation :exec
UPDATE group_conversations gc
SET name = $2, avatar = LEFT($2, 1)
WHERE gc.conversation_id = $1
  AND EXISTS (
    SELECT 1 FROM participants p
    WHERE p.conversation_id = gc.conversation_id
      AND p.user_id = $3
      AND p.is_active = TRUE
      AND p.deleted_at IS NULL
  )
  AND gc.owner_id = $3
  AND gc.deleted_at IS NULL
  AND (SELECT is_active FROM conversations WHERE id = $1 AND deleted_at IS NULL) = TRUE;
```

Alternative: Use Row-Level Security (RLS):
```sql
CREATE POLICY group_owner_policy ON group_conversations
    FOR UPDATE
    TO authenticated_role
    USING (
        owner_id IN (
            SELECT user_id FROM participants
            WHERE conversation_id = group_conversations.conversation_id
              AND user_id = current_user_id()
              AND is_active = TRUE
        )
    );
```

### Benefits
- Reduced data transfer (no full object fetch)
- Closer to data integrity
- Business logic enforced at database level
- Better for multi-server consistency

---

## 4. Participant Uniqueness (Medium Priority)

### Current Implementation
- **File**: `internal/domain/groupConversation.go:122-132`
- **Method**: `Invite`

Currently checked in domain logic:
```go
if inviter.UserID == invitee.ID {
    return nil, ErrorCannotInviteOneself
}
```

### Proposed Database Solution
Add unique constraint and handle conflicts:

```sql
-- Add unique constraint
ALTER TABLE participants
ADD CONSTRAINT unique_conversation_user
UNIQUE (conversation_id, user_id)
WHERE deleted_at IS NULL;

-- Handle with INSERT ... ON CONFLICT
INSERT INTO participants (id, conversation_id, user_id, is_active)
VALUES ($1, $2, $3, TRUE)
ON CONFLICT (conversation_id, user_id)
DO UPDATE SET is_active = TRUE, deleted_at = NULL;
```

### Benefits
- Database-level integrity (impossible to create duplicates)
- Simplifies application logic
- Automatic conflict resolution
- Consistent across all application entry points

---

## 5. Materialized Views for Conversation Lists (Low Priority)

### Current Implementation
- **File**: `internal/infra/postgres/queries.sql:221-259`
- **Query**: `GetUserConversations`

Complex CTE executed on every request:
```sql
WITH last_messages AS (
    SELECT conversation_id, MAX(created_at) as max_created_at
    FROM messages
    GROUP BY conversation_id
)
SELECT c.id, m.content, u.name, ...
FROM conversations c
JOIN participants p ON ...
LEFT JOIN last_messages lm ON ...
LEFT JOIN messages m ON ...
LEFT JOIN users u ON ...
LEFT JOIN group_conversations gc ON ...
```

### Proposed Database Solution
Create materialized view and refresh periodically or trigger-based:

```sql
CREATE MATERIALIZED VIEW user_conversations_mv AS
WITH last_messages AS (
    SELECT conversation_id, MAX(created_at) as max_created_at
    FROM messages
    GROUP BY conversation_id
)
SELECT
    p.user_id,
    c.id as conversation_id,
    c.created_at,
    c.type,
    m.id as message_id,
    m.type as message_type,
    m.content as message_content,
    m.created_at as message_created_at,
    m.user_id as message_user_id,
    u.name as message_user_name,
    u.avatar as message_user_avatar,
    gc.avatar as group_avatar,
    gc.name as group_name,
    ou.id as other_user_id,
    ou.name as other_user_name,
    ou.avatar as other_user_avatar
FROM conversations c
JOIN participants p ON p.conversation_id = c.id
LEFT JOIN last_messages lm ON lm.conversation_id = c.id
LEFT JOIN messages m ON m.conversation_id = c.id AND m.created_at = lm.max_created_at
LEFT JOIN users u ON u.id = m.user_id
LEFT JOIN group_conversations gc ON gc.conversation_id = c.id
LEFT JOIN participants op ON op.conversation_id = c.id
    AND op.user_id <> p.user_id
    AND op.is_active = TRUE
LEFT JOIN users ou ON ou.id = op.user_id
WHERE c.is_active = TRUE
  AND c.deleted_at IS NULL
  AND p.is_active = TRUE
  AND p.deleted_at IS NULL;

-- Create index for querying
CREATE INDEX idx_user_conversations_mv_user_id ON user_conversations_mv(user_id);

-- Refresh function
CREATE OR REPLACE FUNCTION refresh_user_conversations_mv()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_conversations_mv;
END;
$$ LANGUAGE plpgsql;

-- Trigger to refresh incrementally
CREATE TRIGGER refresh_user_conversations_mv
AFTER INSERT OR UPDATE OR DELETE ON messages
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_user_conversations_mv();
```

Query becomes:
```sql
SELECT * FROM user_conversations_mv
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
```

### Benefits
- Faster reads (no complex joins on every request)
- Reduced CPU load
- Better scalability for large datasets
- Can be refreshed asynchronously

---

## 6. Soft Delete Cascades (Low Priority)

### Current Implementation
- **Files**: Multiple services
- **Pattern**: Manual `is_active` flag updates

Example from `conversationService.DeleteGroupConversation:114-140`:
```go
conversation.IsActive = false
s.groupConversations.Update(ctx, conversation)
```

### Proposed Database Solution
Use database triggers to cascade soft deletes:

```sql
CREATE OR REPLACE FUNCTION deactivate_related_records()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.is_active = TRUE AND NEW.is_active = FALSE THEN
        -- Deactivate all participants
        UPDATE participants
        SET is_active = FALSE, updated_at = NOW()
        WHERE conversation_id = NEW.id
          AND is_active = TRUE
          AND deleted_at IS NULL;

        -- Mark messages as inactive
        UPDATE messages
        SET deleted_at = NOW()
        WHERE conversation_id = NEW.id
          AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cascade_soft_delete_conversation
AFTER UPDATE ON conversations
FOR EACH ROW
WHEN (OLD.is_active = TRUE AND NEW.is_active = FALSE)
EXECUTE FUNCTION deactivate_related_records();
```

### Benefits
- Automatic consistency
- No missing updates
- Simplifies application logic
- Data integrity guaranteed

---

## Summary Table

| Operation | Current Location | Move To | Priority | Impact |
|-----------|-----------------|---------|----------|---------|
| System messages (join/leave/etc) | Application layer + N+1 queries | Database triggers | **High** | High performance gain |
| Direct conversation upsert | Check then insert pattern | `INSERT ... ON CONFLICT` | **High** | Race condition prevention |
| Permission checks | Fetch full objects in Go | Stored procedures | Medium | Reduced data transfer |
| Participant uniqueness | Domain logic validation | Unique constraints | Medium | Data integrity |
| Conversation list queries | CTE per request | Materialized views | Low | Scalability improvement |
| Soft delete cascades | Manual updates in services | Database triggers | Low | Automatic consistency |

---

## Implementation Priority

### Phase 1 (Immediate - High Impact)
1. **System Messages Triggers** - Biggest performance improvement
2. **Direct Conversation Upsert** - Eliminates race conditions

### Phase 2 (Medium Priority)
3. **Permission Checks in Stored Procedures** - Reduces data transfer
4. **Participant Uniqueness Constraints** - Improves data integrity

### Phase 3 (Future Optimization)
5. **Materialized Views** - Better scalability
6. **Soft Delete Cascades** - Better consistency

---

## Migration Considerations

### Testing
- Write integration tests for each trigger/stored procedure
- Test edge cases (concurrent operations, null values, etc.)
- Ensure backward compatibility during migration

### Performance
- Benchmark before and after changes
- Monitor query execution times
- Check for performance regressions

### Code Changes
- Remove corresponding application-level logic after migration
- Update unit tests to mock database behavior appropriately
- Add integration tests for trigger logic
- Test edge cases (concurrent operations, null values, etc.)