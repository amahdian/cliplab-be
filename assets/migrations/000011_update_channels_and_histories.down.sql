-- Revert column types to INTEGER
ALTER TABLE channel_histories ALTER COLUMN average_likes TYPE INTEGER;
ALTER TABLE channel_histories ALTER COLUMN average_comments TYPE INTEGER;
ALTER TABLE channel_histories ALTER COLUMN average_video_views TYPE INTEGER;
ALTER TABLE channel_histories ALTER COLUMN average_video_plays TYPE INTEGER;

-- Remove added columns from channel_histories
ALTER TABLE channel_histories DROP COLUMN IF EXISTS latest_post_publish_date;
ALTER TABLE channel_histories DROP COLUMN IF EXISTS average_engagement_rate;
ALTER TABLE channel_histories DROP COLUMN IF EXISTS profile_descriptor;
ALTER TABLE channel_histories DROP COLUMN IF EXISTS profile_image;

-- Remove added columns from channels
ALTER TABLE channels DROP COLUMN IF EXISTS profile_descriptor;
ALTER TABLE channels DROP COLUMN IF EXISTS profile_image;
