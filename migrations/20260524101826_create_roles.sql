-- +goose Up
CREATE TABLE role (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    `desc` VARCHAR(500) NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_role_name (name)
);

CREATE TABLE role_menu (
    role_id BIGINT UNSIGNED NOT NULL,
    menu_id BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (role_id, menu_id),

    CONSTRAINT fk_role_menu_role
        FOREIGN KEY (role_id)
        REFERENCES role(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_role_menu_menu
        FOREIGN KEY (menu_id)
        REFERENCES menu(id)
        ON DELETE CASCADE
);

CREATE TABLE role_api (
    role_id BIGINT UNSIGNED NOT NULL,
    api_id BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (role_id, api_id),

    CONSTRAINT fk_role_api_role
        FOREIGN KEY (role_id)
        REFERENCES role(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_role_api_api
        FOREIGN KEY (api_id)
        REFERENCES api(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS role_api;
DROP TABLE IF EXISTS role_menu;
DROP TABLE IF EXISTS role;