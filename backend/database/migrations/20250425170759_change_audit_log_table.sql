-- +goose Up
-- +goose StatementBegin
ALTER TABLE `audit_log` 
    DROP KEY `idx_audit_type`,
    CHANGE COLUMN `audit_type_id` `audit_log_type` VARCHAR(255) NOT NULL,
    CHANGE COLUMN `outcoming_id` `payout_id` INT DEFAULT NULL,
    ADD COLUMN `payin_id` INT DEFAULT NULL AFTER `payout_id`,
    ADD KEY `idx_audit_type` (`audit_log_type`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `audit_log` 
    DROP KEY `idx_audit_type`,
    DROP COLUMN `payin_id`,
    CHANGE COLUMN `payout_id` `outcoming_id` INT DEFAULT NULL,
    CHANGE COLUMN `audit_log_type` `audit_type_id` TINYINT NOT NULL,
    ADD KEY `idx_audit_type` (`audit_type_id`);
-- +goose StatementEnd
