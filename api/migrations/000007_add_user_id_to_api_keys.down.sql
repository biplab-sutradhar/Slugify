ALTER TABLE api_keys DROP CONSTRAINT IF EXISTS fk_api_keys_user_id;
ALTER TABLE api_keys DROP COLUMN IF EXISTS user_id;
