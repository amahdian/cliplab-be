ALTER TABLE posts DROP COLUMN IF EXISTS user_id;
ALTER TABLE posts DROP COLUMN IF EXISTS user_ip;
ALTER TABLE posts DROP COLUMN IF EXISTS status;
ALTER TABLE posts DROP COLUMN IF EXISTS fail_reason;

ALTER TABLE post_contents DROP COLUMN IF EXISTS status;
ALTER TABLE post_contents DROP COLUMN IF EXISTS fail_reason;

CREATE TABLE IF NOT EXISTS analyze_requests
(
    id           uuid PRIMARY KEY      DEFAULT uuid_generate_v4(),
    user_id      uuid         REFERENCES users (id) ON DELETE SET NULL,
    user_ip      INET,
    post_id      VARCHAR(255) NOT NULL REFERENCES posts (id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    link         TEXT         NOT NULL,
    status       VARCHAR(100) NOT NULL DEFAULT 'pending',
    fail_reason  TEXT,
    llm_request  TEXT,
    llm_response TEXT
);

CREATE INDEX IF NOT EXISTS idx_analyze_requests_user_id ON analyze_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_analyze_requests_post_id ON analyze_requests(post_id);