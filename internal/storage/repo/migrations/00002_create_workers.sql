-- +goose Up
CREATE TYPE worker_role AS ENUM('owner', 'worker');

CREATE TABLE
    workers (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name VARCHAR(255) NOT NULL,
        phone VARCHAR(50),
        role worker_role NOT NULL DEFAULT 'worker',
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- +goose Down
DROP TABLE IF EXISTS workers;

DROP TYPE IF EXISTS worker_role;