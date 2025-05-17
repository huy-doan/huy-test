-- +goose Up
-- +goose StatementBegin
CREATE TABLE `transaction` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '振込集約データID(主キー)',
  `shop_id` int NOT NULL COMMENT 'メイクショップの加盟店ID',
  `transaction_status` int NOT NULL COMMENT '振込集約データの状態\n1:処理中, 2:承認待ち, 3:承認済み, 4:送金依頼済, 5:送金依頼失敗',
  `payout_id` int NOT NULL COMMENT '振込ID',
  `payout_record_id` int NOT NULL COMMENT '振込明細レコードID',
  `created_at` datetime DEFAULT NULL COMMENT 'レコード作成日時',
  `updated_at` datetime DEFAULT NULL COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  KEY `idx_shop_id` (`shop_id`),
  KEY `idx_outcoming_id` (`payout_id`),
  KEY `idx_outcoming_record_id` (`payout_record_id`),
  KEY `idx_shop_status` (`shop_id`,`transaction_status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='取引集約';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `transaction`;
-- +goose StatementEnd
