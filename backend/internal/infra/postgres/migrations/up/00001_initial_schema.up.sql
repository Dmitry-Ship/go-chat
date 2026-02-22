-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY,
    avatar TEXT,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    refresh_token TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_name ON users(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE TABLE conversations (
    id UUID PRIMARY KEY,
    type INTEGER NOT NULL,
    name TEXT,
    avatar TEXT,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT conversations_type_chk CHECK (type IN (0, 1)),
    CONSTRAINT conversations_group_metadata_chk CHECK (
        (type = 0 AND owner_id IS NOT NULL AND name IS NOT NULL)
        OR (type = 1 AND owner_id IS NULL)
    )
);

CREATE INDEX idx_conversations_deleted_at ON conversations(deleted_at);
CREATE INDEX idx_conversations_owner_id ON conversations(owner_id) WHERE deleted_at IS NULL;

CREATE TABLE participants (
    id UUID PRIMARY KEY,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_participants_conversation_user ON participants(conversation_id, user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_participants_conversation_id ON participants(conversation_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_participants_user_id ON participants(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_participants_deleted_at ON participants(deleted_at);

CREATE TABLE messages (
    id UUID PRIMARY KEY,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    type INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_messages_conversation_id ON messages(conversation_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_messages_user_id ON messages(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_messages_created_at ON messages(created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_messages_deleted_at ON messages(deleted_at);
CREATE INDEX messages_conversation_created_id_idx
ON messages (conversation_id, created_at DESC, id DESC)
WHERE deleted_at IS NULL;
-- +goose StatementEnd
