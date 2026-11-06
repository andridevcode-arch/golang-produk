-- 000006_soft_delete_users.up.sql
ALTER TABLE users
ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX idx_users_is_deleted ON users (is_deleted);