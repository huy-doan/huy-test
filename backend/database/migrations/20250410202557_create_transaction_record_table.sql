-- +goose Up
-- +goose StatementBegin
CREATE TABLE `transaction_record` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '振込明細データID（主キー）',
  `transaction_id` int NOT NULL COMMENT '振込集約データID',
  `merchant_id` int DEFAULT NULL COMMENT '決済会社が発行した加盟店 ID',
  `payin_detail_id` int NOT NULL COMMENT '決済会社から送付した取引明細ファイルのレコードID',
  `payin_summary_id` int DEFAULT NULL COMMENT '決済会社から送付した入金レポート（法人パート）レコードID',
  `transaction_record_type` int NOT NULL COMMENT '振込データ種別\n1:入金, 2:手数料, 3:振込手数料',
  `title` varchar(255) NOT NULL COMMENT '振込データのタイトル',
  `amount` decimal(18,2) NOT NULL COMMENT '振込データの金額（入金金額、振込手数料、MS手数料）',
  `created_at` datetime DEFAULT NULL COMMENT 'レコード作成日時',
  `updated_at` datetime DEFAULT NULL COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_merchant_id` (`merchant_id`),
  KEY `idx_master_transaction_record_type_id` (`payin_detail_id`),
  KEY `idx_payin_summary_id` (`payin_summary_id`),
  KEY `idx_amount` (`amount`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='集約取引詳細';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `transaction_record`;
-- +goose StatementEnd
