-- 将系统默认分类统一为 user_id=0 的全局记录

-- Step 1: 插入全局系统默认分类（如果还不存在）
INSERT IGNORE INTO `expense_categories` (`user_id`, `name`, `type`, `created_at`, `updated_at`)
VALUES (0, '早饭', 1, NOW(), NOW()),
       (0, '午饭', 1, NOW(), NOW()),
       (0, '晚饭', 1, NOW(), NOW()),
       (0, '杂项', 1, NOW(), NOW());

-- Step 2: 删除所有之前为每个用户创建的旧系统默认分类
DELETE FROM `expense_categories`
WHERE `type` = 1 AND `user_id` != 0;

-- Step 3: 如果有支出记录引用了已删除的旧分类，更新 category_id 指向新的全局分类
-- 早饭
UPDATE `expense_logs` el
JOIN `expense_categories` ec_old ON ec_old.id = el.category_id
JOIN `expense_categories` ec_new ON ec_new.name = ec_old.name AND ec_new.user_id = 0
SET el.category_id = ec_new.id
WHERE ec_old.user_id != 0 AND ec_old.type = 1;

-- 午饭
UPDATE `expense_logs` el
JOIN `expense_categories` ec_old ON ec_old.id = el.category_id
JOIN `expense_categories` ec_new ON ec_new.name = ec_old.name AND ec_new.user_id = 0
SET el.category_id = ec_new.id
WHERE ec_old.user_id != 0 AND ec_old.type = 1 AND el.category_id NOT IN (SELECT id FROM expense_categories WHERE user_id = 0);

-- 上面的两步 UPDATE 已在 Step 3 中统一处理，以下 DELETE 确保旧分类引用已清理

-- Step 4: 添加唯一索引（如果表已有数据，需要先确认无冲突）
-- ALTER TABLE `expense_categories` ADD UNIQUE INDEX `idx_user_name_del` (`user_id`, `name`, `deleted_at`);
