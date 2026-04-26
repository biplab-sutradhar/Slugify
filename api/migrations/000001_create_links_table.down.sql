-- Down migration: rollback the links table
DROP INDEX IF EXISTS idx_links_short_code;
DROP TABLE IF EXISTS links;