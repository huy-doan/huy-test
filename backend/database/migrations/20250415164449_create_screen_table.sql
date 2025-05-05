-- +goose Up
-- +goose StatementBegin
CREATE TABLE `screen` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '画面ID（主キー）',
  `name` varchar(45) NOT NULL COMMENT '画面名',
  `screen_code` varchar(45) NOT NULL COMMENT '画面コード, for example: USER_MANAGEMENT_SCREEN, SYSTEM_LOG_SCREEN, PROFILE_SCREEN',
  `screen_path` varchar(45) NOT NULL COMMENT '画面パス /add/*, /edit/*',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `screen_code_UNIQUE` (`screen_code`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='画面定義';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `screen`;
-- +goose StatementEnd
