-- Rollback migration: drop 'post_date' column from 'posts' table
ALTER TABLE posts DROP COLUMN IF EXISTS post_date;

