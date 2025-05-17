-- +goose Up
-- +goose StatementBegin
CREATE TABLE `paypay_payin_detail` (
  `id` int NOT NULL AUTO_INCREMENT,
  `payin_file_id` int NOT NULL,
  `payment_merchant_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '加盟店ID: PayPay側で管理されている加盟店ID',
  `merchant_business_name` varchar(100) DEFAULT NULL COMMENT '屋号: PayPay側に登録されている屋号',
  `cutoff_date` date DEFAULT NULL COMMENT '締め日: 入金サイクルに応じた締め日',
  `transaction_amount` decimal(18,2) DEFAULT NULL COMMENT '取引額: 入金処理対象期間内の取引額',
  `refund_amount` decimal(18,2) DEFAULT NULL COMMENT '返金額: 入金処理対象期間内の返金額',
  `usage_fee` decimal(18,2) DEFAULT NULL COMMENT '利用料: 取引および返金対象の利用料（負/正あり）',
  `platform_fee` decimal(18,2) DEFAULT NULL COMMENT 'プラットフォーム使用料: 月額サービス利用料',
  `initial_fee` decimal(18,2) DEFAULT NULL COMMENT '初期費用: 初期サービス利用料',
  `tax` decimal(18,2) DEFAULT NULL COMMENT '税: 上記に対する税額',
  `cashback` decimal(18,2) DEFAULT NULL COMMENT 'キャッシュバック: ユーザー向け控除額',
  `adjustment` decimal(18,2) DEFAULT NULL COMMENT '調整額: その他調整金額',
  `fee` decimal(18,2) DEFAULT NULL COMMENT '入金手数料: 入金処理にかかる手数料',
  `amount` decimal(18,2) DEFAULT NULL COMMENT '支払金額: 入金金額',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT '削除日時（論理削除用）',
  PRIMARY KEY (`id`),
  KEY `idx_payin_file_id` (`payin_file_id`),
  KEY `idx_payment_merchant_id` (`payment_merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='入金レポート（加盟店パート）';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `paypay_payin_detail`;
-- +goose StatementEnd
