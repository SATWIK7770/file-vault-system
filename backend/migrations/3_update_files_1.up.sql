ALTER TABLE files 
ADD COLUMN mime_type VARCHAR(255) NOT NULL;

ALTER TABLE files 
ADD COLUMN ref_count INT DEFAULT 1;

ALTER TABLE users
ADD COLUMN actual_storage BIGINT NOT NULL DEFAULT 0;

ALTER TABLE users
ADD COLUMN expected_storage BIGINT NOT NULL DEFAULT 0;

ALTER TABLE files
DROP COLUMN user_id;

ALTER TABLE files
DROP COLUMN uploaded_At;

CREATE TABLE user_files (
    id             SERIAL PRIMARY KEY,
    user_id        INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_id        INT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    file_name      VARCHAR(255),
    uploaded_at    TIMESTAMP DEFAULT NOW(),
    download_times INT DEFAULT 0,
    UNIQUE(user_id, file_id)
);

