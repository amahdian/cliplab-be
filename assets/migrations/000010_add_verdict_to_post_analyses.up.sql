-- Add verdict column to post_analyses
ALTER TABLE post_analyses ADD COLUMN IF NOT EXISTS verdict JSONB NOT NULL DEFAULT '{}';
