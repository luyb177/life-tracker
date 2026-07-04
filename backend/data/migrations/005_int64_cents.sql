-- 将金额从 decimal(10,2) 转为 bigint（单位：分）
-- 例如: 29.90 元 → 2990 分

-- Step 1: 转换现有数据（元 * 100 = 分）
ALTER TABLE `expense_logs` MODIFY COLUMN `amount` BIGINT NOT NULL DEFAULT 0;

-- Step 2: 如果表中有数据，执行:
-- UPDATE `expense_logs` SET `amount` = ROUND(`amount` * 100) WHERE `amount` > 0;
