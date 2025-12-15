CREATE TABLE IF NOT EXISTS community_announcements (
    id SERIAL PRIMARY KEY,
    community_id INT NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_by INT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
