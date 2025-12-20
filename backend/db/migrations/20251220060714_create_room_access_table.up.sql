CREATE TABLE IF NOT EXISTS room_access (
  id SERIAL PRIMARY KEY,
  room_id INT REFERENCES rooms(id),
  user_id INT REFERENCES users(id),
  order_id INT REFERENCES room_join_orders(id),
  granted_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);