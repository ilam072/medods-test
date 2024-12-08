CREATE TABLE users
(
    user_uuid UUID NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE sessions
(
    id UUID NOT NULL UNIQUE,
    user_uuid UUID NOT NULL,
    refresh_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP,
    used BOOLEAN
);