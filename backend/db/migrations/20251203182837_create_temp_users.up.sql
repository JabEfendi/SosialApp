CREATE TABLE IF NOT EXISTS temp_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255),
    gender VARCHAR(50),
    birthdate DATE,
    phone VARCHAR(50),
    bio TEXT,
    country VARCHAR(100),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
