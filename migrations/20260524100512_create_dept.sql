-- +goose Up
CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    "desc" VARCHAR(500) NULL,  -- Using double quotes for reserved keyword
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    "order" INT NOT NULL DEFAULT 0,  -- Using double quotes for reserved keyword
    parent_id INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_departments_name ON departments(name);
CREATE INDEX idx_departments_is_deleted ON departments(is_deleted);
CREATE INDEX idx_departments_order ON departments("order");
CREATE INDEX idx_departments_parent_id ON departments(parent_id);

-- +goose Down
DROP TABLE IF EXISTS departments;