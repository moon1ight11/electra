-- +goose Up
INSERT INTO
    workers (name, phone, role, password_hash)
VALUES
    ('Алексей Владелец', '79001234567', 'owner', 'owner123');

-- +goose Down
DELETE FROM workers
WHERE
    phone = '79001234567'
    AND role = 'owner';