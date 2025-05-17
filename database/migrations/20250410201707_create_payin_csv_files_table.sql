-- +goose Up
-- +goose StatementBegin
CREATE TABLE `payin_csv_files` (
  `id` int NOT NULL AUTO_INCREMENT,
  `payment_provider_id` int NOT NULL,
  `payin_csv_group_id` int DEFAULT NULL,
  `file_name` varchar(255) NOT NULL,
  `file_content_key` varchar(255) NOT NULL,
  `payin_file_type_id` int DEFAULT NULL COMMENT 'CSVファイル種類を特定するためのカラム。Paypay連携時に複数なファイルタイプが存在しているため',
  `has_data_record` tinyint(1) NOT NULL,
  `added_manually` tinyint(1) NOT NULL,
  `content_added_manually` text NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_file_name` (`file_name`),
  KEY `fk_payment_provider_id` (`payment_provider_id`),
  KEY `fk_payin_csv_group_id` (`payin_csv_group_id`),
  KEY `fk_payin_file_type_id` (`payin_file_type_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='インポートされたCSVファイルの詳細';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `payin_csv_files`;
-- +goose StatementEnd
