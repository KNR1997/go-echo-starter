-- +goose Up
RENAME TABLE api TO apis;

-- +goose Down
RENAME TABLE apis TO api;