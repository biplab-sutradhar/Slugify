CREATE TABLE IF NOT EXISTS click_events (
    id SERIAL PRIMARY KEY,
    link_id TEXT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    short_code TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    ip_hash TEXT,
    user_agent TEXT,
    referrer TEXT,
    country TEXT DEFAULT 'unknown'
);

CREATE INDEX IF NOT EXISTS idx_click_events_link_id ON click_events (link_id);
CREATE INDEX IF NOT EXISTS idx_click_events_short_code ON click_events (short_code);
CREATE INDEX IF NOT EXISTS idx_click_events_timestamp ON click_events (timestamp);