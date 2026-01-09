-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_conversations_deleted_at;
DROP INDEX IF EXISTS idx_conversations_is_active;
DROP TABLE IF EXISTS conversations;
-- +goose StatementEnd
