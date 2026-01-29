-- 000002_create_ranges_table.up.sql
CREATE TABLE IF NOT EXISTS ranges (
    range_id SERIAL PRIMARY KEY,
    start_id BIGINT NOT NULL,
    end_id BIGINT NOT NULL,
    current_id BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_ranges_active ON ranges (is_active);