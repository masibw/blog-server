ALTER TABLE tags ADD post_id CHAR(26) COLLATE utf8mb4_unicode_ci NOT NULL;
ALTER TABLE tags ADD CONSTRAINT tags_ibfk_1 FOREIGN KEY (post_id) REFERENCES posts(id);
