CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";

CREATE TABLE IF NOT EXISTS posts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid,
    user_ip INET,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    link TEXT NOT NULL,
    caption TEXT,
    image_url TEXT,
    video_url TEXT,
    platform VARCHAR(50) NOT NULL DEFAULT 'instagram',
    format VARCHAR(50) NOT NULL DEFAULT 'image',
    user_name varchar(255) NOT NULL,
    user_anchor VARCHAR(255) NOT NULL,
    user_profile_link VARCHAR(255),
    user_profile_image VARCHAR(1024)
);

CREATE TABLE IF NOT EXISTS post_contents (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    type TEXT NOT NULL DEFAULT 'caption',
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    text TEXT,
    metadata TEXT NOT NULL DEFAULT '{}',
    embedding vector(1536)
);

CREATE INDEX IF NOT EXISTS idx_post_contents_post_id ON post_contents(post_id);

