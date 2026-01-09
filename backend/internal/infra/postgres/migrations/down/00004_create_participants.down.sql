-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_participants_deleted_at;
DROP INDEX IF EXISTS idx_participants_user_id;
DROP INDEX IF EXISTS idx_participants_conversation_id;
DROP INDEX IF EXISTS idx_participants_conversation_user;
DROP TABLE IF EXISTS participants;
-- +goose StatementEnd
