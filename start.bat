@echo off
chcp 65001 >nul
REM Telegram Bot 主控启动脚本 (Windows)

echo 🚀 启动 Telegram Bot 主控系统...
echo.

REM 检查配置文件
if not exist "config.env" (
    echo ❌ 错误: config.env 文件不存在
    echo 请复制 config.example.env 为 config.env 并填写配置
    pause
    exit /b 1
)

REM 检查数据目录
if not exist "data" (
    echo 📁 创建数据目录...
    mkdir data
)

REM 编译程序
echo 🔨 编译程序...
go build -o master.exe cmd/master/main.go

if errorlevel 1 (
    echo ❌ 编译失败
    pause
    exit /b 1
)

echo ✅ 编译完成
echo 🚀 启动主控Bot...
echo.
echo ================================================
echo   Telegram Bot 主控系统 (Windows)
echo ================================================
echo.
echo 首次使用步骤：
echo 1. 向Bot发送 /start
echo 2. 使用 /setadmin 设置管理员
echo 3. 使用 /admin_help 查看管理员命令
echo.
echo 按 Ctrl+C 停止运行
echo ================================================
echo.

REM 启动程序
master.exe
