CREATE TABLE IF NOT EXISTS room_commissions (
  id SERIAL PRIMARY KEY,
  room_id INT REFERENCES rooms(id),
  order_id INT REFERENCES room_join_orders(id),
  platform_fee DECIMAL(10,2) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);