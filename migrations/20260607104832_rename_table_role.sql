-- +goose Up
RENAME TABLE role TO roles;

-- +goose Down
RENAME TABLE roles TO role;