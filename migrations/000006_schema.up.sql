CREATE TABLE IF NOT EXISTS `posts_tags` (
  `id` CHAR(26) NOT NULL,
  `post_id` CHAR(26) COLLATE utf8mb4_unicode_ci NOT NULL,
  `tag_id` CHAR(26) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE(`post_id`,`tag_id`),
  FOREIGN KEY(`post_id`) REFERENCES  posts(id),
  FOREIGN KEY(`tag_id`) REFERENCES  tags(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
