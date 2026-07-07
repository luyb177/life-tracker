# Life Tracker

Life Tracker 是一个个人生活记录与支出复盘应用。它把日常记录、支出明细、趋势图表和 AI 日报/周报/月报/年报放在同一个工作台里，支持桌面端和移动端 Web 使用。

## 当前能力

- 邮箱验证码注册、邮箱密码登录。
- JWT access token + refresh token 轮换，前端支持定时刷新和 401 自动刷新重试。
- 今日页：快捷生活记录、快捷支出、当日时间线、AI 日报预览。
- 生活记录：按日期查看、编辑、删除，支持标签。
- 支出记录：按日期查看、编辑、删除、退款，自定义分类，日/月/年趋势图。
- 总结：AI 生成与用户手写总结，支持 Markdown 渲染、编辑、删除。
- 设置：修改资料、头像 URL/默认头像、修改密码、退出登录。

## 技术栈

| 层级 | 技术 |
| --- | --- |
| 前端 | Vue 3 + TypeScript + Vite + Naive UI + Pinia + ECharts |
| 后端 | Go 1.26 + go-zero + GORM |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis |
| AI | DeepSeek/OpenAI 兼容接口 |
| 鉴权 | JWT access token + refresh token，Redis Lua 原子轮换 |

## 项目结构

```text
life-tracker/
├── backend/                 # Go 后端服务
│   ├── cmd/cron/            # AI 总结定时任务 CLI
│   ├── common/              # 通用错误、响应、JWT、数据库等
│   ├── data/                # migration、ip2region 数据
│   ├── etc/                 # 配置文件，真实 yaml 不提交
│   ├── internal/            # handler -> logic -> repo
│   ├── life_tracker.api     # go-zero API 定义
│   └── life_tracker.go      # API 服务入口
├── frontend/                # Vue 前端
│   ├── src/api/             # API 封装
│   ├── src/components/      # 可复用组件
│   ├── src/pages/           # 页面
│   ├── src/stores/          # Pinia 状态
│   └── src/utils/           # 日期、金额、Markdown、总结工具
├── docs/                    # 需求文档
├── docker-compose.yaml      # 部署编排
├── build.sh                 # 后端 + 前端镜像构建脚本
└── README.md
```

## 快速启动

### 1. 后端依赖

准备 MySQL 和 Redis。数据库 DSN 需要包含 `loc=Asia%2FShanghai`。

```bash
mysql -u root -p life_tracker < backend/data/migrations/000_init_schema.sql
```

复制配置并填入真实值：

```bash
cp backend/etc/life_tracker.example.yaml backend/etc/life_tracker.yaml
```

启动后端：

```bash
cd backend
go run life_tracker.go -f etc/life_tracker.yaml
```

### 2. 前端

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器默认在 `5173`，`/api/v1` 请求代理到后端。

## 常用命令

```bash
# 后端编译
cd backend && go build ./...

# 后端从 API 定义生成代码
cd backend && goctl api go -api life_tracker.api -dir . --style go_zero

# 手动执行 AI 总结任务
cd backend && go run cmd/cron/main.go cron summary -t 1

# 前端类型检查 + 构建
cd frontend && npm run build
```

## Docker 部署

`build.sh` 会同时构建并推送两份镜像：

- `life-tracker`：Go API 服务与 cron 命令。
- `life-tracker-web`：已内置前端 `dist` 和 nginx 配置的 Web 镜像。

本地构建并推送：

```bash
./build.sh
```

服务器更新：

```bash
docker compose pull
docker compose up -d
```

前端已经打进 `life-tracker-web` 镜像，不需要再手动复制 `frontend/dist`，也不需要在服务器运行 `npm run dev`。

## 自动 AI 总结

`docker-compose.yaml` 默认启动 `life-tracker-scheduler`，容器时区设置为 `Asia/Shanghai`，会按北京时间自动执行：

- 每天 08:00 生成日报。
- 每周一 08:00 生成周报。
- 每月 1 日 08:00 生成月报。
- 每年 1 月 1 日 08:00 生成年报。

手动验证：

```bash
docker compose run --rm life-tracker-cron cron summary -t 1 --config etc/life_tracker.yaml
```

## 数据与约定

- 金额全链路用 int64 分，前端展示时转换为元。
- 生活记录和总结是两个独立模块。
- 总结按 `user_id + period_type + period_start + source` 去重：AI 一条、用户一条。
- 记录、支出、总结都有 `last_updated_by` 和 `last_updated_at` 审计字段。
- `period_end` 是开区间，日期格式统一 `YYYY-MM-DD`。
- 真实配置、构建产物、依赖目录、系统文件都应被 `.gitignore` 忽略。

## 文档

- [后端 README](backend/README.md)
- [前端 README](frontend/README.md)

# TODO
- [ ] AI总结还需要继续优化
- [ ] 用户需求：假设我有一笔钱，在当日支出，之后会AA，这种情况要如何处理呢
