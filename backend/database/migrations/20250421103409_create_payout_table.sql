-- +goose Up
-- +goose StatementBegin
CREATE TABLE `payout` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主キー',
  `payout_status` int NOT NULL COMMENT '振込状態　1:ドラフト, 2:振込データ作成済み, 3:送金手続き済み',
  `total` decimal(18,2) NOT NULL COMMENT '振込金額合計',
  `total_count` int NOT NULL COMMENT '振込の件数合計',
  `sending_date` date NOT NULL COMMENT '振込予定日',
  `sent_date` datetime NOT NULL COMMENT '実際振込の実施日時',
  `aozora_transfer_apply_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'あおぞらネット銀行の振込番号',
  `approval_id` int DEFAULT NULL COMMENT '振込承認ID',
  `user_id` int NOT NULL COMMENT '申請者ID',
  `created_at` datetime DEFAULT NULL COMMENT 'レコード作成日時',
  `updated_at` datetime DEFAULT NULL COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  KEY `idx_total` (`total`),
  KEY `idx_total_count` (`total_count`),
  KEY `idx_sent_date` (`sent_date`),
  KEY `idx_ganb_transfer_apply_no` (`aozora_transfer_apply_no`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='出金取引集約';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `payout`;
-- +goose StatementEnd
