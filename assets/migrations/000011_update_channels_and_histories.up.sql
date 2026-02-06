-- Add columns to channels table
ALTER TABLE channels ADD COLUMN IF NOT EXISTS profile_image TEXT;
ALTER TABLE channels ADD COLUMN IF NOT EXISTS profile_descriptor TEXT;

-- Update channel_histories table
ALTER TABLE channel_histories ADD COLUMN IF NOT EXISTS profile_image TEXT;
ALTER TABLE channel_histories ADD COLUMN IF NOT EXISTS profile_descriptor TEXT;
ALTER TABLE channel_histories ADD COLUMN IF NOT EXISTS average_engagement_rate DOUBLE PRECISION;
ALTER TABLE channel_histories ADD COLUMN IF NOT EXISTS latest_post_publish_date TIMESTAMPTZ;

-- Change column types to DOUBLE PRECISION for averages
ALTER TABLE channel_histories ALTER COLUMN average_likes TYPE DOUBLE PRECISION;
ALTER TABLE channel_histories ALTER COLUMN average_comments TYPE DOUBLE PRECISION;
ALTER TABLE channel_histories ALTER COLUMN average_video_views TYPE DOUBLE PRECISION;
ALTER TABLE channel_histories ALTER COLUMN average_video_plays TYPE DOUBLE PRECISION;
