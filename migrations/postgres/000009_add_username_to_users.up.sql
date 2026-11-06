-- 000009_add_username_to_users.up.sql

ALTER TABLE users ADD COLUMN username VARCHAR(50);
CREATE UNIQUE INDEX uq_users_username_lower ON users (LOWER(username)) WHERE username IS NOT NULL;
