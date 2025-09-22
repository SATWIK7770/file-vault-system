ALTER TABLE user_files ADD COLUMN is_owner BOOLEAN DEFAULT FALSE;

ALTER TABLE user_files ADD COLUMN visibility TEXT DEFAULT 'private'; 
ALTER TABLE user_files ADD COLUMN public_token TEXT UNIQUE;        
