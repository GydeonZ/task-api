CREATE TABLE IF NOT EXISTS cards (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    position INTEGER DEFAULT 0,
    due_date TIMESTAMPTZ,
    list_id INTEGER NOT NULL REFERENCES lists(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_cards_deleted_at ON cards(deleted_at);
CREATE INDEX IF NOT EXISTS idx_cards_list_id ON cards(list_id);
