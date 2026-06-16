-- +goose Up
-- Create ENUM type for HTTP methods
CREATE TYPE http_method_enum AS ENUM ('GET', 'POST', 'PUT', 'DELETE', 'PATCH');

CREATE TABLE apis (
    id BIGSERIAL PRIMARY KEY,
    path VARCHAR(100) NOT NULL,
    method http_method_enum NOT NULL,
    summary VARCHAR(500) NOT NULL,
    tags VARCHAR(100) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes separately
CREATE INDEX idx_apis_path ON apis(path);
CREATE INDEX idx_apis_method ON apis(method);
CREATE INDEX idx_apis_summary ON apis(summary);
CREATE INDEX idx_apis_tags ON apis(tags);

-- +goose Down
DROP TABLE IF EXISTS apis;
DROP TYPE IF EXISTS http_method_enum;