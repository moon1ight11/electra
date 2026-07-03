-- +goose Up
CREATE TABLE
    order_workers (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        order_id UUID NOT NULL REFERENCES orders (id),
        worker_id UUID NOT NULL REFERENCES workers (id),
        time_spent INT,
        earned_amount DECIMAL(12, 2),
        materials_used TEXT,
        notes TEXT,
        UNIQUE (order_id, worker_id)
    );

CREATE INDEX idx_order_workers_order ON order_workers (order_id);

CREATE INDEX idx_order_workers_worker ON order_workers (worker_id);

-- +goose Down
DROP TABLE IF EXISTS order_workers;