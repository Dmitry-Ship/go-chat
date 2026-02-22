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

-- Conversation queries

-- name: StoreConversation :exec
INSERT INTO conversations (id, type, name, avatar, owner_id)
VALUES ($1, $2, $3, $4, $5);

-- name: DeleteConversation :exec
UPDATE conversations
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: RenameGroupConversation :exec
UPDATE conversations
SET name = sqlc.arg(name)::text, avatar = sqlc.arg(avatar)::text, updated_at = NOW()
WHERE id = sqlc.arg(conversation_id)
  AND type = 0
  AND deleted_at IS NULL;

-- Participant queries

-- name: StoreParticipant :exec
INSERT INTO participants (id, conversation_id, user_id)
VALUES ($1, $2, $3);

-- name: StoreParticipantsBatch :exec
INSERT INTO participants (id, conversation_id, user_id, created_at)
SELECT unnest($1::uuid[]), $2, unnest($3::uuid[]), NOW();

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
    m.user_id
FROM messages m
WHERE m.conversation_id = $1
  AND m.deleted_at IS NULL
  AND (
    sqlc.arg(cursor_created_at)::timestamptz IS NULL
    OR m.created_at < sqlc.arg(cursor_created_at)
    OR (
      m.created_at = sqlc.arg(cursor_created_at)
      AND m.id < sqlc.arg(cursor_id)
    )
  )
ORDER BY m.created_at DESC, m.id DESC
LIMIT sqlc.arg(page_limit);

-- name: GetUserConversations :many
WITH last_messages AS (
    SELECT DISTINCT ON (m.conversation_id)
        m.conversation_id,
        m.id,
        m.type,
        m.content,
        m.created_at,
        m.user_id
    FROM messages m
    WHERE m.deleted_at IS NULL
    ORDER BY m.conversation_id, m.created_at DESC, m.id DESC
)
SELECT
    c.id as conversation_id,
    c.created_at,
    c.type,
    lm.id as message_id,
    lm.type as message_type,
    lm.content as message_content,
    lm.created_at as message_created_at,
    lm.user_id as message_user_id,
    mu.name as message_user_name,
    mu.avatar as message_user_avatar,
    c.avatar as group_avatar,
    c.name as group_name,
    (
        SELECT u.id
        FROM participants op
        JOIN users u ON u.id = op.user_id
        WHERE op.conversation_id = c.id
          AND op.user_id <> $1
          AND op.deleted_at IS NULL
          AND u.deleted_at IS NULL
        ORDER BY op.created_at ASC
        LIMIT 1
    ) as other_user_id,
    COALESCE((
        SELECT u.name
        FROM participants op
        JOIN users u ON u.id = op.user_id
        WHERE op.conversation_id = c.id
          AND op.user_id <> $1
          AND op.deleted_at IS NULL
          AND u.deleted_at IS NULL
        ORDER BY op.created_at ASC
        LIMIT 1
    ), '') as other_user_name,
    (
        SELECT u.avatar
        FROM participants op
        JOIN users u ON u.id = op.user_id
        WHERE op.conversation_id = c.id
          AND op.user_id <> $1
          AND op.deleted_at IS NULL
          AND u.deleted_at IS NULL
        ORDER BY op.created_at ASC
        LIMIT 1
    ) as other_user_avatar
FROM conversations c
JOIN participants p ON p.conversation_id = c.id
  AND p.user_id = $1
  AND p.deleted_at IS NULL
LEFT JOIN last_messages lm ON lm.conversation_id = c.id
LEFT JOIN users mu ON mu.id = lm.user_id
  AND mu.deleted_at IS NULL
WHERE c.deleted_at IS NULL
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
    (
        SELECT u.id
        FROM participants op
        JOIN users u ON u.id = op.user_id
        WHERE op.conversation_id = c.id
          AND op.user_id <> $2
          AND op.deleted_at IS NULL
          AND u.deleted_at IS NULL
        ORDER BY op.created_at ASC
        LIMIT 1
    ) as other_user_id,
    COALESCE((
        SELECT u.name
        FROM participants op
        JOIN users u ON u.id = op.user_id
        WHERE op.conversation_id = c.id
          AND op.user_id <> $2
          AND op.deleted_at IS NULL
          AND u.deleted_at IS NULL
        ORDER BY op.created_at ASC
        LIMIT 1
    ), '') as other_user_name,
    (
        SELECT u.avatar
        FROM participants op
        JOIN users u ON u.id = op.user_id
        WHERE op.conversation_id = c.id
          AND op.user_id <> $2
          AND op.deleted_at IS NULL
          AND u.deleted_at IS NULL
        ORDER BY op.created_at ASC
        LIMIT 1
    ) as other_user_avatar,
    c.avatar as group_avatar,
    c.name as group_name,
    c.owner_id as group_owner_id,
    pc.count as participants_count,
    up.id as user_participant_id
FROM conversations c
LEFT JOIN participants up ON up.conversation_id = c.id
    AND up.user_id = $2
    AND up.deleted_at IS NULL
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
    SELECT c.id as conv_id
    FROM conversations c
    WHERE c.id = sqlc.arg(conversation_id)
      AND c.type = 0
      AND c.deleted_at IS NULL
),
valid_invitee AS (
    SELECT u.id as user_id FROM users u WHERE u.id = sqlc.arg(invitee_id) AND u.deleted_at IS NULL
),
new_participant AS (
    INSERT INTO participants (id, conversation_id, user_id, created_at)
    SELECT sqlc.arg(participant_id), vc.conv_id, vi.user_id, NOW()
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
    SELECT 1 FROM conversations c
    JOIN participants p ON p.conversation_id = c.id AND p.user_id = sqlc.arg(user_id)
    WHERE c.id = sqlc.arg(conversation_id)
      AND c.type = 0
      AND c.owner_id = sqlc.arg(user_id)
      AND c.deleted_at IS NULL
      AND p.deleted_at IS NULL
);

-- Message queries

-- name: StoreMessageAndReturn :one
WITH new_message AS (
    INSERT INTO messages (id, conversation_id, user_id, content, type, created_at)
    VALUES ($1, $2, $3, $4, $5, NOW())
    RETURNING id, type, created_at, conversation_id, content, user_id
)
SELECT
    nm.id, nm.type, nm.created_at, nm.conversation_id,
    nm.content as formatted_text,
    u.id as user_id, u.name as user_name, u.avatar as user_avatar
FROM new_message nm
JOIN users u ON u.id = nm.user_id
WHERE u.deleted_at IS NULL;
