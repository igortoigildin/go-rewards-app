CREATE TABLE IF NOT EXISTS orders (
    number VARCHAR(15) PRIMARY KEY,
    status VARCHAR(10),
    user_id BIGINT,
    accrual INT,
    uploaded_at timestamp
);