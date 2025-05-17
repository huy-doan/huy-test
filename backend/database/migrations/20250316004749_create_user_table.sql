-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `role_id` int NOT NULL,
  `enabled_mfa` tinyint(1) NOT NULL DEFAULT '1',
  `mfa_type` int DEFAULT '1' COMMENT '2段階認証種類\n1:メール',
  `full_name` varchar(200) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `deleted_at` (`deleted_at`),
  KEY `fk_user_role` (`role_id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd
