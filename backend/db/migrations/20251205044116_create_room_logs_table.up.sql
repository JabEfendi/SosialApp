CREATE TABLE IF NOT EXISTS room_logs (
  id SERIAL PRIMARY KEY,
  room_id INT REFERENCES rooms(id),
  user_id INT REFERENCES users(id),
  kyc_status VARCHAR(20),
  summary_json JSON,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);