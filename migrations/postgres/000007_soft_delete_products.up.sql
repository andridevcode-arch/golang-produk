-- 000007_soft_delete_products.up.sql

ALTER TABLE products ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT false;
CREATE INDEX idx_products_is_deleted ON products(is_deleted);
