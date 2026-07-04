-- 添加退款状态字段
ALTER TABLE `expense_logs`
    ADD COLUMN `status`      TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=正常, 1=已退款' AFTER `location`,
    ADD COLUMN `refunded_at` DATETIME(3)       NULL     DEFAULT NULL COMMENT '退款时间' AFTER `status`;
