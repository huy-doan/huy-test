-- +goose Up
-- +goose StatementBegin
CREATE TABLE `payment_provider` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主キー',
  `code` varchar(255) NOT NULL COMMENT '決済会社コード（例：PAYPAY）',
  `name` varchar(255) DEFAULT NULL COMMENT '決済会社名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='決済プロバイダーマスター';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `payment_provider`;
-- +goose StatementEnd
