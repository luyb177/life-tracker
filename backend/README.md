# life-tracker-backend
个人生活账本后端

## 命令行
生成文件
```bash
goctl api go -api life_tracker.api -dir . --style go_zero
```

生成swag
```bash
goctl api swagger -api life_tracker.api -dir ./docs -filename swagger
```