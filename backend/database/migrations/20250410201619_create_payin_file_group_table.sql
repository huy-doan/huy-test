-- +goose Up
-- +goose StatementBegin
CREATE TABLE `payin_file_group` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主キー',
  `payment_provider_id` int NOT NULL COMMENT '決済会社ID',
  `import_target_date` date NOT NULL COMMENT '取り込み予定日',
  `imported_at` datetime DEFAULT NULL COMMENT '実際に取り込んだ日時',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日',
  PRIMARY KEY (`id`),
  KEY `idx_payment_provider_id` (`payment_provider_id`),
  KEY `idx_import_target_date` (`import_target_date`),
  KEY `idx_imported_at` (`imported_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='インポートされたCSVファイルのグルーピング';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `payin_file_group`;
-- +goose StatementEnd
