CREATE TABLE IF NOT EXISTS community_members (
    id SERIAL PRIMARY KEY,
    community_id INT NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member', -- admin, moderator, member
    status VARCHAR(20) DEFAULT 'approved', -- pending / approved
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    UNIQUE (community_id, user_id)
);
