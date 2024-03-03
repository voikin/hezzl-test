-- +goose Up
-- +goose StatementBegin

CREATE TABLE
    IF NOT EXISTS goods (
        id SERIAL PRIMARY KEY,
        project_id INTEGER REFERENCES projects (id),
        name VARCHAR(255) NOT NUll,
        description VARCHAR(255) NOT NULL DEFAULT '',
        priority SERIAL,
        removed BOOLEAN DEFAULT false,
        created_at TIMESTAMP NOT NULL default now ()
    );

CREATE index ON goods USING btree (name);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS goods;

-- +goose StatementEnd