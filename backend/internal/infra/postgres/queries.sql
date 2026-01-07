-- User queries

-- name: StoreUser :exec
INSERT INTO users (id, avatar, name, password, refresh_token)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateUser :exec
UPDATE users
SET avatar = $2, name = $3, password = $4, refresh_token = $5, updated_at = NOW()
WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: FindUserByUsername :one
SELECT * FROM users
WHERE name = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateUserRefreshToken :exec
UPDATE users
SET refresh_token = $2, updated_at = NOW()
WHERE id = $1;

-- Conversation queries

-- name: StoreConversation :exec
INSERT INTO conversations (id, type, is_active)
VALUES ($1, $2, $3);

-- name: UpdateConversation :exec
UPDATE conversations
SET type = $2, is_active = $3, updated_at = NOW()
WHERE id = $1;

-- name: GetConversationByID :one
SELECT * FROM conversations
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: DeactivateConversation :exec
UPDATE conversations
SET is_active = FALSE, updated_at = NOW()
WHERE id = $1;

-- GroupConversation queries

-- name: StoreGroupConversation :exec
INSERT INTO group_conversations (id, name, avatar, conversation_id, owner_id)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateGroupConversation :exec
UPDATE group_conversations
SET name = $2, avatar = $3, updated_at = NOW()
WHERE conversation_id = $1;

-- name: GetGroupConversationByID :one
SELECT * FROM group_conversations
WHERE conversation_id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetGroupConversationByOwnerID :many
SELECT gc.*, c.*
FROM group_conversations gc
JOIN conversations c ON c.id = gc.conversation_id
WHERE gc.owner_id = $1 AND gc.deleted_at IS NULL AND c.deleted_at IS NULL;

-- name: GetGroupConversationWithOwner :one
SELECT
    gc.id,
    gc.name,
    gc.avatar,
    gc.conversation_id,
    gc.owner_id,
    c.type as conversation_type,
    c.is_active as conversation_is_active,
    p.id as owner_participant_id,
    p.user_id as owner_user_id,
    p.conversation_id as owner_conversation_id,
    p.is_active as owner_is_active
FROM group_conversations gc
JOIN conversations c ON c.id = gc.conversation_id
JOIN participants p ON p.conversation_id = gc.conversation_id AND p.user_id = gc.owner_id
WHERE gc.conversation_id = $1
  AND gc.deleted_at IS NULL
  AND c.deleted_at IS NULL
  AND p.deleted_at IS NULL
LIMIT 1;

-- Participant queries

-- name: StoreParticipant :exec
INSERT INTO participants (id, conversation_id, user_id, is_active)
VALUES ($1, $2, $3, $4);

-- name: UpdateParticipant :exec
UPDATE participants
SET is_active = $2, updated_at = NOW()
WHERE id = $1;

-- name: GetParticipantByID :one
SELECT * FROM participants
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: FindParticipantByConversationAndUser :one
SELECT * FROM participants
WHERE conversation_id = $1 AND user_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: GetParticipantsForConversation :many
SELECT p.*, u.name, u.avatar
FROM participants p
JOIN users u ON u.id = p.user_id
WHERE p.conversation_id = $1 AND p.is_active = TRUE AND p.deleted_at IS NULL AND u.deleted_at IS NULL
ORDER BY p.created_at ASC
LIMIT $2 OFFSET $3;

-- name: GetParticipantsIDsByConversationID :many
SELECT user_id
FROM participants
WHERE conversation_id = $1 AND deleted_at IS NULL
ORDER BY created_at ASC;

-- name: DeactivateParticipant :exec
UPDATE participants
SET is_active = FALSE, updated_at = NOW()
WHERE id = $1;

-- Message queries

-- name: StoreMessage :exec
INSERT INTO messages (id, conversation_id, user_id, content, type)
VALUES ($1, $2, $3, $4, $5);

-- name: GetMessagesForConversation :many
SELECT m.*, u.name as user_name, u.avatar as user_avatar
FROM messages m
JOIN users u ON u.id = m.user_id
WHERE m.conversation_id = $1 AND m.deleted_at IS NULL
ORDER BY m.created_at ASC
LIMIT $2 OFFSET $3;

-- name: GetMessageByID :one
SELECT m.*, u.name as user_name, u.avatar as user_avatar
FROM messages m
JOIN users u ON u.id = m.user_id
WHERE m.id = $1 AND m.deleted_at IS NULL
LIMIT 1;

-- Complex queries for read model

-- name: GetContacts :many
SELECT id, name, avatar
FROM users
WHERE deleted_at IS NULL AND id != $1
LIMIT $2 OFFSET $3;

-- name: GetParticipantsByConversationID :many
SELECT u.id, u.name, u.avatar
FROM users u
JOIN participants p ON p.user_id = u.id
WHERE p.conversation_id = $1
  AND p.is_active = TRUE
  AND p.deleted_at IS NULL
  AND u.deleted_at IS NULL
LIMIT $2 OFFSET $3;

-- name: GetPotentialInvitees :many
SELECT u.id, u.name, u.avatar
FROM users u
WHERE u.deleted_at IS NULL
  AND u.id NOT IN (
    SELECT user_id
    FROM participants
    WHERE conversation_id = $1
      AND is_active = TRUE
      AND deleted_at IS NULL
  )
LIMIT $2 OFFSET $3;

-- name: GetUserByIDDTO :one
SELECT id, name, avatar
FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetConversationMessagesWithUser :many
SELECT
    m.id,
    m.type,
    m.created_at,
    m.conversation_id,
    m.content,
    u.id as user_id,
    u.name as user_name,
    u.avatar as user_avatar
FROM messages m
LEFT JOIN users u ON u.id = m.user_id
WHERE m.conversation_id = $1 AND m.deleted_at IS NULL
ORDER BY m.created_at ASC
LIMIT $2 OFFSET $3;

-- name: GetNotificationMessageWithUser :one
SELECT
    m.id,
    m.type,
    m.created_at,
    m.conversation_id,
    m.content,
    u.id as user_id,
    u.name as user_name,
    u.avatar as user_avatar
FROM messages m
LEFT JOIN users u ON u.id = m.user_id
WHERE m.id = $1 AND m.deleted_at IS NULL
LIMIT 1;

-- name: GetUserConversations :many
WITH last_messages AS (
    SELECT conversation_id, MAX(created_at) as max_created_at
    FROM messages
    GROUP BY conversation_id
)
SELECT
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
    AND op.user_id <> $1
    AND op.is_active = TRUE
LEFT JOIN users ou ON ou.id = op.user_id
WHERE p.user_id = $1
  AND c.is_active = TRUE
  AND c.deleted_at IS NULL
  AND p.is_active = TRUE
  AND p.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetConversationFull :one
WITH participants_count AS (
    SELECT COUNT(*) as count
    FROM participants
    WHERE conversation_id = $1
      AND is_active = TRUE
      AND deleted_at IS NULL
)
SELECT
    c.id as conversation_id,
    c.created_at,
    c.type,
    ou.id as other_user_id,
    ou.name as other_user_name,
    ou.avatar as other_user_avatar,
    gc.avatar as group_avatar,
    gc.name as group_name,
    gc.owner_id as group_owner_id,
    pc.count as participants_count,
    up.id as user_participant_id
FROM conversations c
LEFT JOIN participants op ON op.conversation_id = c.id
    AND op.user_id <> $2
    AND op.is_active = TRUE
LEFT JOIN users ou ON ou.id = op.user_id
LEFT JOIN group_conversations gc ON gc.conversation_id = c.id
LEFT JOIN participants pc_sub ON pc_sub.conversation_id = c.id AND pc_sub.user_id = $2 AND pc_sub.is_active = TRUE
LEFT JOIN participants up ON up.conversation_id = c.id AND up.user_id = $2 AND up.is_active = TRUE
CROSS JOIN participants_count pc
WHERE c.id = $1 AND c.is_active = TRUE AND c.deleted_at IS NULL
LIMIT 1;

-- name: GetDirectConversationBetweenUsers :one
SELECT c.*
FROM conversations c
JOIN participants p1 ON p1.conversation_id = c.id AND p1.user_id = $1 AND p1.is_active = TRUE
JOIN participants p2 ON p2.conversation_id = c.id AND p2.user_id = $2 AND p2.is_active = TRUE
WHERE c.type = 1
  AND c.is_active = TRUE
  AND c.deleted_at IS NULL
  AND p1.deleted_at IS NULL
  AND p2.deleted_at IS NULL
LIMIT 1;
