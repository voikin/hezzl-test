-- +goose Up
-- +goose StatementBegin

CREATE TABLE
    IF NOT EXISTS projects (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL);

INSERT INTO
    projects (name)
VALUES
    ('Первая запись');

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS projects;

-- +goose StatementEnd