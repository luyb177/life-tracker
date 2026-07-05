-- ============================================================================
-- Life Tracker 数据库初始化 DDL
-- MySQL 8.0+ / InnoDB / utf8mb4
--
-- 软删除约定:
--   deleted_at = 0 → 活跃记录
--   deleted_at > 0 → 已删除（Unix nano 时间戳）
--
-- 审计字段:
--   last_updated_by → 最后更新人 (0=系统/AI)
--   last_updated_at → 最后更新时间
-- ============================================================================

-- ############################################################################
-- 1. users — 用户表
-- ############################################################################
-- email  唯一，作为登录凭证
-- password 存储 bcrypt hash
-- avatar  头像 URL
-- name    用户名
-- ############################################################################
CREATE TABLE `users` (
    `id`         bigint unsigned auto_increment,
    `created_at` datetime(3)  NULL,
    `updated_at` datetime(3)  NULL,
    `deleted_at` bigint unsigned,
    `avatar`     varchar(512),
    `name`       varchar(100),
    `email`      varchar(255),
    `password`   varchar(255),
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_users_email` (`email`)
);

-- ############################################################################
-- 2. expense_categories — 支出分类表
-- ############################################################################
-- user_id = 0 → 系统全局默认分类，所有用户可见
-- user_id > 0 → 用户自定义分类，仅该用户可见
-- type:       1=系统默认  2=用户自定义
--
-- 唯一索引 (user_id, name, deleted_at):
--   同一用户不能有同名活跃分类（deleted_at=0），软删后可以重建同名
-- ############################################################################
CREATE TABLE `expense_categories` (
    `id`         bigint unsigned auto_increment,
    `created_at` datetime(3)  NULL,
    `updated_at` datetime(3)  NULL,
    `deleted_at` bigint unsigned,
    `user_id`    bigint unsigned,
    `name`       varchar(50),
    `type`       tinyint unsigned,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_user_name_del` (`user_id`, `name`, `deleted_at`)
);

-- 初始化系统默认分类（user_id=0, type=1, deleted_at=0）
INSERT IGNORE INTO `expense_categories` (`user_id`, `name`, `type`, `created_at`, `updated_at`, `deleted_at`) VALUES
    (0, '早饭', 1, NOW(), NOW(), 0),
    (0, '午饭', 1, NOW(), NOW(), 0),
    (0, '晚饭', 1, NOW(), NOW(), 0),
    (0, '杂项', 1, NOW(), NOW(), 0);

-- ############################################################################
-- 3. expense_logs — 支出记录表
-- ############################################################################
-- amount       单位：分 (2990 = 29.90 元)
-- location     由 IP 中间件自动填充（格式: 国家/省份/城市）
-- status       0=正常  1=已退款 (退款后统计自动排除)
-- refunded_at  退款时间 (status=1 时非空)
-- last_updated_by  最后更新人 (0=系统)
-- last_updated_at  最后更新时间
-- ############################################################################
CREATE TABLE `expense_logs` (
    `id`              bigint unsigned auto_increment,
    `created_at`      datetime(3)  NULL,
    `updated_at`      datetime(3)  NULL,
    `deleted_at`      bigint unsigned,
    `user_id`         bigint unsigned,
    `category_id`     bigint unsigned,
    `amount`          bigint,
    `note`            varchar(255),
    `location`        varchar(255),
    `occurred_at`     datetime,
    `status`          tinyint unsigned DEFAULT 0,
    `refunded_at`     datetime,
    `last_updated_by` bigint unsigned DEFAULT 0,
    `last_updated_at` datetime(3),
    PRIMARY KEY (`id`),
    INDEX `idx_user_date`                 (`user_id`, `occurred_at`),
    INDEX `idx_expense_logs_category_id`  (`category_id`)
);

-- ############################################################################
-- 4. summaries — 周期总结表
-- ############################################################################
-- period_type        1=日报  2=周报  3=月报  4=年报
-- period_start/end   统一 YYYY-MM-DD 格式，end 为开区间
-- source             1=AI 生成  2=用户手动
-- summary_content    总结正文
-- suggestion_content 改进建议
-- location           AI 生成时汇总的地点分布文本
-- last_updated_by    最后更新人 (0=系统/AI, >0=用户ID)
-- last_updated_at    最后更新时间
--
-- 去重规则: 同 (user_id, period_type, period_start, source) 仅保留一条
-- AI 重跑时更新原记录，不新建
-- ############################################################################
CREATE TABLE `summaries` (
    `id`                 bigint unsigned auto_increment,
    `created_at`         datetime(3)  NULL,
    `updated_at`         datetime(3)  NULL,
    `deleted_at`         bigint unsigned,
    `user_id`            bigint unsigned,
    `period_type`        tinyint unsigned,
    `period_start`       varchar(32),
    `period_end`         varchar(32),
    `source`             tinyint unsigned,
    `summary_content`    text,
    `suggestion_content` text,
    `title`              varchar(255),
    `location`           varchar(255),
    `status`             tinyint unsigned DEFAULT 1,
    `last_updated_by`    bigint unsigned DEFAULT 0,
    `last_updated_at`    datetime(3),
    PRIMARY KEY (`id`),
    INDEX `idx_user_period` (`user_id`, `period_type`, `period_start`)
);

-- ############################################################################
-- 5. life_logs — 生活记录表
-- ############################################################################
-- 用户日常活动记录，同一天允许多条
-- content   记录内容（logic 层限制最长 10000 字符）
-- 标签      通过 life_log_tags → tags 关联表管理，不存于此表
-- last_updated_by  最后更新人 (0=系统)
-- last_updated_at  最后更新时间
-- ############################################################################
CREATE TABLE `life_logs` (
    `id`              bigint unsigned auto_increment,
    `created_at`      datetime(3)  NULL,
    `updated_at`      datetime(3)  NULL,
    `deleted_at`      bigint unsigned,
    `user_id`         bigint unsigned,
    `content`         text,
    `occurred_at`     datetime,
    `last_updated_by` bigint unsigned DEFAULT 0,
    `last_updated_at` datetime(3),
    PRIMARY KEY (`id`),
    INDEX `idx_user_date` (`user_id`, `occurred_at`)
);

-- ############################################################################
-- 6. tags — 全局标签表
-- ############################################################################
-- 所有用户共享同一套标签（类似小红书的话题/#标签机制）
-- name 唯一，通过 FindOrCreate 自动扩充
-- 创建记录时前端传 [{id:0, name:"新标签"}] 即可自动创建
-- ############################################################################
CREATE TABLE `tags` (
    `id`         bigint unsigned auto_increment,
    `created_at` datetime(3)  NULL,
    `name`       varchar(50),
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_tags_name` (`name`)
);

-- ############################################################################
-- 7. life_log_tags — 生活记录 ↔ 标签关联表
-- ############################################################################
-- 多对多关联：一条生活记录可以有多个标签
-- 唯一索引 (life_log_id, tag_id) 保证不重复关联同一标签
-- tag_id 单独索引用于按标签反查生活记录 (FindLifeLogIDsByTagID)
-- ############################################################################
CREATE TABLE `life_log_tags` (
    `id`          bigint unsigned auto_increment,
    `life_log_id` bigint unsigned,
    `tag_id`      bigint unsigned,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_lifelog_tag`     (`life_log_id`, `tag_id`),
    INDEX        `idx_lifelog_tag_tag` (`tag_id`)
);

-- ############################################################################
-- 8. summary_tags — 总结 ↔ 标签关联表
-- ############################################################################
-- 同上模式，用于总结的标签关联
-- ############################################################################
CREATE TABLE `summary_tags` (
    `id`         bigint unsigned auto_increment,
    `summary_id` bigint unsigned,
    `tag_id`     bigint unsigned,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_summary_tag`     (`summary_id`, `tag_id`),
    INDEX        `idx_summary_tag_tag` (`tag_id`)
);
