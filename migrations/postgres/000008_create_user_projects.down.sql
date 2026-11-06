-- 000008_create_user_projects.down.sql

DROP INDEX IF EXISTS idx_user_projects_template_version_id;
DROP INDEX IF EXISTS idx_user_projects_template_id;
DROP INDEX IF EXISTS idx_user_projects_user_status;
DROP INDEX IF EXISTS idx_user_projects_user_created_at;
DROP TABLE IF EXISTS user_projects;
