#!/bin/bash

echo "========================================"
echo "   WeNote 一键启动脚本 (Linux/Mac)"
echo "========================================"
echo ""

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "[错误] Docker未运行，请先启动Docker"
    exit 1
fi

# 检查.env文件
if [ ! -f .env ]; then
    echo "[提示] 未找到.env文件，从.env.example复制..."
    cp .env.example .env
    echo "[完成] 已创建.env文件，请根据需要修改配置"
    echo ""
fi

echo "[1/3] 停止并清理旧容器..."
docker-compose down

echo ""
echo "[2/3] 构建并启动服务..."
docker-compose up -d --build

echo ""
echo "[3/3] 等待服务启动..."
sleep 10

echo ""
echo "========================================"
echo "   启动完成！"
echo "========================================"
echo ""
echo "前端地址: http://localhost"
echo "后端地址: http://localhost:8080"
echo "MySQL端口: 3307"
echo ""
echo "查看日志: docker-compose logs -f"
echo "停止服务: docker-compose down"
echo ""
