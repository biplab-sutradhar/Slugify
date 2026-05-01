ALTER TABLE links ADD COLUMN user_id TEXT;

ALTER TABLE links ADD CONSTRAINT fk_links_user
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_links_user_id ON links(user_id);