# Life Tracker Backend

Go 后端服务，提供认证、生活记录、支出、总结、用户设置和 AI 总结生成接口。

## 技术栈

- Go 1.26
- go-zero REST
- GORM + MySQL
- Redis
- DeepSeek/OpenAI 兼容 AI 接口
- Cobra cron CLI

## 目录结构

```text
backend/
├── cmd/cron/                 # 定时任务 CLI，手动/定时生成 AI 总结
├── common/                   # errorx、respx、jwtx、database、cache 等通用能力
├── data/
│   ├── migrations/           # SQL 初始化和增量迁移
│   └── ip2region_*.xdb       # IP 地理位置库
├── docs/                     # swagger 等后端文档产物
├── etc/                      # 配置文件，真实 yaml 不提交
├── internal/
│   ├── config/               # 配置结构体
│   ├── constvar/             # 常量
│   ├── handler/              # HTTP 层，保持薄层
│   ├── logic/                # 业务流程
│   │   ├── auth/             # 注册、登录、刷新 token、验证码
│   │   ├── cron/             # AI 总结生成主流程
│   │   ├── expense/          # 支出、分类、统计
│   │   ├── lifelog/          # 生活记录
│   │   ├── summary/          # 总结
│   │   └── user/             # 用户资料、密码
│   ├── middleware/           # JWT、IP 地理位置
│   ├── pkg/                  # AI、邮件、验证码、分页 token、密码工具
│   ├── repo/                 # 数据访问层，GORM/Redis
│   ├── svc/                  # ServiceContext 依赖注入
│   └── types/                # goctl 生成类型
├── life_tracker.api          # API 定义
├── life_tracker.go           # API 服务入口
└── Dockerfile
```

## 架构约定

### Handler -> Logic -> Repo

- `handler` 只解析请求、调用 logic、写统一响应。
- `logic` 做参数校验、业务规则和事务编排。
- `repo` 封装 GORM/Redis 访问，并暴露接口。
- `svc.ServiceContext` 在启动时初始化所有依赖，失败直接 panic。

### API-first

接口和类型先改 `life_tracker.api`，再生成：

```bash
goctl api go -api life_tracker.api -dir . --style go_zero
```

注意：goctl 不覆盖已有文件，重新生成时需要先删除目标文件；如果生成了 `lifetracker.go` 且已有 `life_tracker.go`，需要删除冲突桩文件。

## 核心模型

### 用户与鉴权

- 邮箱验证码注册。
- 邮箱密码登录。
- Access token 默认短期有效。
- Refresh token 默认 7 天有效。
- Refresh token 每次刷新都会轮换，Redis 保存当前有效 JTI。
- Redis Lua 脚本保证 refresh token 校验和替换原子执行。
- 修改密码会删除 Redis refresh key，使旧 refresh token 失效。

### 生活记录

- 表：`life_logs`
- 同一天可多条，按 `occurred_at` 查询和排序。
- 支持标签，关联表为 `life_log_tags`。
- 有 `last_updated_by`、`last_updated_at` 审计字段。

### 支出

- 表：`expense_logs`
- 金额单位为分，`2990 = 29.90 元`。
- 分类表：`expense_categories`，`user_id=0` 为系统默认分类，其他为用户自定义。
- 支持退款，退款后 `status=1`，统计查询自动排除。
- 有 `last_updated_by`、`last_updated_at` 审计字段。

### 总结

- 表：`summaries`
- `period_type`：1 日报，2 周报，3 月报，4 年报。
- `period_start` / `period_end` 为 `YYYY-MM-DD` 字符串，`period_end` 是开区间。
- `source`：1 AI，2 用户。
- 同一用户、同一周期、同一来源只保留一条总结。
- 总结内容支持 Markdown。
- 有 `last_updated_by`、`last_updated_at` 审计字段；系统定时任务更新时 `last_updated_by=0`。

AI 总结上下文层级：

```text
日报 <- 支出 + 生活记录
周报 <- 日报 + 支出 + 生活记录
月报 <- 周报 + 支出 + 生活记录
年报 <- 月报 + 支出 + 生活记录
```

## 配置

复制示例：

```bash
cp etc/life_tracker.example.yaml etc/life_tracker.yaml
```

关键配置：

```yaml
MySQLConf:
  DSN: user:password@tcp(127.0.0.1:3306)/life_tracker?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

RedisConf:
  Addr: 127.0.0.1:6379

JWTConf:
  Secret: your-secret
  ExpireS: 7200
  RefreshExpireS: 604800

AIConf:
  Endpoint: https://api.deepseek.com
  APIKey: sk-xxx
  Model: deepseek-chat
```

真实 `etc/*.yaml` 已被忽略，不要提交。

## 本地开发

```bash
# 初始化数据库
mysql -u root -p life_tracker < data/migrations/000_init_schema.sql

# 同步依赖
go mod tidy

# 编译检查
go build ./...

# 启动 API
go run life_tracker.go -f etc/life_tracker.yaml

# 手动执行 AI 总结任务，1=日报 2=周报 3=月报 4=年报
go run cmd/cron/main.go cron summary -t 1
```

## 自动 AI 总结调度

生产部署使用根目录 `docker-compose.yaml` 中的 `life-tracker-scheduler` 服务。该服务使用 `crond`，并设置 `TZ=Asia/Shanghai`，调度时间按北京时间计算：

| 周期 | crontab | 说明 |
| --- | --- | --- |
| 日报 | `0 8 * * *` | 每天 08:00 |
| 周报 | `0 8 * * 1` | 每周一 08:00 |
| 月报 | `0 8 1 * *` | 每月 1 日 08:00 |
| 年报 | `0 8 1 1 *` | 每年 1 月 1 日 08:00 |

手动验证：

```bash
docker compose run --rm life-tracker-cron cron summary -t 1 --config etc/life_tracker.yaml
```

## API 分组

所有接口前缀为 `/api/v1`。

| 分组 | 说明 |
| --- | --- |
| `/auth` | 发送验证码、注册、登录、刷新 token |
| `/user` | 修改资料、修改密码 |
| `/life_log` | 生活记录 CRUD、按日查询、分页查询 |
| `/expense` | 支出 CRUD、退款、分类、趋势和统计 |
| `/summary` | 总结 CRUD、范围查询、AI 生成 |

## 迁移

迁移文件在 `data/migrations/`。

- `000_init_schema.sql`：初始化表结构。
- `001_add_last_updated_audit.sql`：补充审计字段和总结唯一约束。

如果已有重复总结数据，执行唯一索引迁移前需要先合并或清理重复数据。

## 维护要点

- 金额一律使用分，避免浮点误差。
- 时间按 `Asia/Shanghai` 处理，MySQL DSN 必须带 `loc=Asia%2FShanghai`。
- 生活记录和总结是独立模块，不要混用表语义。
- 支出和总结列表使用游标分页，不使用 offset。
- 手动改 goctl 生成文件前，先确认 `.api` 是否也需要同步。
- 当前项目没有 `*_test.go`，主要验证方式为 `go build ./...`。
