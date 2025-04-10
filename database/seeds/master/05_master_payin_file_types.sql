-- +goose Up
INSERT INTO `master_payin_file_types` (id, code, title, created_at, updated_at) VALUES
(1, 'PAYPAY_PAYOUT_INCOMI_CSV_FILE', '入金レポート', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'PAYPAY_PAYOUT_TRANSACTION_FILE', '入金明細', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON DUPLICATE KEY UPDATE
    code = VALUES(code),
    title = VALUES(title),
    updated_at = IF(
        code != VALUES(code) OR title != VALUES(title),
        CURRENT_TIMESTAMP,
        updated_at
    );
-- +goose Down
TRUNCATE TABLE `master_payin_file_types`;
