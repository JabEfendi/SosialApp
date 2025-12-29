CREATE TABLE IF NOT EXISTS admin_audit_logs (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    target_id INTEGER NULL,
    before_data JSON NULL,
    after_data JSON NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);