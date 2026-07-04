-- 创建 tags 表（全局标签池）
CREATE TABLE IF NOT EXISTS `tags` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` DATETIME(3) NOT NULL,
    `name`       VARCHAR(50) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建 life_log_tags 关联表（需要 MySQL 8.0+ 以使用 JSON_TABLE）
-- 如果使用 MySQL 5.7，请手动拆分逗号分隔的标签后迁移
CREATE TABLE IF NOT EXISTS `life_log_tags` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `life_log_id` BIGINT UNSIGNED NOT NULL,
    `tag_id`      BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_lifelog_tag` (`life_log_id`, `tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 迁移现有逗号分隔的 tags 到关联表
-- Step 1: 提取所有不重复的标签名并插入 tags 表
INSERT IGNORE INTO `tags` (`name`, `created_at`)
SELECT DISTINCT TRIM(tag_name) AS tag_name, NOW()
FROM `life_logs`
CROSS JOIN JSON_TABLE(
    CONCAT('["', REPLACE(COALESCE(`tags`, ''), ',', '","'), '"]'),
    '$[*]' COLUMNS (tag_name VARCHAR(50) PATH '$')
) jt
WHERE COALESCE(`tags`, '') != '';

-- Step 2: 创建关联记录
INSERT INTO `life_log_tags` (`life_log_id`, `tag_id`)
SELECT DISTINCT t.id, tags.id
FROM `life_logs` t
CROSS JOIN JSON_TABLE(
    CONCAT('["', REPLACE(COALESCE(t.tags, ''), ',', '","'), '"]'),
    '$[*]' COLUMNS (tag_name VARCHAR(50) PATH '$')
) jt
JOIN `tags` ON tags.name = TRIM(jt.tag_name)
WHERE COALESCE(t.tags, '') != '';

-- Step 3: 确认迁移正确后，删除 life_logs 的 tags 列（手动执行）
-- ALTER TABLE `life_logs` DROP COLUMN `tags`;
