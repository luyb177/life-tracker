# Life Tracker Backend

个人生活追踪后端服务 — 支出记录 + AI 智能总结。

## 技术栈

- **框架**: go-zero (REST API)
- **数据库**: MySQL + GORM
- **缓存**: Redis (JWT refresh token / 验证码)
- **AI**: DeepSeek (日报/周报/月报/年报自动生成)
- **IP 定位**: ip2region (自动记录消费地点)
- **构建**: Docker multi-stage build

## 项目结构

```
backend/
├── cmd/cron/              # 定时任务入口 (cobra CLI)
├── internal/
│   ├── config/            # 配置结构体
│   ├── handler/           # HTTP 处理器 (goctl 生成)
│   ├── logic/             # 业务逻辑
│   │   ├── cron/          # AI 总结生成
│   │   ├── expense/       # 支出管理
│   │   ├── summary/       # 总结管理
│   │   └── user/          # 用户认证
│   ├── middleware/         # JWT / IP 定位中间件
│   ├── repo/              # 数据访问层
│   │   ├── expense/       # 支出 + 分类
│   │   ├── summary/       # AI 总结
│   │   ├── token/         # Refresh token
│   │   ├── user/          # 用户
│   │   └── verify/        # 验证码 (Redis + Lua)
│   ├── svc/               # ServiceContext (DI)
│   └── types/             # API 类型定义 (goctl 生成)
├── data/                  # ip2region 数据库
├── etc/                   # 配置文件
├── Dockerfile             # 构建镜像
└── life_tracker.api       # API 定义文件
```

## 本地开发

```bash
# 同步依赖
go mod tidy

# 生成代码（修改 .api 文件后执行）
goctl api go -api life_tracker.api -dir . --style go_zero

# 生成 Swagger 文档
goctl api swagger -api life_tracker.api -dir ./docs -filename swagger

# 运行
go run life_tracker.go -f etc/life_tracker.yaml
```

## 构建镜像

```bash
# 在仓库根目录执行
./build.sh
```

## 部署

### 1. 服务器准备

```bash
mkdir -p /opt/life-tracker/data
# 复制文件到服务器
scp docker-compose.yaml       root@server:/opt/life-tracker/
scp etc/life_tracker.yaml     root@server:/opt/life-tracker/etc/
scp data/ip2region_*.xdb      root@server:/opt/life-tracker/data/
```

### 2. 配置 `etc/life_tracker.yaml`

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
  Endpoint: https://api.deepseek.com/v1
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

### 3. 一键启动

```bash
cd /opt/life-tracker
docker compose up -d
```

启动后自动运行以下服务：

| 服务 | 说明 |
|---|---|
| `life-tracker` | 主 API 服务，端口 8888 |
| `life-tracker-scheduler` | 定时调度器 |

### 4. 定时任务

`life-tracker-scheduler` 容器内置 crond，自动按以下时间执行 AI 总结：

| 任务 | 时间 |
|---|---|
| 日报 | 每天 8:00 |
| 周报 | 每周一 8:00 |
| 月报 | 每月 1 号 8:00 |
| 年报 | 每年 1 月 1 号 8:00 |

### 5. 手动触发总结

```bash
# 日报
docker compose run --rm life-tracker-cron cron summary -t 1

# 周报
docker compose run --rm life-tracker-cron cron summary -t 2

# 月报
docker compose run --rm life-tracker-cron cron summary -t 3

# 年报
docker compose run --rm life-tracker-cron cron summary -t 4
```

## API 接口

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /api/user/register | 注册 |
| POST | /api/user/login | 登录 |
| POST | /api/user/send-code | 发送验证码 |
| POST | /api/user/change-password | 修改密码 |
| PUT | /api/user/info | 更新用户信息 |
| POST | /api/expense/category | 创建支出分类 |
| GET | /api/expense/category/list | 分类列表 |
| DELETE | /api/expense/category | 删除分类 |
| POST | /api/expense/log | 创建支出记录 |
| GET | /api/expense/log/list | 支出列表（游标分页） |
| DELETE | /api/expense/log | 删除支出记录 |
| GET | /api/expense/daily-total | 当日支出汇总 |
| POST | /api/summary | 创建总结 |
| GET | /api/summary/list | 总结列表 |
| PUT | /api/summary | 更新总结 |
| DELETE | /api/summary | 删除总结 |
| POST | /api/summary/generate | 手动生成 AI 总结 |

## 数据库迁移

首次部署或模型变更时执行：

```sql
-- Location 字段（v0.0.1 新增）
ALTER TABLE expense_logs ADD COLUMN location VARCHAR(255) DEFAULT '' AFTER note;
ALTER TABLE summaries ADD COLUMN location VARCHAR(255) DEFAULT '' AFTER suggestion_content;
```
