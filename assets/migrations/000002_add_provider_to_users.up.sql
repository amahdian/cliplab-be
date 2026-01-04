ALTER TABLE users ADD COLUMN provider VARCHAR(50) DEFAULT 'local';
ALTER TABLE users ADD COLUMN provider_id VARCHAR(255);
CREATE INDEX idx_users_provider_id ON users(provider_id);
