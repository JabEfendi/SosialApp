CREATE TABLE IF NOT EXISTS email_campaigns (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    subject VARCHAR(150) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    scheduled_at TIMESTAMP,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);