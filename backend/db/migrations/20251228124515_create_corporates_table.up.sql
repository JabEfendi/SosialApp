CREATE TABLE IF NOT EXISTS corporates (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    logo VARCHAR(255),
    reffcorporate VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active',
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
