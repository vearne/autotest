# 使用更新的 golang 基础镜像
FROM golang:1.21 AS builder

ARG PROJECT="github.com/vearne/autotest/example/grpc_api"

# Copy代码
ADD . $GOPATH/src/$PROJECT
# 设置工作目录
WORKDIR $GOPATH/src/$PROJECT

# 下载依赖
RUN go mod download

# 使用 CGO_DISABLED=0 来构建 Go 二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/grpcapi .

# 使用轻量级的 Alpine Linux 基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段将二进制文件复制到最终镜像中
COPY --from=builder /app/grpcapi .

# 暴露应用运行的端口
EXPOSE 50031

# 运行二进制文件
CMD ["./grpcapi"]
