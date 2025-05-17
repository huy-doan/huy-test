-- +goose Up
INSERT INTO `payment_provider` (id, code, name, created_at, updated_at) VALUES
(1, 'PAYPAY', 'PayPay', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON DUPLICATE KEY UPDATE
    code = VALUES(code),
    name = VALUES(name);
-- +goose Down
TRUNCATE TABLE `payment_provider`;
