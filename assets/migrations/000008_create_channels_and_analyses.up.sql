-- Create channels table
CREATE TABLE IF NOT EXISTS channels (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name TEXT NOT NULL,
    handler TEXT NOT NULL,
    platform VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create channel_histories table
CREATE TABLE IF NOT EXISTS channel_histories (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id uuid NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    followers_count INTEGER NOT NULL DEFAULT 0,
    following_count INTEGER NOT NULL DEFAULT 0,
    media_count INTEGER NOT NULL DEFAULT 0,
    average_likes INTEGER NOT NULL DEFAULT 0,
    average_comments INTEGER NOT NULL DEFAULT 0,
    average_video_views INTEGER NOT NULL DEFAULT 0,
    average_video_plays INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add missing columns to posts
ALTER TABLE posts ADD COLUMN IF NOT EXISTS channel_id uuid REFERENCES channels(id) ON DELETE SET NULL;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS like_count BIGINT NOT NULL DEFAULT 0;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS comment_count BIGINT NOT NULL DEFAULT 0;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS video_view_count BIGINT NOT NULL DEFAULT 0;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS video_play_count BIGINT NOT NULL DEFAULT 0;

-- Create post_analyses table
CREATE TABLE IF NOT EXISTS post_analyses (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id VARCHAR(255) NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    viral_score INTEGER NOT NULL DEFAULT 0,
    big_idea TEXT,
    why_viral TEXT,
    audience_sentiment TEXT,
    sentiment_score INTEGER NOT NULL DEFAULT 0,
    metrics JSONB NOT NULL DEFAULT '[]',
    strengths JSONB NOT NULL DEFAULT '[]',
    weaknesses JSONB NOT NULL DEFAULT '[]',
    hook_ideas JSONB NOT NULL DEFAULT '[]',
    script_ideas JSONB NOT NULL DEFAULT '[]',
    captions JSONB NOT NULL DEFAULT '{}',
    hashtags JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_post_analyses_post_id ON post_analyses(post_id);
CREATE INDEX IF NOT EXISTS idx_channel_histories_channel_id ON channel_histories(channel_id);
CREATE INDEX IF NOT EXISTS idx_posts_channel_id ON posts(channel_id);
