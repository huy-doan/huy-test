-- +goose Up
-- +goose StatementBegin
CREATE TABLE `payin_csv_groups` (
  `id` int NOT NULL AUTO_INCREMENT,
  `payment_provider_id` int NOT NULL,
  `import_target_date` date NOT NULL,
  `imported_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_payment_provider_id` (`payment_provider_id`),
  KEY `idx_import_target_date` (`import_target_date`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='インポートされたCSVファイルのグルーピング';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `payin_csv_groups`;
-- +goose StatementEnd
