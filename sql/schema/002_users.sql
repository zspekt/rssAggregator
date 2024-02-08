-- +goose Up
ALTER TABLE users
    ADD COLUMN apikey varchar(64) NOT NULL UNIQUE;

UPDATE
    users
SET
    apikey = encode(sha256 (random()::text::bytea), 'hex');

-- +goose Down
ALTER TABLE users
    DROP COLUMN apikey;
