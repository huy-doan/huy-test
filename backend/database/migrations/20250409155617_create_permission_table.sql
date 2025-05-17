-- +goose Up
-- +goose StatementBegin
CREATE TABLE `permission` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '権限ID（主キー）',
  `name` varchar(45) NOT NULL COMMENT '権限名',
  `code` varchar(45) NOT NULL COMMENT 'Role permission, for example: READ_ONLY, APPROVAL, ADMIN',
  `screen_id` int NOT NULL COMMENT '画面ID（外部キー）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `code_UNIQUE` (`code`),
  KEY `fk_permission_screen_idx` (`screen_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='権限定義';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `permission`;
-- +goose StatementEnd
