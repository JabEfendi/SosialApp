CREATE TABLE admin_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO admin_roles (name, description) VALUES
('superadmin', 'Full access to all admin features'),
('admin', 'Limited admin access');