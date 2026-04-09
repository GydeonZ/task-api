CREATE TABLE IF NOT EXISTS boards (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_boards_deleted_at ON boards(deleted_at);
CREATE INDEX IF NOT EXISTS idx_boards_user_id ON boards(user_id);
