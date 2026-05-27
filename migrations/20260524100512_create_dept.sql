-- +goose Up
CREATE TABLE dept (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    `desc` VARCHAR(500) NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    `order` INT NOT NULL DEFAULT 0,
    parent_id INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_dept_name (name),
    INDEX idx_dept_is_deleted (is_deleted),
    INDEX idx_dept_order (`order`),
    INDEX idx_dept_parent_id (parent_id)
);

-- +goose Down
DROP TABLE IF EXISTS dept;
