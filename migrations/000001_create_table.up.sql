CREATE TABLE IF NOT EXISTS users (
    user_id bigserial PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash bytea NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    number VARCHAR(15) PRIMARY KEY,
    status VARCHAR(10),
    user_id BIGINT,
    accrual INT,
    uploaded_at timestamp
);

CREATE TABLE IF NOT EXISTS withdrawals (
    order VARCHAR(15) PRIMARY KEY,
    sum INT,
    user_id BIGINT,
    date timestamp
);