ALTER TABLE posts modify permalink VARCHAR(256) COLLATE utf8mb4_unicode_ci NOT NULL;
ALTER TABLE posts modify LONGTEXT COLLATE utf8mb4_unicode_ci NOT NULL;
ALTER TABLE posts ADD UNIQUE permalink;