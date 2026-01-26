#!/bin/bash
# WeNote 自动部署脚本
# 在服务器上执行: ./deploy.sh

set -e

echo "🚀 开始部署 WeNote..."

# 拉取最新代码
echo "📥 拉取最新代码..."
git pull origin main

# 停止旧容器
echo "🛑 停止旧容器..."
docker compose down

# 重新构建并启动
echo "🔨 构建并启动服务..."
docker compose up -d --build

# 清理旧镜像
echo "🧹 清理旧镜像..."
docker image prune -f

# 查看状态
echo "✅ 部署完成！服务状态："
docker compose ps

echo ""
echo "🌐 访问地址: http://$(curl -s ifconfig.me)"
