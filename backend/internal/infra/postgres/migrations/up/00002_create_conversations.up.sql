-- +goose Up
-- +goose StatementBegin
CREATE TABLE conversations (
    id UUID PRIMARY KEY,
    type INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_conversations_deleted_at ON conversations(deleted_at);
-- +goose StatementEnd
