# 使用更新的 golang 基础镜像
FROM golang:1.21 AS builder

# 设置工作目录
WORKDIR /app

# 复制 Go 模块文件和主文件
COPY go.mod go.sum ./
COPY main.go .

# 下载依赖
RUN go mod download

# 使用 CGO_DISABLED=0 来构建 Go 二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o httpapi .

# 使用轻量级的 Alpine Linux 基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段将二进制文件复制到最终镜像中
COPY --from=builder /app/httpapi .

# 暴露应用运行的端口
EXPOSE 8080

# 运行二进制文件
CMD ["./httpapi"]
