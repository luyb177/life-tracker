# Life Tracker Backend

个人生活记录后端服务，覆盖：

1. **认证与鉴权**
2. **支出记录与统计**
3. **生活记录 / 周期总结统一承载层**
4. **AI 日报 / 周报 / 月报 / 年报 / 人生总结**

---

## 1. 技术栈

- **框架**：go-zero REST API
- **语言**：Go
- **数据库**：MySQL + GORM
- **缓存**：Redis
- **AI**：DeepSeek
- **地理位置**：ip2region
- **鉴权**：JWT Access Token + Refresh Token 轮换
- **部署**：Docker / docker compose

---

## 2. 项目结构

```text
backend/
├── cmd/
│   └── cron/                 # 定时任务入口
├── common/
│   ├── errorx/               # 统一业务错误
│   ├── jwtx/                 # JWT 工具
│   └── respx/                # 统一响应包装
├── data/                     # ip2region 数据文件
├── etc/                      # 配置文件
├── internal/
│   ├── config/               # 配置结构体
│   ├── constvar/             # 常量定义
│   ├── handler/              # HTTP handler（goctl 生成）
│   ├── logic/                # 业务逻辑
│   │   ├── auth/             # 注册 / 登录 / 刷新令牌
│   │   ├── cron/             # AI 总结主流程
│   │   ├── expense/          # 支出 CRUD 与统计
│   │   ├── summary/          # 生活记录 / 总结
│   │   └── user/             # 用户资料 / 修改密码
│   ├── middleware/           # JWT / IP 中间件
│   ├── pkg/
│   │   ├── ai/               # AI 客户端
│   │   ├── email/            # 邮件发送
│   │   └── ip/               # IP 解析
│   ├── repo/                 # 数据访问层
│   │   ├── expense/
│   │   ├── summary/
│   │   ├── token/
│   │   ├── user/
│   │   └── verify/
│   ├── svc/                  # ServiceContext 依赖注入
│   └── types/                # API 类型（goctl 生成）
├── Dockerfile
├── life_tracker.api          # API 单一事实来源
└── life_tracker.go           # 服务入口
```

---

## 3. 架构说明

项目使用典型的 **handler → logic → repo** 分层：

1. **handler**
   - 负责解析请求、调用 logic、输出响应
   - 不写业务逻辑

2. **logic**
   - 负责参数校验、业务规则、流程编排
   - 比如 summary 周期规则、AI 总结触发、支出统计拼装

3. **repo**
   - 负责数据库 / Redis 访问
   - 不承担 HTTP 层和页面语义

4. **svc / ServiceContext**
   - 统一管理 DB、Redis、Repo、AI 配置、中间件等依赖

---

## 4. 核心业务模型

## 4.1 Summary 是统一记录层

本项目**不单独再维护 `life_logs` 表**。

`summary` 统一承载两类内容：

1. **`source=user`**
   - 用户手动写入的生活记录 / 日报 / 主动总结

2. **`source=AI`**
   - AI 自动生成的周期总结

也就是说，`summary` 既是“生活记录”，也是“AI 总结”的统一数据层。

---

## 4.2 Summary 重复规则

当前业务规则如下：

1. **`source=user + day`**
   - 同一天允许多条
   - 用于记录多个生活片段 / 多条日报内容

2. **`source=user + week/month/year/life`**
   - 同周期仅允许一条
   - 若已存在，应走 `/summary/update`

3. **`source=AI`**
   - 同用户、同周期、同 `period_start` 仅保留一条
   - 重跑 AI 时更新原记录，不新建重复记录

4. **AI 不覆盖用户记录**
   - AI 与用户记录通过 `source` 分层共存

---

## 4.3 Summary 周期模型

当前 `summary.period_start` / `summary.period_end` 已统一为：

1. **全部使用 `YYYY-MM-DD`**
2. **`period_end` 为开区间结束日期**

### 周期语义

| PeriodType | 含义 | `period_start` 规则 | `period_end` 规则 |
|---|---|---|---|
| 1 | 日报 | 任意合法日期 | `start + 1 天` |
| 2 | 周报 | 必须是周一 | `start + 7 天` |
| 3 | 月报 | 必须是当月第一天 | 下月第一天 |
| 4 | 年报 | 必须是当年 1 月 1 日 | 下一年 1 月 1 日 |
| 5 | 人生总结 | 任意合法日期 | 必须晚于 `start` |

### 示例

| 类型 | period_start | period_end |
|---|---|---|
| 日报 | `2026-06-17` | `2026-06-18` |
| 周报 | `2026-06-15` | `2026-06-22` |
| 月报 | `2026-06-01` | `2026-07-01` |
| 年报 | `2026-01-01` | `2027-01-01` |
| 人生总结 | `2020-01-01` | `2026-06-17` |

---

## 4.4 AI 总结输入来源

AI 总结不是只看支出，它会综合使用：

1. **支出汇总**
2. **分类支出明细**
3. **地点分布**
4. **用户日报内容**
5. **下级周期总结**
6. **人生总结场景下的长期标签偏好**
7. **人生总结场景下最近 12 个月的记录 / 支出 / 标签趋势**

### 下级周期关系

| 当前总结 | 使用的下级上下文 |
|---|---|
| 日报 | 无 |
| 周报 | 日报 |
| 月报 | 周报 |
| 年报 | 月报 |
| 人生总结 | 年报 |

---

## 5. 鉴权与安全

## 5.1 Access / Refresh Token

当前鉴权方案：

1. Access Token
   - 短期有效
   - 用于请求 API

2. Refresh Token
   - 长期有效
   - 用于换新 Access Token

3. Refresh Token 轮换
   - 每次刷新都会生成新的 Refresh Token
   - Redis 中只保留当前有效 JTI
   - 旧 Refresh Token 会立即失效

4. 修改密码后
   - 会清除对应 Refresh Token 状态
   - 达到强制重新登录效果

## 5.2 验证码

验证码使用 Redis 存储，并带有：

1. TTL 过期控制
2. 发送冷却
3. Lua 原子校验 / 删除，防止重复消费

---

## 6. 本地开发

```bash
# 进入 backend 目录
cd backend

# 整理依赖
go mod tidy

# 修改 life_tracker.api 后重新生成
goctl api go -api life_tracker.api -dir . --style go_zero

# 可选：生成 Swagger
goctl api swagger -api life_tracker.api -dir ./docs -filename swagger

# 运行服务
go run life_tracker.go -f etc/life_tracker.yaml
```

### 说明

1. `life_tracker.api` 是 API 合同的单一事实来源
2. 修改接口类型或路由时，优先改 `.api`
3. goctl 不会覆盖已有业务逻辑文件，生成冲突时需人工清理 stub

---

## 7. 配置说明

`etc/life_tracker.yaml` 示例：

```yaml
Name: life_tracker
Host: 0.0.0.0
Port: 8888

MySQL:
  DataSource: user:password@tcp(mysql-host:3306)/life_tracker?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

Redis:
  Addr: redis-host:6379
  Password: ""
  DB: 0

AIConf:
  Endpoint: https://api.deepseek.com
  APIKey: sk-xxxxxxxxxxxxxxxx
  Model: deepseek-chat

Auth:
  AccessSecret: your-access-secret
  AccessExpire: 900
  RefreshSecret: your-refresh-secret
  RefreshExpire: 604800
  RefreshRemember: 2592000

IP2Region:
  V4: data/ip2region_v4.xdb
  V6: data/ip2region_v6.xdb
```

### 配置要点

1. `MySQL.DataSource` 必须带 `parseTime=true`
2. `loc=Asia%2FShanghai` 要和代码中的 `TimeLocation` 一致
3. `AIConf` 必须填可用模型和密钥
4. `IP2Region` 需要本地 `.xdb` 文件

---

## 8. 构建与部署

## 8.1 构建镜像

在仓库根目录执行：

```bash
./build.sh
```

## 8.2 部署准备

```bash
mkdir -p /opt/life-tracker/data

scp docker-compose.yaml   root@server:/opt/life-tracker/
scp etc/life_tracker.yaml root@server:/opt/life-tracker/etc/
scp data/ip2region_*.xdb  root@server:/opt/life-tracker/data/
```

## 8.3 启动

```bash
cd /opt/life-tracker
docker compose up -d
```

### 运行服务

| 服务 | 说明 |
|---|---|
| `life-tracker` | 主 API 服务 |
| `life-tracker-scheduler` | AI 总结定时任务 |

---

## 9. 定时任务

调度器按固定周期自动生成 AI 总结：

| 任务 | 时间 |
|---|---|
| 日报 | 每天 08:00 |
| 周报 | 每周一 08:00 |
| 月报 | 每月 1 日 08:00 |
| 年报 | 每年 1 月 1 日 08:00 |
| 人生总结 | 每年 1 月 1 日 09:00 |

### 手动触发

```bash
# 1=日报 2=周报 3=月报 4=年报 5=人生总结
docker compose run --rm life-tracker-cron cron summary -t 1
docker compose run --rm life-tracker-cron cron summary -t 2
docker compose run --rm life-tracker-cron cron summary -t 3
docker compose run --rm life-tracker-cron cron summary -t 4
docker compose run --rm life-tracker-cron cron summary -t 5
```

---

## 10. 统一响应格式

所有接口统一通过 `common/respx` 输出：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

### 约定

1. `code = 0` 表示成功
2. 非 0 表示业务错误
3. handler 层负责写响应
4. logic 层返回 `(*Resp, error)`

---

## 11. API 总览

## 11.1 Auth

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/v1/auth/send_code` | 发送验证码 |
| POST | `/api/v1/auth/register` | 注册 |
| POST | `/api/v1/auth/login` | 登录 |
| POST | `/api/v1/auth/refresh` | 刷新令牌 |

## 11.2 User

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/v1/user/update_info` | 更新用户信息 |
| POST | `/api/v1/user/change_password` | 修改密码 |

## 11.3 Expense

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/expense/categories` | 支出分类列表 |
| POST | `/api/v1/expense/category/create` | 创建自定义分类 |
| POST | `/api/v1/expense/category/delete` | 删除自定义分类 |
| POST | `/api/v1/expense/create` | 创建支出记录 |
| POST | `/api/v1/expense/update` | 更新支出记录 |
| POST | `/api/v1/expense/delete` | 删除支出记录 |
| GET | `/api/v1/expense/list` | 支出列表（游标分页） |
| GET | `/api/v1/expense/by_date?date=YYYY-MM-DD` | 某天支出明细 + 总额 |
| GET | `/api/v1/expense/daily_total?date=YYYY-MM-DD` | 某日支出总额 |
| GET | `/api/v1/expense/stats/range?start=YYYY-MM-DD&end=YYYY-MM-DD` | 区间总额 |
| GET | `/api/v1/expense/stats/category?start=YYYY-MM-DD&end=YYYY-MM-DD` | 分类占比 |
| GET | `/api/v1/expense/stats/trend?start=YYYY-MM-DD&end=YYYY-MM-DD` | 按日趋势 |
| GET | `/api/v1/expense/stats/monthly?start=YYYY-MM-DD&end=YYYY-MM-DD` | 按月趋势 |

## 11.4 Summary

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/v1/summary/create` | 创建用户记录 / 主动总结 |
| POST | `/api/v1/summary/update` | 更新 summary |
| POST | `/api/v1/summary/delete` | 删除 summary |
| GET | `/api/v1/summary/list` | 游标分页列表 |
| GET | `/api/v1/summary/day?date=YYYY-MM-DD` | 查询某天日报（用户 + AI） |
| GET | `/api/v1/summary/range?start=YYYY-MM-DD&end=YYYY-MM-DD` | 时间范围查询 |
| GET | `/api/v1/summary/stats/tags?start=YYYY-MM-DD&end=YYYY-MM-DD` | 标签频次统计 |
| GET | `/api/v1/summary/stats/tag_trend?start=YYYY-MM-DD&end=YYYY-MM-DD` | 标签按月趋势 |
| POST | `/api/v1/summary/generate` | 手动生成 AI 总结 |

---

## 12. 常见请求示例

## 12.1 登录

```bash
curl -X POST http://127.0.0.1:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "target": "demo@example.com",
    "channel": 1,
    "password": "123456"
  }'
```

## 12.2 创建一条用户日报

```bash
curl -X POST http://127.0.0.1:8888/api/v1/summary/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "period_type": 1,
    "period_start": "2026-06-17",
    "period_end": "2026-06-18",
    "summary_content": "今天主要写代码、看房、晚饭吃得比较简单。",
    "title": "普通工作日",
    "tags": "工作,看房,晚饭"
  }'
```

## 12.3 创建一条用户周报

```bash
curl -X POST http://127.0.0.1:8888/api/v1/summary/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "period_type": 2,
    "period_start": "2026-06-15",
    "period_end": "2026-06-22",
    "summary_content": "这周节奏比较紧，支出控制还可以。",
    "title": "六月第三周周报",
    "tags": "工作,支出,复盘"
  }'
```

## 12.4 手动生成 AI 月报

```bash
curl -X POST http://127.0.0.1:8888/api/v1/summary/generate \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "period_type": 3,
    "period_start": "2026-06-01"
  }'
```

## 12.5 查询标签趋势

```bash
curl "http://127.0.0.1:8888/api/v1/summary/stats/tag_trend?start=2026-01-01&end=2026-12-31" \
  -H "Authorization: Bearer <token>"
```

---

## 13. 前端联调要点

1. `summary/day` 返回用户记录和 AI 日报混合列表，前端按 `source` 分组展示
2. 同一天多条用户日报是合法行为，不要在前端假设“每天只有一条”
3. 周/月/年/人生总结若已存在，前端应优先进入编辑态，而不是重复创建
4. `period_end` 是开区间，不是“含当天结束”
5. 标签当前以逗号分隔字符串存储，前端展示时建议自行 split / trim

---

## 14. 数据库迁移

首次部署或模型变更时执行：

```sql
ALTER TABLE expense_logs ADD COLUMN location VARCHAR(255) DEFAULT '' AFTER note;

ALTER TABLE summaries ADD COLUMN location VARCHAR(255) DEFAULT '' AFTER suggestion_content;
ALTER TABLE summaries ADD COLUMN title VARCHAR(255) DEFAULT '' AFTER suggestion_content;
ALTER TABLE summaries ADD COLUMN tags VARCHAR(500) DEFAULT '' AFTER title;
```

如果数据库里存在历史旧数据，并且旧版 `period_start` / `period_end` 不是 `YYYY-MM-DD`，建议在上线统一周期规则后补一次数据清洗。

---

## 15. 维护建议

1. 新增接口时优先修改 `life_tracker.api`
2. 新增 summary 相关能力时，优先复用当前统一周期模型
3. 不要重新引入 `YYYY-Wxx`、`YYYY-MM`、`YYYY` 这类混合存储格式
4. 新的 AI 总结能力优先在 `internal/logic/cron/summary.go` 汇总上下文
