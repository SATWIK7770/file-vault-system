CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    size BIGINT NOT NULL,
    hash TEXT NOT NULL,
    storage_path TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT now()
);