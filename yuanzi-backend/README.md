# Yuanzi Backend - Go + Gin

小园子母婴记录应用后端服务

## 技术栈

- **框架**: Gin (v1.9.1)
- **数据库**: MySQL 8.0 + GORM
- **缓存**: Redis + go-redis/v9
- **日志**: Zap
- **配置**: Viper
- **认证**: JWT
- **文档**: Swagger

## 项目结构

```
yuanzi-backend/
├── handler/          # 控制器/HTTP 处理
├── biz/              # 业务逻辑
├── model/            # 数据模型
├── middleware/       # 中间件
├── router/           # 路由配置
├── config/           # 配置管理
├── logger/           # 日志
├── pkg/              # 公共工具包
│   └── gredis/       # Redis 封装
├── docs/             # Swagger 文档
├── mysql/            # MySQL 配置
│   └── init/         # 初始化脚本
├── build/            # 构建产物
├── config.yaml       # 配置文件
├── go.mod            # 项目依赖
├── Dockerfile        # Docker 镜像
└── docker-compose.yml # Docker Compose 配置
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 7+

### 本地开发

```bash
# 1. 安装依赖
go mod download

# 2. 配置数据库
# 编辑 config/config.yaml

# 3. 运行
go run main.go

# 4. 访问 Swagger 文档
# http://localhost:8080/swagger/index.html
```

### Docker 部署

```bash
# 1. 构建并运行
docker-compose up -d

# 2. 运行数据库初始化
docker exec -i yuanzi-mysql mysql -uyuanzi -pyuanzi123 yuanzi < mysql/init/yuanzi.sql

# 3. 查看日志
docker-compose logs -f app
```

## 开发命令

```bash
# 下载依赖
go mod download

# 生成 Swagger 文档
swag init -g main.go -o docs

# 运行测试
go test ./...

# 构建二进制
go build -o yuanzi-backend

# 格式化代码
gofmt -w .
```

## 配置说明

参考 `config/config.yaml`：

```yaml
server:
  run_mode: debug  # debug/release
  http_port: 8080
  read_timeout: 60
  write_timeout: 60

database:
  type: mysql
  host: localhost
  port: 3306
  name: yuanzi
  user: root
  password: your_password

redis:
  host: localhost
  port: 6379
```

## RESTful API 文档

启动服务后访问: `http://localhost:8080/swagger/index.html`

## 部署

参见 `Dockerfile` 和 `docker-compose.yml`

## 开发日志

- 2026-03-08: 初始化项目，搭建基础框架

## 联系方式

技术问题请联系 BE 团队。
