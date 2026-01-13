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
INSERT INTO conversations (id, type)
VALUES ($1, $2);

-- name: UpdateConversation :exec
UPDATE conversations
SET type = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteConversation :exec
UPDATE conversations
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- GroupConversation queries

-- name: StoreGroupConversation :exec
INSERT INTO group_conversations (id, name, avatar, conversation_id, owner_id)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateGroupConversation :exec
UPDATE group_conversations
SET name = $2, avatar = $3, updated_at = NOW()
WHERE conversation_id = $1;

-- name: RenameGroupConversation :exec
UPDATE group_conversations
SET name = $2, updated_at = NOW()
WHERE conversation_id = $1;

-- name: GetGroupConversationWithOwner :one
SELECT
    gc.id,
    gc.name,
    gc.avatar,
    gc.conversation_id,
    gc.owner_id,
    c.type as conversation_type,
    p.id as owner_participant_id,
    p.user_id as owner_user_id,
    p.conversation_id as owner_conversation_id
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
INSERT INTO participants (id, conversation_id, user_id)
VALUES ($1, $2, $3);

-- name: StoreParticipantsBatch :exec
INSERT INTO participants (id, conversation_id, user_id, created_at)
SELECT unnest($1::uuid[]), $2, unnest($3::uuid[]), NOW();

-- name: DeleteParticipant :exec
UPDATE participants
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: FindParticipantByConversationAndUser :one
SELECT * FROM participants
WHERE conversation_id = $1 AND user_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: GetParticipantsIDsByConversationID :many
SELECT user_id
FROM participants
WHERE conversation_id = $1 AND deleted_at IS NULL
ORDER BY created_at ASC;

-- name: GetDirectConversationWithParticipants :one
SELECT c.id, c.type, c.created_at, c.updated_at,
       ARRAY_AGG(p.user_id) as participant_user_ids
FROM conversations c
JOIN participants p ON p.conversation_id = c.id
WHERE c.id = $1 AND c.deleted_at IS NULL AND p.deleted_at IS NULL
GROUP BY c.id;

-- Message queries

-- name: StoreMessage :exec
INSERT INTO messages (id, conversation_id, user_id, content, type, created_at)
VALUES ($1, $2, $3, $4, $5, NOW());

-- name: GetMessageWithUser :one
SELECT
    m.id, m.type, m.created_at, m.conversation_id, m.content,
    u.id as user_id, u.name as user_name, u.avatar as user_avatar
FROM messages m
JOIN users u ON u.id = m.user_id
WHERE m.id = $1 AND u.deleted_at IS NULL
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
      AND deleted_at IS NULL
  )
LIMIT $2 OFFSET $3;

-- name: GetUsersByIDs :many
SELECT id, name, avatar
FROM users
WHERE id = ANY($1::uuid[])
  AND deleted_at IS NULL;

-- name: GetConversationMessagesRaw :many
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

-- name: GetNotificationMessageRaw :one
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
LEFT JOIN users ou ON ou.id = op.user_id
WHERE p.user_id = $1
  AND c.deleted_at IS NULL
  AND p.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetConversationFull :one
WITH participants_count AS (
    SELECT COUNT(*) as count
    FROM participants
    WHERE conversation_id = $1
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
LEFT JOIN users ou ON ou.id = op.user_id
LEFT JOIN group_conversations gc ON gc.conversation_id = c.id
LEFT JOIN participants pc_sub ON pc_sub.conversation_id = c.id AND pc_sub.user_id = $2
LEFT JOIN participants up ON up.conversation_id = c.id AND up.user_id = $2
CROSS JOIN participants_count pc
WHERE c.id = $1 AND c.deleted_at IS NULL
LIMIT 1;

-- name: GetDirectConversationBetweenUsers :one
SELECT c.*
FROM conversations c
JOIN participants p1 ON p1.conversation_id = c.id AND p1.user_id = $1
JOIN participants p2 ON p2.conversation_id = c.id AND p2.user_id = $2
WHERE c.type = 1
  AND c.deleted_at IS NULL
  AND p1.deleted_at IS NULL
  AND p2.deleted_at IS NULL
LIMIT 1;

-- name: GetConversationIDsByUserID :many
SELECT conversation_id
FROM participants
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: LeaveConversationAtomic :execrows
UPDATE participants
SET deleted_at = NOW(), updated_at = NOW()
WHERE conversation_id = $1
  AND user_id = $2
  AND deleted_at IS NULL;

-- name: InviteToConversationAtomic :one
WITH valid_conversation AS (
    SELECT gc.conversation_id as conv_id
    FROM group_conversations gc
    JOIN conversations c ON c.id = gc.conversation_id
    WHERE gc.conversation_id = $1 
        AND c.deleted_at IS NULL
        AND gc.deleted_at IS NULL
),
valid_invitee AS (
    SELECT u.id as user_id FROM users u WHERE u.id = $2 AND u.deleted_at IS NULL
),
new_participant AS (
    INSERT INTO participants (id, conversation_id, user_id, created_at)
    SELECT $3, vc.conv_id, vi.user_id, NOW()
    FROM valid_conversation vc, valid_invitee vi
    ON CONFLICT DO NOTHING
    RETURNING user_id
)
SELECT user_id FROM new_participant;

-- name: IsMember :one
SELECT EXISTS(
    SELECT 1 FROM participants
    WHERE conversation_id = $1 AND user_id = $2 AND deleted_at IS NULL
);

-- name: IsMemberOwner :one
SELECT EXISTS(
    SELECT 1 FROM group_conversations gc
    JOIN conversations c ON c.id = gc.conversation_id
    JOIN participants p ON p.conversation_id = gc.conversation_id AND p.user_id = $2
    WHERE gc.conversation_id = $1 AND gc.owner_id = $2
      AND gc.deleted_at IS NULL AND c.deleted_at IS NULL AND p.deleted_at IS NULL
);

-- name: StoreMessageAndReturn :one
WITH new_message AS (
    INSERT INTO messages (id, conversation_id, user_id, content, type, created_at)
    VALUES ($1, $2, $3, $4, $5, NOW())
    RETURNING id, type, created_at, conversation_id, content
)
SELECT
    nm.id, nm.type, nm.created_at, nm.conversation_id,
    nm.content as formatted_text,
    u.id as user_id, u.name as user_name, u.avatar as user_avatar
FROM new_message nm
JOIN users u ON u.id = nm.user_id
WHERE u.deleted_at IS NULL;

-- name: RenameConversationAndReturn :execrows
UPDATE group_conversations
SET name = $2, updated_at = NOW()
WHERE conversation_id = $1
  AND deleted_at IS NULL
  AND EXISTS (
    SELECT 1 FROM conversations c WHERE c.id = conversation_id AND c.deleted_at IS NULL
  );