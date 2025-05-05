-- +goose Up
-- +goose StatementBegin
CREATE TABLE `role_permission` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ロール・権限の連携テーブル',
  `role_id` int NOT NULL COMMENT 'ロールID',
  `permission_id` int NOT NULL COMMENT '権限ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_id` (`role_id`,`permission_id`),
  KEY `permission_id` (`permission_id`),
  CONSTRAINT `fk_role_perms_permissions` FOREIGN KEY (`permission_id`) REFERENCES `permission` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_role_perms_roles` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ロールと権限の紐付け';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `role_permission`;
-- +goose StatementEnd
