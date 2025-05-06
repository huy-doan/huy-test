-- +goose Up
-- +goose StatementBegin
CREATE TABLE `payment_provider_review` (
  `id` int NOT NULL AUTO_INCREMENT,
  `payment_provider_id` int NOT NULL COMMENT '決済サービスプロバイダーid',
  `merchant_review_status` int NOT NULL COMMENT '''1:審査中, 2:審査通過, 3:審査否認''',
  `payment_merchant_id` int DEFAULT NULL,
  `merchant_id` int DEFAULT NULL,
  `registered_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '申請日時',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT '削除日時（論理削除用）',
  PRIMARY KEY (`id`),
  KEY `shop_payment_reviews_merchants_FK` (`merchant_id`),
  CONSTRAINT `shop_payment_reviews_merchants_FK` FOREIGN KEY (`merchant_id`) REFERENCES `merchant` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='加盟店ごとの決済プロバイダー利用状況テーブル';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `payment_provider_review`;
-- +goose StatementEnd
