DROP TABLE IF EXISTS transaction.products;
DROP INDEX IF EXISTS idx_transaction_products_name ON transaction.products(name);
DROP INDEX IF EXISTS idx_transaction_products_brand_id ON transaction.products(brand_id);
