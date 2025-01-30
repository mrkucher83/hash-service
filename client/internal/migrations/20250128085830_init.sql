-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hash_storage (
    id         SERIAL PRIMARY KEY,
    hash_value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE hash_storage;
-- +goose StatementEnd
