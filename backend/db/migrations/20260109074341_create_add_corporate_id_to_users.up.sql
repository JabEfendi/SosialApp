ALTER TABLE users
ADD COLUMN corporate_id BIGINT NULL,
ADD COLUMN joined_corporate_at TIMESTAMP NULL;

ALTER TABLE users
ADD CONSTRAINT users_corporate_id_fkey
FOREIGN KEY (corporate_id)
REFERENCES corporates(id)
ON DELETE SET NULL;