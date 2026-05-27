-- +goose Up
CREATE TABLE menu (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    remark JSON NULL,
    menu_type ENUM('catalog', 'menu') NULL,
    icon VARCHAR(100) NULL,
    path VARCHAR(100) NOT NULL,
    `order` INT NOT NULL DEFAULT 0,
    parent_id INT NOT NULL DEFAULT 0,
    is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    component VARCHAR(100) NOT NULL,
    keepalive BOOLEAN NOT NULL DEFAULT TRUE,
    redirect VARCHAR(100) NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_menu_name (name),
    INDEX idx_menu_path (path),
    INDEX idx_menu_order (`order`),
    INDEX idx_menu_parent_id (parent_id)
);

-- +goose Down
DROP TABLE IF EXISTS menu;