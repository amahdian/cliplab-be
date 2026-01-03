DROP INDEX IF EXISTS idx_posts_channel_id;
DROP INDEX IF EXISTS idx_channel_histories_channel_id;
DROP INDEX IF EXISTS idx_post_analyses_post_id;

DROP TABLE IF EXISTS post_analyses;

ALTER TABLE posts DROP COLUMN IF EXISTS video_play_count;
ALTER TABLE posts DROP COLUMN IF EXISTS video_view_count;
ALTER TABLE posts DROP COLUMN IF EXISTS comment_count;
ALTER TABLE posts DROP COLUMN IF EXISTS like_count;
ALTER TABLE posts DROP COLUMN IF EXISTS channel_id;

DROP TABLE IF EXISTS channel_histories;
DROP TABLE IF EXISTS channels;
