CREATE TABLE users
(
    --id SERIAL PRIMARY KEY,
    user_uuid UUID NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE sessions
(
    id UUID NOT NULL UNIQUE,
    user_uuid UUID NOT NULL,
    --user_id INT,
    refresh_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP
    --FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- TODO: unique constraint refresh token ???