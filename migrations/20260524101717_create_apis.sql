-- +goose Up
CREATE TABLE api (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    path VARCHAR(100) NOT NULL,
    method ENUM('GET', 'POST', 'PUT', 'DELETE', 'PATCH') NOT NULL,
    summary VARCHAR(500) NOT NULL,
    tags VARCHAR(100) NOT NULL,
    
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_api_path (path),
    INDEX idx_api_method (method),
    INDEX idx_api_summary (summary),
    INDEX idx_api_tags (tags)
);

-- +goose Down
DROP TABLE IF EXISTS api;