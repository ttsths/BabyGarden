.PHONY: all build run test clean docs docker-up docker-down docker-init

# 默认目标
all: build

# 构建
build:
	go build -o yuanzi-backend

# 运行开发环境
run:
	go run main.go

# 测试
test:
	go test ./... -v

# 格式化代码
fmt:
	go fmt ./...

# 静态检查
lint:
	golangci-lint run ./...

# 生成 Swagger 文档
docs:
	swag init -g main.go -o docs

# 清理
clean:
	rm -rf yuanzi-backend
	go clean -cache

# Docker 相关
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-stop:
	docker-compose stop

docker-restart:
	docker-compose restart

# 初始化数据库
docker-init:
	docker exec -i yuanzi-mysql mysql -uroot -prootpassword yuanzi < mysql/init/yuanzi.sql

# 构建并推送 Docker 镜像
docker-build:
	docker build -t yuanzi-backend:latest .

docker-push:
	docker push yuanzi-backend:latest

# 查看日志
logs:
	docker-compose logs -f

# 进入容器
shell:
	docker exec -it yuanzi-api sh
