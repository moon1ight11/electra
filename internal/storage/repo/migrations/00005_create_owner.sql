-- +goose Up
INSERT INTO
    workers (name, phone, role, password_hash)
VALUES
    (
        'Дмитрий Мягков',
        '89009231984',
        'owner',
        '$2a$10$C0x/Tj/BlsVwQ.FICakdzezUZT0RLILj9M3QKGaf3E7iXDH7GgFdi'
    );

-- +goose Down
DELETE FROM workers
WHERE
    phone = '79001234567'
    AND role = 'owner';