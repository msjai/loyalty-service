CREATE TABLE IF NOT EXISTS users(
    id      SERIAL PRIMARY KEY,
    login    VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    balance NUMERIC DEFAULT 0
);

CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY,
    number      VARCHAR NOT NULL UNIQUE,
    status      VARCHAR NOT NULL,
    user_id     INTEGER NOT NULL
        CONSTRAINT orders_users_id_fk
            REFERENCES users  ON UPDATE CASCADE ON delete CASCADE,

    accrual_sum NUMERIC DEFAULT 0,
    uploaded_at           DATE    NOT NULL
);

CREATE TABLE IF NOT EXISTS writes_off
(
    id             SERIAL PRIMARY KEY,
    order_woff_num VARCHAR NOT NULL,
    sum            NUMERIC NOT NULL,
    user_id        INTEGER NOT NULL
        CONSTRAINT writes_off_users_id_fk
            REFERENCES users
            ON  UPDATE CASCADE ON DELETE CASCADE ,
    date           DATE    NOT NULL
);