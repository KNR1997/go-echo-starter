-- +goose Up
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,  -- Explicitly set default to NULL
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_posts_title ON posts(title);

-- +goose Down
DROP TABLE IF EXISTS posts;
