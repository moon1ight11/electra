-- +goose Up
CREATE TYPE request_status AS ENUM ('new', 'converted', 'cancelled');

CREATE TABLE
    requests (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name VARCHAR(255) NOT NULL,
        phone VARCHAR(50) NOT NULL,
        comment TEXT,
        status request_status NOT NULL DEFAULT 'new',
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- +goose Down
DROP TABLE IF EXISTS requests;
DROP TYPE IF EXISTS request_status;