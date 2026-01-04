CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";

CREATE TABLE IF NOT EXISTS posts (
    id VARCHAR(255) PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    user_ip INET,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    link TEXT NOT NULL,
    image_url TEXT,
    video_url TEXT,
    format VARCHAR(50) NOT NULL DEFAULT 'image',
    user_name varchar(255) NOT NULL,
    user_anchor VARCHAR(255) NOT NULL,
    user_profile_link VARCHAR(255),
    user_profile_image VARCHAR(1024)
);

CREATE TABLE IF NOT EXISTS post_contents (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id VARCHAR(255) NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    type TEXT NOT NULL DEFAULT 'caption',
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    text TEXT,
    metadata TEXT NOT NULL DEFAULT '{}',
    embedding vector(1536)
);

CREATE INDEX IF NOT EXISTS idx_post_contents_post_id ON post_contents(post_id);

