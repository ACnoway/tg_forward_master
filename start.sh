#!/bin/bash

# Telegram Bot 主控启动脚本 (Linux/Mac)

echo "🚀 启动 Telegram Bot 主控系统..."

# 检查配置文件
if [ ! -f "config.env" ]; then
    echo "❌ 错误: config.env 文件不存在"
    echo "请复制 config.example.env 为 config.env 并填写配置"
    exit 1
fi

# 检查数据目录
if [ ! -d "data" ]; then
    echo "📁 创建数据目录..."
    mkdir -p data
fi

# 编译程序
echo "🔨 编译程序..."
go build -o master cmd/master/main.go

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译完成"
echo "🚀 启动主控Bot..."
echo ""
echo "================================================"
echo "  Telegram Bot 主控系统 (Linux/Mac)"
echo "================================================"
echo ""
echo "首次使用步骤："
echo "1. 向Bot发送 /start"
echo "2. 使用 /setadmin 设置管理员"
echo "3. 使用 /admin_help 查看管理员命令"
echo ""
echo "按 Ctrl+C 停止运行"
echo "================================================"
echo ""

# 启动程序
./master
