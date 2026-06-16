-- +goose Up
-- Create roles table (plural name)
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    description VARCHAR(500) NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create role_menus junction table
CREATE TABLE role_menus (
    role_id BIGINT NOT NULL,
    menus_id BIGINT NOT NULL,

    PRIMARY KEY (role_id, menus_id),

    CONSTRAINT fk_role_menus_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_role_menus_menus
        FOREIGN KEY (menus_id)
        REFERENCES menus(id)
        ON DELETE CASCADE
);

-- Create role_apis junction table
CREATE TABLE role_apis (
    role_id BIGINT NOT NULL,
    apis_id BIGINT NOT NULL,

    PRIMARY KEY (role_id, apis_id),

    CONSTRAINT fk_role_apis_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_role_apis_apis
        FOREIGN KEY (apis_id)
        REFERENCES apis(id)
        ON DELETE CASCADE
);

-- Create indexes separately
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_role_menus_role_id ON role_menus(role_id);
CREATE INDEX idx_role_menus_menus_id ON role_menus(menus_id);
CREATE INDEX idx_role_apis_role_id ON role_apis(role_id);
CREATE INDEX idx_role_apis_apis_id ON role_apis(apis_id);

-- +goose Down
-- Drop indexes first (good practice)
DROP INDEX IF EXISTS idx_role_apis_apis_id;
DROP INDEX IF EXISTS idx_role_apis_role_id;
DROP INDEX IF EXISTS idx_role_menus_menus_id;
DROP INDEX IF EXISTS idx_role_menus_role_id;
DROP INDEX IF EXISTS idx_roles_name;

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS role_apis;
DROP TABLE IF EXISTS role_menus;
DROP TABLE IF EXISTS roles;