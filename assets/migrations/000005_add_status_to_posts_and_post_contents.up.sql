ALTER TABLE posts ADD COLUMN status VARCHAR(100) NOT NULL DEFAULT 'pending';
ALTER TABLE posts ADD COLUMN fail_reason TEXT;

ALTER TABLE post_contents ADD COLUMN status TEXT NOT NULL DEFAULT 'pending';
ALTER TABLE post_contents ADD COLUMN fail_reason TEXT;

