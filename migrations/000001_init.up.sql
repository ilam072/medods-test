CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    user_uuid UUID NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE sessions
(
    id UUID,
    user_id INT,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);