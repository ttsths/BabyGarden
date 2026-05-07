FROM golang:1.24-alpine AS builder

# 设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /build

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o yuanzi-backend .

# 生产环境镜像
FROM alpine:latest

# 安装 CA 证书
RUN apk add --no-cache ca-certificates

# 设置工作目录
WORKDIR /app

# 复制配置文件
COPY --from=builder /build/config ./config

# 复制可执行文件
COPY --from=builder /build/yuanzi-backend ./yuanzi-backend

# 创建日志目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 8080

# 启动应用 (CF Containers 只读 ENTRYPOINT，不读 CMD)
ENTRYPOINT ["./yuanzi-backend"]
