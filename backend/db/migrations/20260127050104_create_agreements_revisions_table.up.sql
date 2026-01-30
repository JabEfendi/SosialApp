CREATE TABLE IF NOT EXISTS agreement_revisions (
    id SERIAL PRIMARY KEY,
    agreement_id INTEGER NOT NULL REFERENCES agreements(id) ON DELETE CASCADE,
    requested_by VARCHAR(20) NOT NULL CHECK (requested_by IN ('user', 'corporate')),
    note TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);