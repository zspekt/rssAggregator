-- +goose Up
-- FOR INTERNAL USE ONLY. don't run this migration
-- resolve duplicates before adding unique flag
CREATE TEMPORARY TABLE feeds_temp AS SELECT DISTINCT ON (url)
    *
FROM
    feeds;

-- Delete all records from the original table
DELETE FROM feeds;

-- Insert unique records back into the original table
INSERT INTO feeds
SELECT
    *
FROM
    feeds_temp;

-- Drop the temporary table
DROP TABLE feeds_temp;

-- Add unique constraint
ALTER TABLE feeds
    ADD UNIQUE (url);

-- +goose Down
-- Remove unique constraint
ALTER TABLE feeds
    DROP CONSTRAINT feeds_url_key;
