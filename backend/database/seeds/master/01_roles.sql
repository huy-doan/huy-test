-- +goose Up
INSERT INTO `role` 
    VALUES 
    (1,'システム管理者','SYSTEM_ADMIN','2025-03-19 00:13:08','2025-03-19 00:13:08',NULL),
    (2,'一般ユーザー','GENERAL_USER','2025-03-19 00:13:08','2025-03-19 00:13:08',NULL),
    (3,'事業担当者','BUSINESS_USER','2025-03-19 00:13:08','2025-03-19 00:13:08',NULL),
    (4,'経理担当者','ACCOUNTING_USER','2025-03-19 00:13:08','2025-03-19 00:13:08',NULL)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    code = VALUES(code),
   deleted_at = VALUES(deleted_at);

-- +goose Down
