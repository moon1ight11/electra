-- +goose Up
CREATE TABLE
    orders (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        request_id UUID REFERENCES requests (id),
        address TEXT NOT NULL,
        description TEXT,
        estimated_price DECIMAL(12, 2),
        planned_date DATE,
        created_by UUID NOT NULL REFERENCES workers (id),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        completed_at TIMESTAMPTZ
    );

CREATE INDEX idx_orders_request_id ON orders (request_id);

CREATE INDEX idx_orders_created_by ON orders (created_by);

CREATE INDEX idx_orders_completed_at ON orders (completed_at);

-- +goose Down
DROP TABLE IF EXISTS orders;