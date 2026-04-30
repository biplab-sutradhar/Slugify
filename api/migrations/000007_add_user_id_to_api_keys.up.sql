ALTER TABLE api_keys ADD COLUMN user_id TEXT;

ALTER TABLE api_keys ADD CONSTRAINT fk_api_keys_user_id 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys (user_id);
