CREATE TABLE IF NOT EXISTS room_owner_earnings (
  id SERIAL PRIMARY KEY,
  room_id INT REFERENCES rooms(id),
  owner_id INT REFERENCES users(id),
  order_id INT REFERENCES room_join_orders(id),
  amount DECIMAL(10,2) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);