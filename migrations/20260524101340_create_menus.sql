-- +goose Up
-- Create ENUM type first
CREATE TYPE menus_type_enum AS ENUM ('catalog', 'menus');

CREATE TABLE menus (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    remark JSONB NULL,
    menus_type menus_type_enum NULL,  -- Using PostgreSQL ENUM
    icon VARCHAR(100) NULL,
    path VARCHAR(100) NOT NULL,
    order_number INT NOT NULL DEFAULT 0,
    parent_id INT NOT NULL DEFAULT 0,
    is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    component VARCHAR(100) NOT NULL,
    keepalive BOOLEAN NOT NULL DEFAULT TRUE,
    redirect VARCHAR(100) NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_menus_name ON menus(name);
CREATE INDEX idx_menus_path ON menus(path);
CREATE INDEX idx_menus_order ON menus(order_number);
CREATE INDEX idx_menus_parent_id ON menus(parent_id);

-- +goose Down
DROP TABLE IF EXISTS menus;
DROP TYPE IF EXISTS menus_type_enum;