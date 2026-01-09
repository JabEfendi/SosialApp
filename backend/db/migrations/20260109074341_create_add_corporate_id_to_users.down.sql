ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_corporate_id_fkey;

ALTER TABLE users
DROP COLUMN IF EXISTS corporate_id,
DROP COLUMN IF EXISTS joined_corporate_at;