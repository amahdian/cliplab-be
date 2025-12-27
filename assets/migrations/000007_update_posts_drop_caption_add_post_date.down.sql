-- Rollback migration: add 'caption' column and drop 'post_date' column from 'posts' table
ALTER TABLE posts ADD COLUMN caption TEXT;
ALTER TABLE posts DROP COLUMN IF EXISTS post_date;

