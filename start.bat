@echo off
chcp 65001 >nul
echo ========================================
echo    WeNote 一键启动脚本 (Windows)
echo ========================================
echo.

REM 检查Docker是否运行
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] Docker未运行，请先启动Docker Desktop
    pause
    exit /b 1
)

REM 检查.env文件
if not exist .env (
    echo [提示] 未找到.env文件，从.env.example复制...
    copy .env.example .env >nul
    echo [完成] 已创建.env文件，请根据需要修改配置
    echo.
)

echo [1/3] 停止并清理旧容器...
docker-compose down

echo.
echo [2/3] 构建并启动服务...
docker-compose up -d --build

echo.
echo [3/3] 等待服务启动...
timeout /t 10 /nobreak >nul

echo.
echo ========================================
echo    启动完成！
echo ========================================
echo.
echo 前端地址: http://localhost
echo 后端地址: http://localhost:8080
echo MySQL端口: 3307
echo.
echo 查看日志: docker-compose logs -f
echo 停止服务: docker-compose down
echo.
pause
