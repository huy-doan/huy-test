-- +goose Up
-- +goose StatementBegin
CREATE TABLE `payout_record` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主キー',
  `shop_id` int NOT NULL COMMENT 'メイクショップの加盟店ID',
  `payout_id` int NOT NULL COMMENT '振込ID',
  `transaction_id` int NOT NULL COMMENT '振込データID',
  `bank_name` varchar(255) NOT NULL COMMENT '加盟店の金融機関名義',
  `bank_code` varchar(255) NOT NULL COMMENT '加盟店の金融機関コード',
  `branch_name` varchar(255) NOT NULL COMMENT '加盟店の金融機関支店名',
  `branch_code` varchar(255) NOT NULL COMMENT '加盟店の金融機関支店コード',
  `bank_account_type` int NOT NULL COMMENT '口座種別　\n1:普通預金, 2:当座預金, 3:定期預金',
  `account_no` varchar(255) NOT NULL COMMENT '加盟店の金融機関の口座番号',
  `account_name` varchar(255) NOT NULL COMMENT '加盟店の金融機関の口座名義',
  `amount` decimal(18,2) NOT NULL COMMENT '振込金額',
  `transfer_status` int NOT NULL COMMENT '振込状態 \n1:振込中, 2:ホワイトリスト追加エラー, 3:振込依頼APIエラー, 4:振込依頼失敗, 5:振込依頼済み, 6:送金手続き済み',
  `sending_date` date DEFAULT NULL COMMENT '振込予定日',
  `aozora_transfer_apply_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'あおぞらネット銀行の振込実施番号',
  `transfer_requested_at` datetime DEFAULT NULL COMMENT '振込実施依頼日時',
  `transfer_executed_at` datetime DEFAULT NULL COMMENT '振込実施日時',
  `transfer_request_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '振込実施エラー内容',
  `idempotency_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '冪等キー(IdempotencyKey)\n二重送金防止するための冪等キー',
  `created_at` datetime DEFAULT NULL COMMENT 'レコード作成日時',
  `updated_at` datetime DEFAULT NULL COMMENT 'レコード更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日時',
  PRIMARY KEY (`id`),
  KEY `idx_shop_id` (`shop_id`),
  KEY `idx_outcoming_id` (`payout_id`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_ganb_transfer_apply_no` (`aozora_transfer_apply_no`),
  KEY `idx_ganb_idempotency_key` (`idempotency_key`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_outcoming_outcoming_record` FOREIGN KEY (`payout_id`) REFERENCES `payout` (`id`),
  CONSTRAINT `fk_outcoming_records_shop` FOREIGN KEY (`shop_id`) REFERENCES `merchant` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='出金取引詳細';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `payout_record`;
-- +goose StatementEnd
