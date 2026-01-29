-- Up migration: create the links table
CREATE TABLE links (
    id TEXT PRIMARY KEY,
    short_code TEXT UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_links_short_code ON links (short_code);