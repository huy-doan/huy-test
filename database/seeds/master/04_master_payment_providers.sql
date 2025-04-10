-- +goose Up
INSERT INTO `master_payment_providers` (id, code, name, created_at, updated_at) VALUES
(1, 'PAYPAY', 'PayPay', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON DUPLICATE KEY UPDATE
    code = VALUES(code),
    name = VALUES(name),
    updated_at = IF(
        code != VALUES(code) OR name != VALUES(name),
        CURRENT_TIMESTAMP,
        updated_at
    );
-- +goose Down
TRUNCATE TABLE `master_payment_providers`;
