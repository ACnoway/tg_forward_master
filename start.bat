@echo off
REM Telegram Bot Master 启动脚本 (Windows)

echo 启动 Telegram Bot Master...

REM 检查配置文件
if not exist "config.env" (
    echo 配置文件不存在，正在创建...
    copy config.example.env config.env
    echo 已创建 config.env，请编辑后重新运行
    pause
    exit /b 1
)

REM 检查数据目录
if not exist "data" (
    echo 创建数据目录...
    mkdir data
)

REM 编译并运行
echo 编译中...
go build -o bin\master.exe cmd\master\main.go

if %ERRORLEVEL% NEQ 0 (
    echo 编译失败
    pause
    exit /b 1
)

echo 编译完成
echo 启动主控Bot...
bin\master.exe
