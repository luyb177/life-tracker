-- Add explicit content audit fields and enforce one summary per user/period/source.
-- If an existing database already has duplicate summaries for the same
-- (user_id, period_type, period_start, source, deleted_at), merge or delete
-- duplicates before creating the unique index.

ALTER TABLE `summaries`
    ADD COLUMN `last_updated_by` bigint unsigned DEFAULT 0 AFTER `status`,
    ADD COLUMN `last_updated_at` datetime(3) NULL AFTER `last_updated_by`;

UPDATE `summaries`
SET `last_updated_by` = 0,
    `last_updated_at` = COALESCE(`updated_at`, `created_at`)
WHERE `last_updated_at` IS NULL;

CREATE UNIQUE INDEX `idx_user_period_source_live`
    ON `summaries` (`user_id`, `period_type`, `period_start`, `source`, `deleted_at`);

ALTER TABLE `life_logs`
    ADD COLUMN `last_updated_by` bigint unsigned DEFAULT 0 AFTER `occurred_at`,
    ADD COLUMN `last_updated_at` datetime(3) NULL AFTER `last_updated_by`;

UPDATE `life_logs`
SET `last_updated_by` = `user_id`,
    `last_updated_at` = COALESCE(`updated_at`, `created_at`)
WHERE `last_updated_at` IS NULL;

ALTER TABLE `expense_logs`
    ADD COLUMN `last_updated_by` bigint unsigned DEFAULT 0 AFTER `refunded_at`,
    ADD COLUMN `last_updated_at` datetime(3) NULL AFTER `last_updated_by`;

UPDATE `expense_logs`
SET `last_updated_by` = `user_id`,
    `last_updated_at` = COALESCE(`updated_at`, `created_at`)
WHERE `last_updated_at` IS NULL;
