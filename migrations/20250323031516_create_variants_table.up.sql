CREATE TABLE variants (
  id SERIAL PRIMARY KEY,
  printful_id BIGINT UNIQUE NOT NULL,
  product_id INTEGER NOT NULL REFERENCES products(id),
  name TEXT,
  sku TEXT,
  retail_price TEXT,
  thumbnail_url TEXT,
  size TEXT,
  color TEXT,
  availability_status TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);
