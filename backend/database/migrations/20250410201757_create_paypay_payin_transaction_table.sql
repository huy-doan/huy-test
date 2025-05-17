-- +goose Up
-- +goose StatementBegin
CREATE TABLE `paypay_payin_transaction` (
  `id` int NOT NULL AUTO_INCREMENT,
  `payin_file_id` int NOT NULL,
  `payment_transaction_id` int DEFAULT NULL COMMENT '決済番号: PayPay側で発行している決済番号',
  `payment_merchant_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '加盟店ID: PayPay側で管理されている加盟店ID',
  `merchant_business_name` varchar(100) DEFAULT NULL COMMENT '屋号: PayPay側に登録されている屋号',
  `shop_id` varchar(100) DEFAULT NULL COMMENT '店舗ID: 決済電文の店舗コード',
  `shop_name` varchar(255) DEFAULT NULL COMMENT '店舗名: MPM方式でのみ設定される項目',
  `terminal_code` varchar(100) DEFAULT NULL COMMENT '端末番号/PosID: 決済電文の端末コード',
  `payment_transaction_status` int DEFAULT NULL COMMENT '1:取引完了, 2:取引受付完了, 3:返金完了, 4:取引取消, 5:取引受付取消, 6:調整, 7:送金完了\n',
  `transaction_at` datetime DEFAULT NULL COMMENT '取引日時: PayPay側での処理日時',
  `transaction_amount` decimal(12,2) DEFAULT NULL COMMENT '取引金額: 入金の場合はマイナスで表記される場合あり',
  `receipt_number` varchar(255) DEFAULT NULL COMMENT 'レシート番号: 決済電文のレシート番号',
  `paypay_payment_method` int DEFAULT NULL COMMENT '支払い方法\\n 1:PayPay（残高）, 2:クレジットカード, 3:Yahoo!マネー廃⽌, 4:Alipay, 5:あと払い（一括のみ）, 6:プリペイドコード, 7:LinePay, 8:PayPay（クレジット）, 9:PayPay商品券, 10:PayPayポイント, 11:PayPay銀行残高',
  `ssid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'SSID: 任意で加盟店から与えられた店舗コード',
  `merchant_order_id` varchar(255) DEFAULT NULL COMMENT '加盟店注文ID: merchant_order_id または merchant_refund_id',
  `payment_detail` json DEFAULT NULL COMMENT '支払い詳細: 支払い方法ごとの追加情報をJSON形式で保存',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  `deleted_at` datetime DEFAULT NULL COMMENT '削除日時（論理削除用）',
  PRIMARY KEY (`id`),
  KEY `idx_payin_file_id` (`payin_file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='PayPay入金明細ファイル取込用テーブル';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `paypay_payin_transaction`;
-- +goose StatementEnd
