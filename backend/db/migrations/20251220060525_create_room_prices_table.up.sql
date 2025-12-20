CREATE TABLE IF NOT EXISTS room_prices (
  id SERIAL PRIMARY KEY,
  room_id INT REFERENCES rooms(id),
  base_price DECIMAL(10,2) NOT NULL,
  platform_fee_percent DECIMAL(5,2) DEFAULT 10.00,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);