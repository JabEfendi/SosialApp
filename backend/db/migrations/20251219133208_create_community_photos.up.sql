CREATE TABLE IF NOT EXISTS community_photos (
    id SERIAL PRIMARY KEY,
    community_id INT NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    photo TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);