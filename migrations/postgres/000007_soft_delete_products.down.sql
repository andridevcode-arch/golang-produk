-- 000007_soft_delete_products.down.sql

DROP INDEX IF EXISTS idx_products_is_deleted;
ALTER TABLE products DROP COLUMN IF EXISTS is_deleted;
