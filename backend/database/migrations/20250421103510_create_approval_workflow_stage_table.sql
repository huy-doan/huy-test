-- +goose Up
-- +goose StatementBegin
CREATE TABLE `approval_workflow_stage` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主キー',
  `workflow_id` int NOT NULL COMMENT '承認ワークフローID（FK）',
  `stage_name` varchar(45) DEFAULT NULL COMMENT '承認ワークフローの段階名',
  `level` int NOT NULL COMMENT '承認ワークフローのレベル',
  `approver_role_id` int NOT NULL COMMENT '承認者のロールID',
  `approver_count` tinyint(1) DEFAULT '1' COMMENT '承認必要な承認者数',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  KEY `idx_workflow_id` (`workflow_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='承認ステップ定義マスター';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `approval_workflow_stage`;
-- +goose StatementEnd
