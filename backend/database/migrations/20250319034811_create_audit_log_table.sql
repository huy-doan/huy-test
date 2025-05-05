-- +goose Up
-- +goose StatementBegin
CREATE TABLE `audit_log` (
    `id` int NOT NULL AUTO_INCREMENT,
    `user_id` int DEFAULT NULL,
    `audit_type_id` tinyint NOT NULL,
    `description` varchar(512) DEFAULT NULL,
    `transaction_id` int DEFAULT NULL,
    `outcoming_id` int DEFAULT NULL,
    `user_agent` varchar(255) DEFAULT NULL,
    `ip_address` varchar(50) DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`,`created_at`),
    KEY `idx_audit_type` (`audit_type_id`),
    KEY `fk_audit_log_user_idx` (`user_id`)
)
    PARTITION BY RANGE (to_days(`created_at`))
(PARTITION p202501 VALUES LESS THAN (739648) ENGINE = InnoDB,
 PARTITION p202502 VALUES LESS THAN (739676) ENGINE = InnoDB,
 PARTITION p202503 VALUES LESS THAN (739707) ENGINE = InnoDB,
 PARTITION p202504 VALUES LESS THAN (739737) ENGINE = InnoDB,
 PARTITION p202505 VALUES LESS THAN (739768) ENGINE = InnoDB,
 PARTITION p202506 VALUES LESS THAN (739798) ENGINE = InnoDB,
 PARTITION p202507 VALUES LESS THAN (739829) ENGINE = InnoDB,
 PARTITION p202508 VALUES LESS THAN (739860) ENGINE = InnoDB,
 PARTITION p202509 VALUES LESS THAN (739890) ENGINE = InnoDB,
 PARTITION p202510 VALUES LESS THAN (739921) ENGINE = InnoDB,
 PARTITION p202511 VALUES LESS THAN (739951) ENGINE = InnoDB,
 PARTITION p202512 VALUES LESS THAN (739982) ENGINE = InnoDB,
 PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB)


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table audit_log;
-- +goose StatementEnd
