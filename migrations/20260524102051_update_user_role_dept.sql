-- +goose Up

ALTER TABLE users
    ADD COLUMN username VARCHAR(20) NOT NULL,
    ADD COLUMN alias VARCHAR(30) NULL,
    -- ADD COLUMN email VARCHAR(255) NOT NULL,
    ADD COLUMN phone VARCHAR(20) NULL,
    -- ADD COLUMN password VARCHAR(128) NULL,
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN is_superuser BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN last_login DATETIME NULL,
    ADD COLUMN dept_id INT NULL;

-- indexes
CREATE UNIQUE INDEX idx_user_username ON users(username);
CREATE UNIQUE INDEX idx_user_email ON users(email);

CREATE INDEX idx_user_alias ON users(alias);
CREATE INDEX idx_user_phone ON users(phone);
CREATE INDEX idx_user_is_active ON users(is_active);
CREATE INDEX idx_user_is_superuser ON users(is_superuser);
CREATE INDEX idx_user_last_login ON users(last_login);
CREATE INDEX idx_user_dept_id ON users(dept_id);

CREATE TABLE user_role (
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (user_id, role_id),

    CONSTRAINT fk_user_role_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_role_role
        FOREIGN KEY (role_id)
        REFERENCES role(id)
        ON DELETE CASCADE
);