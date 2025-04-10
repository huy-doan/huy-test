-- +goose Up
INSERT INTO `master_mfa_types` (id, no, title, is_active, created_at, updated_at, deleted_at)
VALUES
   (1, 1, 'OTP', 1, NOW(), NOW(), NULL),
   (2, 2, 'メール', 1, NOW(), NOW(), NULL),
   (3, 3, 'SMS', 1, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE
   no = VALUES(no),
   title = VALUES(title),
   is_active = VALUES(is_active),
   deleted_at = VALUES(deleted_at);

-- +goose Down
