CREATE TABLE links (
    id TEXT PRIMARY KEY,
    short_code TEXT UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_short_code ON links (short_code);