-- +goose Up
-- +goose StatementBegin
CREATE TABLE `approval_workflow` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主キー',
  `name` varchar(255) NOT NULL COMMENT '承認ワークフロー名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='承認ワークフロー定義マスター';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `approval_workflow`;
-- +goose StatementEnd
