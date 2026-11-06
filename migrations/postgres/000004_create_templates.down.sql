-- 000004_create_templates.down.sql
DROP INDEX IF EXISTS idx_templates_is_premium;

DROP INDEX IF EXISTS idx_templates_is_active;

DROP INDEX IF EXISTS idx_templates_product_id;

DROP TABLE IF EXISTS templates;