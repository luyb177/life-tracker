# Life Tracker Backend

个人生活记录后端服务：

1. **认证与鉴权**
2. **生活记录** — 日常活动记录（life_logs）
3. **支出记录** — 支出 CRUD + 分类管理 + 退款 + 统计
4. **周期总结** — AI 日报 / 周报 / 月报 / 年报 + 用户手动总结

---

## 1. 技术栈

- **框架**：go-zero REST API
- **语言**：Go 1.26
- **数据库**：MySQL 8.0 + GORM
- **缓存**：Redis
- **AI**：DeepSeek（OpenAI 兼容接口）
- **地理位置**：ip2region
- **鉴权**：JWT Access Token + Refresh Token 轮换（Redis Lua 防重放）
- **部署**：Docker / docker compose

---

## 2. 项目结构

```text
backend/
├── cmd/cron/                    # 定时任务 CLI（Cobra）
├── common/
│   ├── cache/                   # Redis 客户端
│   ├── database/                # MySQL 客户端
│   ├── errorx/                  # 统一错误码 + AppError
│   ├── jwtx/                    # JWT 签发 / 校验
│   ├── mail/                    # SMTP 邮件
│   └── respx/                   # 统一响应 {code, msg, data}
├── data/
│   └── migrations/              # 数据库迁移 SQL
├── docs/
│   └── schema.md                # 数据模型 ↔ SQL 对照
├── etc/                         # 配置文件
├── internal/
│   ├── config/                  # 配置结构体
│   ├── constvar/                # 业务常量
│   ├── handler/                 # HTTP 层（goctl 生成，薄层）
│   ├── logic/                   # 业务逻辑
│   │   ├── auth/                #   注册 / 登录 / 刷新令牌
│   │   ├── cron/                #   AI 总结生成流程
│   │   ├── expense/             #   支出 CRUD + 分类 + 统计
│   │   ├── lifelog/             #   生活记录 CRUD
│   │   ├── summary/             #   周期总结 CRUD
│   │   └── user/                #   用户资料 / 改密码
│   ├── middleware/              # JWT 鉴权 + IP 地理位置
│   ├── pkg/
│   │   ├── ai/                  #   DeepSeek 客户端
│   │   ├── code/                #   验证码生成
│   │   ├── email/               #   邮件模板 + 发送
│   │   ├── pagetoken/           #   游标分页 token 编解码
│   │   └── password/            #   bcrypt
│   ├── repo/                    # 数据访问层
│   │   ├── expense/
│   │   ├── lifelog/
│   │   ├── summary/
│   │   ├── tag/                 #   标签（全局标签池 + 关联表）
│   │   ├── token/               #   Redis refresh token
│   │   ├── user/
│   │   └── verify/              #   Redis 验证码
│   ├── svc/                     # ServiceContext 依赖注入
│   └── types/                   # API 类型（goctl 从 .api 生成）
├── Dockerfile
├── life_tracker.api             # API 定义的单一事实来源
└── life_tracker.go              # 主服务入口
```

---

## 3. 架构

### 3.1 Handler → Logic → Repo

- **handler** — 解析请求，调用 logic，写 `respx` 响应。不写业务逻辑。
- **logic** — 参数校验、业务规则、流程编排。返回 `(*Resp, error)`。
- **repo** — 数据库 / Redis 访问接口 + 实现。不承担 HTTP 语义。

### 3.2 ServiceContext

`internal/svc/service_context.go` 在启动时统一初始化 MySQL、Redis、Mailer、JWT、Repos、中间件，启动失败直接 panic。

### 3.3 API-first

`life_tracker.api` 是类型的**单一事实来源**。新增/修改接口时先改 `.api`，再 `goctl api go` 生成 `types.go`、`routes.go` 和 handler stub。goctl 不会覆盖已有文件。

```bash
cd backend
goctl api go -api life_tracker.api -dir . --style go_zero
```

---

## 4. 核心数据模型

### 4.1 生活记录（life_logs）

独立表，与总结完全分离。同一天允许多条，按 `occurred_at` 排序。

| 字段 | 说明 |
|------|------|
| content | 记录内容（最长 10000 字符） |
| occurred_at | 事件发生时间 |
| tags | 通过 `life_log_tags` → `tags` 关联表（`[]TagInfo`） |

### 4.2 支出记录（expense_logs）

| 字段 | 说明 |
|------|------|
| amount | **单位：分**（int64，2990 = 29.90 元） |
| category_id | 分类（系统默认 user_id=0 或用户自定义） |
| status | 0=正常, 1=已退款（退款后统计自动排除） |
| refunded_at | 退款时间 |
| note | 备注（`*string`，可传 `""` 清空） |

### 4.3 支出分类（expense_categories）

系统默认分类 `user_id=0` 全局共享，用户自定义分类 `user_id>0` 按用户隔离。

唯一索引 `(user_id, name, deleted_at)` — 同一用户不能有同名活跃分类，软删后可重建。

### 4.4 周期总结（summaries）

| 字段 | 说明 |
|------|------|
| period_type | 1=日报, 2=周报, 3=月报, 4=年报 |
| period_start/end | YYYY-MM-DD，end 为开区间 |
| source | 1=AI 生成, 2=用户手动 |
| tags | 通过 `summary_tags` → `tags` 关联表 |

**周期约束**（用户手动创建时放宽，AI 使用固定规则）：

| 类型 | period_start | period_end | 规则 |
|------|-------------|-----------|------|
| 日报 | 任意日期 | `end > start` | 仅需 end 晚于 start |
| 周报 | 任意日期 | `end > start` | 不限制必须是周一 |
| 月报 | 任意日期 | `end > start` | 不限制必须是 1 号 |
| 年报 | 任意日期 | `end > start` | 不限制必须是 1 月 1 日 |

**去重**：同 `(user_id, period_type, period_start, source)` 仅保留一条，AI 重跑时更新。

### 4.5 标签（tags）

全局标签池（`tags` 表），所有用户共享。类似小红书的话题/#机制：

- 前端传 `tags: [{id: 0, name: "新标签"}]` → 后端自动 `FindOrCreate`
- 前端传 `tags: [{id: 1, name: "工作"}]` → 后端校验后关联
- `life_log_tags` + `summary_tags` 两个关联表

---

## 5. 鉴权

- **Access Token** — 短期有效（默认 2h），Bearer 方式传递
- **Refresh Token** — 长期有效（默认 7d），含 JTI
- **Token 轮换** — Redis 存当前有效 JTI，Lua 脚本原子性替换，旧 Token 立即失效
- **改密码** — 清除 Redis 中的 refresh token，强制重新登录
- **验证码** — Redis 存储，SHA256(email) 做 key，Lua 原子校验+删除防重放

---

## 6. 配置

`etc/life_tracker.yaml`：

```yaml
Name: life_tracker
Host: 0.0.0.0
Port: 8888

MySQLConf:
  DSN: user:password@tcp(mysql:3306)/life_tracker?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

RedisConf:
  Addr: redis:6379
  Password: ""
  DB: 0

JWTConf:
  Secret: your-secret
  ExpireS: 7200           # Access Token 2h
  RefreshExpireS: 604800  # Refresh Token 7d

AIConf:
  Endpoint: https://api.deepseek.com
  APIKey: sk-xxx
  Model: deepseek-chat

EmailConf:
  From: noreply@example.com
  Password: xxx
  SMTPHost: smtp.example.com
  SMTPPort: 587

IP2RegionConf:
  V4: data/ip2region_v4.xdb
  V6: data/ip2region_v6.xdb
```

---

## 7. 本地开发

```bash
cd backend

# 依赖
go mod tidy

# 生成代码（修改 .api 后）
goctl api go -api life_tracker.api -dir . --style go_zero

# 运行
go run life_tracker.go -f etc/life_tracker.yaml

# 手动 cron
go run cmd/cron/main.go cron summary -t 1 --config etc/life_tracker.yaml
```

---

## 8. API 总览

### Auth

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/send_code` | 发送验证码 |
| POST | `/api/v1/auth/register` | 注册 |
| POST | `/api/v1/auth/login` | 登录 |
| POST | `/api/v1/auth/refresh` | 刷新令牌 |

### User

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/user/update_info` | 更新用户信息 |
| POST | `/api/v1/user/change_password` | 修改密码 |

### Expense

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/expense/categories` | 分类列表（系统+自定义） |
| POST | `/expense/category/create` | 创建自定义分类 |
| POST | `/expense/category/delete` | 删除自定义分类 |
| POST | `/expense/create` | 记录支出（amount 单位分） |
| POST | `/expense/update` | 更新支出（可清空备注） |
| POST | `/expense/delete` | 删除支出 |
| POST | `/expense/refund` | 退款（status→1，统计排除） |
| GET | `/expense/list?page_size=N` | 支出列表（游标分页） |
| GET | `/expense/by_date?date=YYYY-MM-DD` | 某天明细 + 总额 |
| GET | `/expense/stats/category?start=&end=` | 分类支出占比 |
| GET | `/expense/stats/trend?start=&end=` | 每日支出趋势 |
| GET | `/expense/stats/monthly?start=&end=` | 每月支出趋势 |

### LifeLog

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/life_log/create` | 创建生活记录 |
| POST | `/life_log/update` | 更新生活记录 |
| POST | `/life_log/delete` | 删除生活记录 |
| GET | `/life_log/list?page_size=N` | 列表（支持 `tag_id` 过滤） |
| GET | `/life_log/by_date?date=YYYY-MM-DD` | 按天查询 |

### Summary

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/summary/create` | 创建用户总结 |
| POST | `/summary/update` | 更新总结 |
| POST | `/summary/delete` | 删除总结 |
| POST | `/summary/generate` | 手动触发 AI 总结 |
| GET | `/summary/list?page_size=N` | 总结列表 |
| GET | `/summary/range?start=&end=` | 时间范围查询 |

---

## 9. curl 示例

### 注册 / 登录

```bash
# 注册
curl -X POST :8888/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"target":"demo@example.com","channel":1,"code":"123456","password":"123456"}'

# 登录
curl -X POST :8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"target":"demo@example.com","channel":1,"password":"123456"}'
```

### 生活记录

```bash
# 创建（带标签）
curl -X POST :8888/api/v1/life_log/create \
  -H "Authorization: Bearer <token>" \
  -d '{"content":"今天写了代码","tags":[{"id":0,"name":"工作"},{"id":0,"name":"学习"}],"occurred_at":"2026-07-04 10:00:00"}'

# 按天查询
curl ":8888/api/v1/life_log/by_date?date=2026-07-04" \
  -H "Authorization: Bearer <token>"
```

### 支出记录

```bash
# 记录支出（amount 单位分，2990 = 29.90 元）
curl -X POST :8888/api/v1/expense/create \
  -H "Authorization: Bearer <token>" \
  -d '{"category_id":1,"amount":2990,"note":"午饭","occurred_at":"2026-07-04 12:00:00"}'

# 退款
curl -X POST :8888/api/v1/expense/refund \
  -H "Authorization: Bearer <token>" \
  -d '{"id":1}'

# 某天统计
curl ":8888/api/v1/expense/by_date?date=2026-07-04" \
  -H "Authorization: Bearer <token>"
```

### AI 总结

```bash
# 手动生成日报
curl -X POST :8888/api/v1/summary/generate \
  -H "Authorization: Bearer <token>" \
  -d '{"period_type":1,"period_start":"2026-07-03"}'
```

---

## 10. 数据库

初始化 DDL 和迁移脚本在 `data/migrations/`。首次部署执行：

```bash
mysql -u root -p life_tracker < data/migrations/000_init_schema.sql
```

完整的数据模型 ↔ SQL 对照见 `docs/schema.md`。

---

## 11. 部署

```bash
./build.sh                    # 构建 + 推送镜像
docker compose up -d          # 启动（API + scheduler）
docker compose run --rm life-tracker-cron cron summary -t 1  # 手动 cron
```

---

## 12. 维护要点

1. **接口变更**：优先改 `life_tracker.api`，再 `goctl api go`
2. **life_logs ≠ summaries**：生活记录和周期总结是完全独立的模块
3. **金额单位**：全链路使用 int64 分（2990 = 29.90 元），前端负责 `/100` 显示
4. **标签**：全局共享，创建时 id=0 自动创建新标签，id>0 校验后关联
5. **goctl 后**：`Note` 字段需要手动改成 `*string`（goctl 不支持指针类型）
