CREATE TABLE IF NOT EXISTS communities (
    id SERIAL PRIMARY KEY,
    creator_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    country_region VARCHAR(150),
    interests VARCHAR(255),
    type VARCHAR(20) NOT NULL CHECK (type IN ('public', 'private')),
    auto_approve BOOLEAN DEFAULT false,
    chat_notifications BOOLEAN DEFAULT true,
    invite_code VARCHAR(10) UNIQUE NOT NULL,
    cover_image TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
