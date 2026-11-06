-- 000006_soft_delete_users.down.sql

DROP INDEX IF EXISTS idx_users_is_deleted;
ALTER TABLE users DROP COLUMN IF EXISTS is_deleted;
