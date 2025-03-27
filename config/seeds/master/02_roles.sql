-- +goose Up
INSERT INTO `roles` (id, name, code, created_at, updated_at, deleted_at)
VALUES
    (1, 'システム管理者', 'SYSTEM_ADMIN', now(), now(), NULL),
    (2, '一般ユーザー', 'GENERAL_USER', now(), now(), NULL),
    (3, '事業担当者', 'BUSINESS_USER', now(), now(), NULL),
    (4, '経理担当者', 'ACCOUNTING_USER', now(), now(), NULL)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    code = VALUES(code),
   deleted_at = VALUES(deleted_at);

-- +goose Down
