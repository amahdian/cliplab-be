-- Migration to add 'post_date' column to 'posts' table
ALTER TABLE posts ADD COLUMN IF NOT EXISTS post_date TIMESTAMPTZ;

