-- +goose Up
CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL,
    stripe_id TEXT UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,
    jwt TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose Down
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS user.sessions;
