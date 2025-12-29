CREATE TABLE IF NOT EXISTS admin_permissions (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255)
);
INSERT INTO admin_permissions (code, description) VALUES
('system.permission.create', 'Create system permission'),
('system.role.permission.update', 'Update role permissions'),
('system.notification.view', 'View notification settings'),
('system.notification.update', 'Update notification settings'),
('system.legal.update', 'Update legal documents'),
('system.email.schedule', 'Schedule email campaigns'),
('system.maintenance.create', 'Create maintenance schedules');