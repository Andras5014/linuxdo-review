# ==================== 阶段1: 构建前端 ====================
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# 复制前端依赖文件
COPY frontend/package*.json ./

# 安装所有依赖（包括 devDependencies）
RUN npm install --include=dev

# 复制前端源码
COPY frontend/ ./

# 构建前端
RUN npm run build

# ==================== 阶段2: 构建后端 ====================
FROM golang:1.21-bookworm AS backend-builder

WORKDIR /app/backend

# 复制后端依赖文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源码
COPY backend/ ./

# 构建后端 (启用 CGO 以支持 SQLite)
RUN CGO_ENABLED=1 go build -o linuxdo-review .

# ==================== 阶段3: 运行时镜像 ====================
FROM debian:bookworm-slim

# 安装运行时依赖
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

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
