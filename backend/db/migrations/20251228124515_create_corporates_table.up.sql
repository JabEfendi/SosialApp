CREATE TABLE corporates (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(50) DEFAULT 'active',
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
ALTER TABLE users
ADD COLUMN corporate_id BIGINT NULL REFERENCES corporates(id) ON DELETE SET NULL,
ADD COLUMN joined_corporate_at TIMESTAMP NULL;