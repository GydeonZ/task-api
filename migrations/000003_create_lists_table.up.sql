CREATE TABLE IF NOT EXISTS lists (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    title VARCHAR(255) NOT NULL,
    position INTEGER DEFAULT 0,
    board_id INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_lists_deleted_at ON lists(deleted_at);
CREATE INDEX IF NOT EXISTS idx_lists_board_id ON lists(board_id);
