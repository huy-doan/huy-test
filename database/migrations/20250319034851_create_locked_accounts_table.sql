-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `locked_accounts` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `count` int NOT NULL,
  `locked_at` datetime DEFAULT NULL,
  `expired_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_fraud_lock_email` (`email`),
  KEY `fk_locked_account_user_idx` (`user_id`),
  CONSTRAINT `fk_locked_account_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table locked_accounts;
-- +goose StatementEnd
