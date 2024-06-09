-- mysql

CREATE TABLE `text` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'PK',
  `content` mediumtext NOT NULL COMMENT 'text content',
  `type` varchar(100) DEFAULT NULL COMMENT 'type, such markdown, golang, c++, java, python, html, javascript etc.',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'the lastest update time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='text';

CREATE TABLE `text_comment` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'PK',
  `text_id` int NOT NULL COMMENT 'PK of text',
  `comment` varchar(500) NOT NULL COMMENT 'Comment for text',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'the lastest update time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Text comment';

CREATE TABLE `text_tag` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `text_id` int NOT NULL COMMENT 'ID of a Text',
  `name` varchar(500) NOT NULL COMMENT 'Tag name',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'the latest update time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

