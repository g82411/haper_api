# 使用 AWS Lambda 的 Go 运行时作为基础镜像
FROM golang:1.21 AS builder

# 设置工作目录
WORKDIR /app

# 将 Go 模块依赖复制到容器中并下载
COPY go.mod go.sum ./
RUN go mod download

# 将整个项目复制到容器中
COPY . .

# 编译应用程序，假设主程序位于 cmd/main.go
RUN GOOS=linux CGO_ENABLED=0 go build -o main ./cmd/api

# 使用相同的 Lambda Go 运行时镜像作为最终阶段的基础
FROM public.ecr.aws/lambda/go:1.2023.11.15.20

# 从构建阶段复制编译后的应用程序
COPY --from=builder /app/main /var/task/main

# 设置 Lambda 的处理程序执行文件
CMD ["main"]