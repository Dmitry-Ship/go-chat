-- Combined schema for sqlc code generation
-- This file is auto-generated from migration files

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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_conversations_deleted_at ON conversations(deleted_at);

CREATE TABLE group_conversations (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    avatar TEXT,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_group_conversations_conversation_id ON group_conversations(conversation_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_group_conversations_owner_id ON group_conversations(owner_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_group_conversations_deleted_at ON group_conversations(deleted_at);

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
