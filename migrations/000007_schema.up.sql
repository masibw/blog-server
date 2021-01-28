ALTER TABLE posts modify permalink VARCHAR(256) COLLATE utf8mb4_unicode_ci;
ALTER TABLE posts modify content LONGTEXT COLLATE utf8mb4_unicode_ci;
ALTER TABLE posts DROP INDEX permalink;