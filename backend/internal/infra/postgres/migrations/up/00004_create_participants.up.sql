-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd
