-- Remove verdict column from post_analyses
ALTER TABLE post_analyses DROP COLUMN IF EXISTS verdict;
