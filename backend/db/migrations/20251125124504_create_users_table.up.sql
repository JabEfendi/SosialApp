CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  username VARCHAR(50) UNIQUE,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  gender VARCHAR(20) NOT NULL,
  birthdate DATE,
  phone VARCHAR(20),
  bio VARCHAR(255),
  country VARCHAR(100),
  address TEXT,
  provider VARCHAR(50) NOT NULL DEFAULT 'local',
  provider_id VARCHAR(255),
  avatar TEXT,
  FCMToken VARCHAR,
  -- coin_balance BIGINT NOT NULL DEFAULT 0,
  referral_code VARCHAR(20) UNIQUE,
  referred_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

INSERT INTO users (id, name, username, email, password, gender, provider)
VALUES (0, 'Viscata Pancen Uye', 'viscata', 'hr.viscata@gmail.com', '$2a$12$N9qo8uLOickgx2ZMRZo5i.ezvFjH8xGq5c6GJ6ZyM5u0X6f3mHk2K', 'male', 'local');