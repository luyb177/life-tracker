-- 创建 life_logs 表
CREATE TABLE IF NOT EXISTS `life_logs` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`  DATETIME(3) NOT NULL,
    `updated_at`  DATETIME(3) NOT NULL,
    `deleted_at`  BIGINT UNSIGNED DEFAULT NULL,
    `user_id`     BIGINT UNSIGNED NOT NULL,
    `content`     TEXT NOT NULL,
    `tags`        VARCHAR(500) DEFAULT '',
    `occurred_at` DATETIME(3) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_user_date` (`user_id`, `occurred_at`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 将现有的 source=2, period_type=1 的 summaries 迁移到 life_logs
INSERT INTO `life_logs` (`user_id`, `content`, `tags`, `occurred_at`, `created_at`, `updated_at`, `deleted_at`)
SELECT `user_id`, `summary_content`, `tags`,
       STR_TO_DATE(`period_start`, '%Y-%m-%d'),
       `created_at`, `updated_at`, `deleted_at`
FROM `summaries`
WHERE `source` = 2 AND `period_type` = 1;

-- 删除已迁移的 summary 记录
DELETE FROM `summaries` WHERE `source` = 2 AND `period_type` = 1;
