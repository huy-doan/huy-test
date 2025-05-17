-- +goose Up
INSERT INTO `permission` 
    VALUES 
    (1,'ユーザー管理','USER_MANAGE',1,'2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (2,'ユーザーのロール変更','USER_ROLE_CHANGE',1,'2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (4,'システム全体のログ閲覧','SYSTEM_LOG_VIEW',2,'2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (5,'自分の個人データ変更','EDIT_OWN_PROFILE',3,'2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (6,'自分の行動ログ確認','VIEW_OWN_LOG',3,'2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (7,'管理画面の参照権限','VIEW_ADMIN_PANEL',4,'2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (8,'振込み承認（事業）','TRANSFER_APPROVE_BUSINESS',5,'2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (9,'振込み承認（経理）','TRANSFER_APPROVE_ACCOUNTANT',5,'2025-03-30 17:51:49','2025-03-30 17:51:49',NULL),
    (10,'手動振込機能','MANUAL_TRANSFER',6,'2025-03-30 17:51:49','2025-03-30 17:51:49',NULL)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    code = VALUES(code),
    screen_id = VALUES(screen_id),
    deleted_at = VALUES(deleted_at);

-- +goose Down
