-- +goose Up
INSERT INTO `screen` 
    (`id`, `name`, `screen_code`, `screen_path`, `created_at`, `updated_at`, `deleted_at`)
    VALUES 
    (1,'ユーザー管理画面','USER_MANAGEMENT_SCREEN','/user/*','2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (2,'システムログ画面','SYSTEM_LOG_SCREEN','/log/*','2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (3,'プロフィール画面','PROFILE_SCREEN','/profile/*','2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (4,'管理者パネル画面','ADMIN_PANEL_SCREEN','/admin/*','2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (5,'振込承認画面','TRANSFER_APPROVAL_SCREEN','/transfer/approval/*','2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (6,'振込操作画面','TRANSFER_OPERATION_SCREEN','/transfer/operation/*','2025-03-30 17:51:49','2025-03-30 17:51:49',NULL)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    screen_code = VALUES(screen_code),
    screen_path = VALUES(screen_path),
    deleted_at = VALUES(deleted_at);

-- +goose Down
TRUNCATE TABLE `screen`;
