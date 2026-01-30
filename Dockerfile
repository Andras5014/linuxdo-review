# ==================== 阶段1: 构建前端 ====================
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# 复制前端依赖文件
COPY frontend/package*.json ./

# 安装依赖
RUN npm install

# 复制前端源码
COPY frontend/ ./

# 构建前端
RUN npm run build

# ==================== 阶段2: 构建后端 ====================
FROM golang:1.21-alpine AS backend-builder

# 安装必要的构建工具
RUN apk add --no-cache gcc musl-dev

WORKDIR /app/backend

# 复制后端依赖文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源码
COPY backend/ ./

# 构建后端 (启用 CGO 以支持 SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o linuxdo-review .

# ==================== 阶段3: 运行时镜像 ====================
FROM alpine:3.19

# 安装时区数据和 ca 证书
RUN apk add --no-cache tzdata ca-certificates

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制后端可执行文件
COPY --from=backend-builder /app/backend/linuxdo-review ./

# 从构建阶段复制前端静态文件
COPY --from=frontend-builder /app/frontend/dist ./static

# 创建数据目录
RUN mkdir -p /app/data

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./linuxdo-review"]
