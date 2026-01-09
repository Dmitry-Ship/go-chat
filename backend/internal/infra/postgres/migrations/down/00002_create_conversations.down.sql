-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_conversations_deleted_at;
DROP TABLE IF EXISTS conversations;
-- +goose StatementEnd
