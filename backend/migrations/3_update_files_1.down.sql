DROP TABLE user_files;

ALTER TABLE files DROP COLUMN mime_type;

ALTER TABLE files DROP COLUMN ref_count;

ALTER TABLE users
DROP COLUMN actual_storage;

ALTER TABLE users
DROP COLUMN expected_storage;

ALTER TABLE files
ADD COLUMN user_id INT NOT NULL REFERENCES users(id);

ALTER TABLE files
ADD COLUMN uploaded_At TIMESTAMP;