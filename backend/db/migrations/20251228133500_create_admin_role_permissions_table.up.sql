CREATE TABLE IF NOT EXISTS admin_role_permissions (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES admin_roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES admin_permissions(id) ON DELETE CASCADE,
    
    UNIQUE (role_id, permission_id)
);
INSERT INTO admin_role_permissions (role_id, permission_id)
SELECT 1, id FROM admin_permissions;

