-- Down migration: rollback the ranges table
DROP INDEX IF EXISTS idx_ranges_active;
DROP TABLE IF EXISTS ranges;