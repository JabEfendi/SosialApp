CREATE TABLE IF NOT EXISTS corporates (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(50) DEFAULT 'active',
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);