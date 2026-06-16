-- +goose Up
-- Add columns to users table
ALTER TABLE users
    ADD COLUMN username VARCHAR(20) NOT NULL,
    ADD COLUMN alias VARCHAR(30) NULL,
    ADD COLUMN phone VARCHAR(20) NULL,
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN is_superuser BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN last_login TIMESTAMP NULL,
    ADD COLUMN dept_id INT NULL;

-- Create indexes separately
CREATE UNIQUE INDEX idx_users_username ON users(username);

CREATE INDEX idx_users_alias ON users(alias);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_is_superuser ON users(is_superuser);
CREATE INDEX idx_users_last_login ON users(last_login);
CREATE INDEX idx_users_dept_id ON users(dept_id);

-- Create user_role junction table
CREATE TABLE user_role (
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,

    PRIMARY KEY (user_id, role_id),

    CONSTRAINT fk_user_role_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_role_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
);

-- Create indexes for the junction table
CREATE INDEX idx_user_role_user_id ON user_role(user_id);
CREATE INDEX idx_user_role_role_id ON user_role(role_id);

-- +goose Down
-- Drop indexes from junction table
DROP INDEX IF EXISTS idx_user_role_role_id;
DROP INDEX IF EXISTS idx_user_role_user_id;

-- Drop junction table
DROP TABLE IF EXISTS user_role;

-- Drop indexes from users table
DROP INDEX IF EXISTS idx_users_dept_id;
DROP INDEX IF EXISTS idx_users_last_login;
DROP INDEX IF EXISTS idx_users_is_superuser;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_alias;
DROP INDEX IF EXISTS idx_users_username;

-- Remove columns from users table
ALTER TABLE users
    DROP COLUMN IF EXISTS username,
    DROP COLUMN IF EXISTS alias,
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS is_superuser,
    DROP COLUMN IF EXISTS last_login,
    DROP COLUMN IF EXISTS dept_id;