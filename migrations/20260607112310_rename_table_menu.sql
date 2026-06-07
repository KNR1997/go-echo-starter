-- +goose Up
RENAME TABLE menu TO menus;

-- +goose Down
RENAME TABLE menus TO menu;