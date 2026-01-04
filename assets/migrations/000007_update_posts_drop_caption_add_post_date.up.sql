-- Migration to add 'post_date' column to 'posts' table
ALTER TABLE posts ADD COLUMN post_date TIMESTAMPTZ;

