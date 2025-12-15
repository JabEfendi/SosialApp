CREATE TABLE IF NOT EXISTS community_files (
    id SERIAL PRIMARY KEY,
    community_id INT REFERENCES communities(id) ON DELETE CASCADE,
    uploaded_by INT REFERENCES users(id),
    file_url TEXT NOT NULL,
    file_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT now()
);
