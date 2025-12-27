-- Migration to drop 'caption' column and add 'post_date' column to 'posts' table
ALTER TABLE posts DROP COLUMN IF EXISTS caption;
ALTER TABLE posts ADD COLUMN post_date TIMESTAMP WITH TIME ZONE;

