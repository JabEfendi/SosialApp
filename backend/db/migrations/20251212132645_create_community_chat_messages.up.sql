CREATE TABLE IF NOT EXISTS community_chat_messages (
    id SERIAL PRIMARY KEY,
    community_id INT NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);