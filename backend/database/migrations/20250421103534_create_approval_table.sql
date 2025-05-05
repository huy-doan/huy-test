-- +goose Up
-- +goose StatementBegin
CREATE TABLE `approval` (
  `id` int NOT NULL AUTO_INCREMENT,
  `approval_workflow_id` int NOT NULL COMMENT '承認ワークフローID(FK)',
  `approval_status` int NOT NULL COMMENT '承認状況　\n''1'', ''承認待ち(PENDING)''\\n''2'', ''承認中(WAIT_APPROVAL)''\\n''3'', ''承認済み(APPROVED)''\\n''4'', ''却下(REJECTED)''',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  KEY `idx_workflow_id` (`approval_workflow_id`),
  KEY `idx_workflow_stage_status` (`approval_status`),
  CONSTRAINT `fk_approval_txn_workflow` FOREIGN KEY (`approval_workflow_id`) REFERENCES `approval_workflow` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='取引承認リクエスト';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `approval`;
-- +goose StatementEnd
