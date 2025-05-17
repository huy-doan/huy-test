-- +goose Up
-- +goose StatementBegin
CREATE TABLE `approval_stage` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主キー',
  `approval_id` int NOT NULL COMMENT '承認ID',
  `approval_workflow_stage_id` int NOT NULL COMMENT '承認ワークフロー段階ID',
  `approver_id` int NOT NULL COMMENT '承認者ID',
  `approval_result` int NOT NULL COMMENT '承認結果 \n1: 承認 (APPROVED), 2: 却下(REJECTED)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  KEY `user_id` (`approver_id`),
  KEY `approval_stage_histories_ibfk_1` (`approval_id`),
  KEY `approval_stages_master_approval_workflow_stages_FK` (`approval_workflow_stage_id`),
  CONSTRAINT `approval_stage_ibfk_1` FOREIGN KEY (`approval_id`) REFERENCES `approval` (`id`) ON DELETE CASCADE,
  CONSTRAINT `approval_stages_master_approval_workflow_stages_FK` FOREIGN KEY (`approval_workflow_stage_id`) REFERENCES `approval_workflow_stage` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='承認ステップ詳細';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `approval_stage`;
-- +goose StatementEnd
