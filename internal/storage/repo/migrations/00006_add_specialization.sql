-- +goose Up
ALTER TABLE workers
ADD COLUMN specialization VARCHAR(100) DEFAULT '';

-- +goose Down
ALTER TABLE workers
DROP COLUMN IF EXISTS specialization;