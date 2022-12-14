CREATE TABLE IF NOT EXISTS transaction.products(
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT NOT NULL,
  price BIGINT NOT NULL,
  brand_id INT NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transaction_products_name ON transaction.products(name);
CREATE INDEX IF NOT EXISTS idx_transaction_products_brand_id ON transaction.products(brand_id);
