    CREATE TABLE user_reports (
        id SERIAL PRIMARY KEY,
        reporter_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        reported_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        reason_id INT NOT NULL REFERENCES report_reasons(id) ON DELETE CASCADE,
        description TEXT,
        status VARCHAR(20) DEFAULT 'pending',
        admin_note TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );