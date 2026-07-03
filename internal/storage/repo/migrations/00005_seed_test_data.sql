-- +goose Up
-- Создаём владельца (пароль: owner123)
INSERT INTO
    workers (name, phone, role, password_hash)
VALUES
    ('Алексей Владелец', '79001234567', 'owner', 'owner123');

-- Создаём исполнителя (пароль: worker123)
INSERT INTO
    workers (name, phone, role, password_hash)
VALUES
    ('Петя Электрик', '79007654321', 'worker', 'worker123');

-- Пара заявок с сайта
INSERT INTO
    requests (
        name,
        phone,
        comment,
        status
    )
VALUES
    ('Иван Клиент', '79161234567', 'Нужно подключить дом на участке', 'new');

INSERT INTO
    requests (
        name,
        phone,
        comment,
        status
    )
VALUES
    ('Мария Клиент', '79262345678', 'Замена проводки в квартире', 'new');

-- +goose Down
DELETE FROM order_workers;

DELETE FROM orders;

DELETE FROM requests;

DELETE FROM workers;