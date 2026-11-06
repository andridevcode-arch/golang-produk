-- 000005_create_template_versions.down.sql

DROP INDEX IF EXISTS uq_template_versions_active_template;
DROP INDEX IF EXISTS idx_template_versions_template_id;
DROP TABLE IF EXISTS template_versions;
