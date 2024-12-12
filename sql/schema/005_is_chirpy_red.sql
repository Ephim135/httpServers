-- +goose Up
Alter table users
Add column is_chirpy_red boolean not null default false;

-- +goose Down
Alter table users drop column is_chirpy_red;