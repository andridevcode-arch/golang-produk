-- 000009_add_username_to_users.down.sql

DROP INDEX IF EXISTS uq_users_username_lower;
ALTER TABLE users DROP COLUMN username;
