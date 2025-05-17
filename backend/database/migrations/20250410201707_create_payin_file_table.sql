-- +goose Up
-- +goose StatementBegin
CREATE TABLE `payin_file` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主キー',
  `payment_provider_id` int NOT NULL COMMENT '決済会社ID',
  `payin_file_group_id` int DEFAULT NULL COMMENT '決済会社の取引・入金明細ファイルID',
  `file_name` varchar(255) NOT NULL COMMENT '決済会社の取引・入金明細ファイル名',
  `file_content_key` varchar(255) NOT NULL COMMENT '決済会社の取引・入金明細ファイルキー',
  `payin_file_type` int DEFAULT NULL COMMENT '決済会社の取引・入金明細ファイル種類　\n''1:入金レポート, 2:入金明細''\\nCSVファイル種類を特定するためのカラム。Paypay連携時に複数なファイルタイプが存在しているため',
  `has_data_record` tinyint(1) NOT NULL COMMENT '決済会社の取引・入金明細ファイルはデータ有無フラグ',
  `added_manually` tinyint(1) NOT NULL COMMENT '決済会社の取引・入金明細ファイルが手動アップロードであるか',
  `content_added_manually` text NOT NULL COMMENT '手動でアップロード時の決済会社の取引・入金明細ファイルの内容',
  `import_status` int DEFAULT NOT NULL COMMENT '決済会社の取引・入金明細ファイルのインポート状況\n''0:未インポート, 1:インポート中, 2:インポート完了, 3:インポート失敗''',
  `download_status` int DEFAULT NOT NULL COMMENT '決済会社の取引・入金明細ファイルのダウンロード状況\n''0:未ダウンロード, 1:ダウンロード中, 2:ダウンロード完了, 3:ダウンロード失敗''',
  `upload_status` int DEFAULT NOT NULL COMMENT '決済会社の取引・入金明細ファイルのアップロード状況\n''0:未アップロード, 1:アップロード中, 2:アップロード完了, 3:アップロード失敗''',
  `created_at` datetime DEFAULT NULL COMMENT 'レコード作成日',
  `updated_at` datetime DEFAULT NULL COMMENT 'レコード更新日',
  `deleted_at` datetime DEFAULT NULL COMMENT 'レコード削除日',
  PRIMARY KEY (`id`),
  KEY `idx_file_name` (`file_name`),
  KEY `idx_payment_provider_id` (`payment_provider_id`),
  KEY `idx_payin_file_group_id` (`payin_file_group_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='インポートされたCSVファイルの詳細';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `payin_file`;
-- +goose StatementEnd
