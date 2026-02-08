CREATE TABLE IF NOT EXISTS usage_histories
(
    id           uuid PRIMARY KEY      DEFAULT uuid_generate_v4(),
    user_id      uuid         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    tool_name         VARCHAR(255) NOT NULL,
    credits_used      INTEGER      NOT NULL DEFAULT 0,
    remaining_credits INTEGER      NOT NULL DEFAULT 0,
    status            VARCHAR(100) NOT NULL DEFAULT 'completed',
    details      TEXT,
    platform     VARCHAR(100) NOT NULL DEFAULT 'general',
    type              VARCHAR(100) NOT NULL DEFAULT 'analysis',
    request_id        uuid         REFERENCES analyze_requests (id) ON DELETE SET NULL,
    response_log      TEXT,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_usage_histories_user_id ON usage_histories(user_id);
CREATE INDEX IF NOT EXISTS idx_usage_histories_created_at ON usage_histories(created_at);
