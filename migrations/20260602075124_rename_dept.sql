-- +goose Up
RENAME TABLE dept TO departments;

-- +goose Down
RENAME TABLE departments TO dept;