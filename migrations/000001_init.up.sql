CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    login      VARCHAR(20) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,  -- bcrypt hash
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE files (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name   VARCHAR(255) NOT NULL,
    path        TEXT NOT NULL UNIQUE,   -- data/{login}/{filename}
    size_bytes  BIGINT,
    mime_type   VARCHAR(100),
    gd          INTEGER DEFAULT 0,      -- счётчик Get/Delete операций
    uploaded_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE sessions (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
