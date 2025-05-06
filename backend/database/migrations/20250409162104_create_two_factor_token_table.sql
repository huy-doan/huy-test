-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS two_factor_token (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    mfa_type INT NOT NULL DEFAULT 1,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    expired_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    KEY idx_user_id (user_id),
    INDEX idx_user_token (user_id, mfa_type, token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE two_factor_token;
-- +goose StatementEnd
