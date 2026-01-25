CREATE INDEX IF NOT EXISTS messages_conversation_created_id_idx
ON messages (conversation_id, created_at DESC, id DESC)
WHERE deleted_at IS NULL;
