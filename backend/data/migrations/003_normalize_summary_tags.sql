-- 创建 summary_tags 关联表
CREATE TABLE IF NOT EXISTS `summary_tags` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `summary_id` BIGINT UNSIGNED NOT NULL,
    `tag_id`     BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_summary_tag` (`summary_id`, `tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 迁移现有 summaries 表的逗号分隔 tags 到关联表（需要 MySQL 8.0+ JSON_TABLE）
-- Step 1: 将 summaries 中的标签名插入全局 tags 表
INSERT IGNORE INTO `tags` (`name`, `created_at`)
SELECT DISTINCT TRIM(jt.tag_name), NOW()
FROM `summaries`
CROSS JOIN JSON_TABLE(
    CONCAT('["', REPLACE(COALESCE(`tags`, ''), ',', '","'), '"]'),
    '$[*]' COLUMNS (tag_name VARCHAR(50) PATH '$')
) jt
WHERE COALESCE(`tags`, '') != '';

-- Step 2: 创建 summary_tags 关联记录
INSERT INTO `summary_tags` (`summary_id`, `tag_id`)
SELECT DISTINCT s.id, tags.id
FROM `summaries` s
CROSS JOIN JSON_TABLE(
    CONCAT('["', REPLACE(COALESCE(s.tags, ''), ',', '","'), '"]'),
    '$[*]' COLUMNS (tag_name VARCHAR(50) PATH '$')
) jt
JOIN `tags` ON tags.name = TRIM(jt.tag_name)
WHERE COALESCE(s.tags, '') != '';

-- Step 3: 确认迁移正确后删除 summaries 的 tags 列（手动执行）
-- ALTER TABLE `summaries` DROP COLUMN `tags`;
