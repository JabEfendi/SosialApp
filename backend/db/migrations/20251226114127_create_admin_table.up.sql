CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES admin_roles(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(100) DEFAULT 'pending',
    -- is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP NULL,
    approved_at TIMESTAMP NULL,
    approved_by INTEGER NULL REFERENCES admins(id),
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO admins (role_id, name, email, password, status) 
VALUES (1, 'Admin Super', 'admin@gmail.com', '$2a$12$u6m7w7bKQ8P7x9W2zY5r8e7cC5Z7nG9QzR1LqQZ0oF4kVnYwX5YpK', 'active');