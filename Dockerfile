# 构建阶段
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 设置GOPROXY环境变量并下载依赖
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

# 复制源代码
COPY . .

# 生成swagger文档
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o chat_log_server

# 运行阶段
FROM alpine:latest

# 安装基本工具
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/chat_log_server .

# 暴露端口
EXPOSE 8090

# 运行应用
CMD ["./chat_log_server"]