-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `master_audit_log_types` (
  `id` tinyint NOT NULL,
  `code` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table master_audit_log_types;
-- +goose StatementEnd
