-- +goose Up
-- +goose StatementBegin
CREATE TABLE `paypay_payin_summaries` (
  `id` int NOT NULL AUTO_INCREMENT,
  `payin_csv_file_id` int NOT NULL,
  `payin_csv_group_id` int DEFAULT NULL,
  `corporate_name` varchar(255) DEFAULT NULL COMMENT '法人名: 本レポート生成の対象となる法人名',
  `cutoff_date` date DEFAULT NULL COMMENT '締め日: 入金対象取引の締め日',
  `payment_date` date DEFAULT NULL COMMENT '支払日: 入金日（入金予定の場合も含む）',
  `transaction_amount` decimal(18,2) DEFAULT NULL COMMENT '取引額: 取引総額（加盟店パートの合計値）',
  `refund_amount` decimal(18,2) DEFAULT NULL COMMENT '返金額: 返金総額（加盟店パートの合計値）',
  `usage_fee` decimal(18,2) DEFAULT NULL COMMENT '利用料: 取引と返金に基づく利用料（合計値）',
  `platform_fee` decimal(18,2) DEFAULT NULL COMMENT 'プラットフォーム使用料: 月額サービス利用料（合計値）',
  `initial_fee` decimal(18,2) DEFAULT NULL COMMENT '初期費用: 初期サービス利用料（合計値）',
  `tax` decimal(18,2) DEFAULT NULL COMMENT '税: 上記に対する税額（合計値）',
  `cashback` decimal(18,2) DEFAULT NULL COMMENT 'キャッシュバック: ユーザー向け控除額（合計値）',
  `adjustment` decimal(18,2) DEFAULT NULL COMMENT '調整額: その他調整金額（合計値）',
  `fee` decimal(18,2) DEFAULT NULL COMMENT '入金手数料: 入金にかかる手数料（合計値）',
  `amount` decimal(18,2) DEFAULT NULL COMMENT '支払金額: 入金額（合計値）',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT '削除日時（論理削除用）',
  PRIMARY KEY (`id`),
  KEY `fk_payin_csv_file_id` (`payin_csv_file_id`),
  KEY `fk_payin_csv_group_id` (`payin_csv_group_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='入金レポート（法人パート）';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `paypay_payin_summaries`;
-- +goose StatementEnd
