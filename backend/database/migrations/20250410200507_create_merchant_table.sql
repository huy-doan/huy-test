-- +goose Up
-- +goose StatementBegin
CREATE TABLE `merchant` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主キー',
  `payment_provider_id` int NOT NULL COMMENT '決済会社ID',
  `payment_merchant_id` varchar(255) DEFAULT NULL COMMENT '決済会社が発行した加盟店ID',
  `merchant_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '加盟店名',
  `shop_id` int NOT NULL COMMENT 'メイクショップのショップID',
  `shop_url` varchar(255) DEFAULT NULL COMMENT 'メイクショップのURL',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='加盟店情報';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `merchant`;
-- +goose StatementEnd
