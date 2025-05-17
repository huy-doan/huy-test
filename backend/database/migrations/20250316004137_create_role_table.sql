-- +goose Up
-- +goose StatementBegin
CREATE TABLE `role` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_name` (`name`),
  KEY `deleted_at` (`deleted_at`)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE role;
-- +goose StatementEnd
