CREATE TABLE user_legal_acceptances (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    legal_document_id INTEGER NOT NULL,
    accepted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (user_id, legal_document_id)
);