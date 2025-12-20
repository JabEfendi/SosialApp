CREATE TABLE IF NOT EXISTS room_join_orders (
  id SERIAL PRIMARY KEY,
  order_code VARCHAR(50) UNIQUE,
  room_id INT REFERENCES rooms(id),
  user_id INT REFERENCES users(id),
  base_price DECIMAL(10,2) NOT NULL,
  platform_fee DECIMAL(10,2) NOT NULL,
  total_price DECIMAL(10,2) NOT NULL,
  status VARCHAR(20) DEFAULT 'pending',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  paid_at TIMESTAMP WITH TIME ZONE
);