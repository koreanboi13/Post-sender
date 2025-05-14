CREATE TABLE IF NOT EXISTS chat_entries (
	id VARCHAR(255) PRIMARY KEY,
	messenger VARCHAR(50) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_chat_entries_messenger ON chat_entries(messenger);