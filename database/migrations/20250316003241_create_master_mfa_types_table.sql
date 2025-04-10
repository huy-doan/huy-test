-- +goose Up
-- +goose StatementBegin
CREATE TABLE `master_mfa_types` (
  `id` int NOT NULL AUTO_INCREMENT,
  `no` int NOT NULL,
  `title` varchar(255) NOT NULL,
  `is_active` int NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE master_mfa_types;
-- +goose StatementEnd
