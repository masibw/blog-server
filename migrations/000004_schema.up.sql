ALTER TABLE posts ADD published_at DATETIME;
ALTER TABLE posts modify title VARCHAR(64) AFTER id;