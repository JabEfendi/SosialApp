CREATE TABLE IF NOT EXISTS login_logs (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  ip_address VARCHAR(100),
  device VARCHAR(255),
  location VARCHAR(255),
  user_agent TEXT,
  logged_in_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
