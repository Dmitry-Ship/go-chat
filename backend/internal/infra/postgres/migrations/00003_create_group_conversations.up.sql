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
