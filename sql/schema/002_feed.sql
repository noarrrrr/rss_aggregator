-- +goose Up
CREATE TABLE feeds (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null,
    url text not null unique,
    user_id uuid not null,
    constraint dc_id
      foreign key (user_id)
      references users(id) on delete cascade
);

-- +goose Down
DROP TABLE feeds;