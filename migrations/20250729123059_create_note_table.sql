-- +goose Up
CREATE TABLE note (
    id serial PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    created_at timestamp NOT NULL default now(),
    updated_at timestamp
);

-- +goose Down
DROP TABLE note;