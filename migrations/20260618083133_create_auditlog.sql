-- +goose Up
-- Create ENUM type for HTTP methods
-- CREATE TYPE http_method_enum AS ENUM ('GET', 'POST', 'PUT', 'DELETE', 'PATCH');

-- Create audit_logs table
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    username VARCHAR(64) NOT NULL DEFAULT '',
    module VARCHAR(64) NOT NULL DEFAULT '',
    summary VARCHAR(128) NOT NULL DEFAULT '',
    method http_method_enum NOT NULL,
    path VARCHAR(255) NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT -1,
    response_time INTEGER NOT NULL DEFAULT 0,
    request_args JSONB,
    response_body JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_username ON audit_logs(username);
CREATE INDEX idx_audit_logs_module ON audit_logs(module);
CREATE INDEX idx_audit_logs_summary ON audit_logs(summary);
CREATE INDEX idx_audit_logs_method ON audit_logs(method);
CREATE INDEX idx_audit_logs_path ON audit_logs(path);
CREATE INDEX idx_audit_logs_status ON audit_logs(status);
CREATE INDEX idx_audit_logs_response_time ON audit_logs(response_time);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Create composite index for common query patterns
CREATE INDEX idx_audit_logs_user_module ON audit_logs(user_id, module);
CREATE INDEX idx_audit_logs_path_method ON audit_logs(path, method);

-- +goose Down
DROP TABLE IF EXISTS audit_logs;
DROP TYPE IF EXISTS http_method_enum;