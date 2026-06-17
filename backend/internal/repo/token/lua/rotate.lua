-- rotate.lua
-- 原子性刷新 refresh token JTI：先验证旧 JTI，匹配则替换为新 JTI
-- KEYS[1] = refresh:{userID}
-- ARGV[1] = expected old JTI
-- ARGV[2] = new JTI
-- ARGV[3] = TTL (seconds)
-- Returns: 1 = rotated successfully, 0 = old JTI mismatch, nil = key not found

local stored = redis.call("GET", KEYS[1])

if not stored then
    return nil
end

if stored ~= ARGV[1] then
    return 0
end

redis.call("SET", KEYS[1], ARGV[2], "EX", ARGV[3])
return 1
