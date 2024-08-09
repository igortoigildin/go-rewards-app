CREATE TABLE IF NOT EXISTS users (
    user_id bigserial PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash bytea NOT NULL
);