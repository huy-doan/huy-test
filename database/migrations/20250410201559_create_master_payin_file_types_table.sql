-- +goose Up
-- +goose StatementBegin
CREATE TABLE `master_payin_file_types` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(100) NOT NULL COMMENT 'ファイル種別コード（英字）',
  `title` varchar(255) NOT NULL COMMENT '日本語ラベル',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT '削除日時（論理削除用）',
  PRIMARY KEY (`id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='入金ファイル種別マスタテーブル';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `master_payin_file_types`;
-- +goose StatementEnd
