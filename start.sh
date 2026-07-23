#!/bin/bash

# Telegram Bot Master 启动脚本

echo "🚀 启动 Telegram Bot Master..."

# 检查配置文件
if [ ! -f "config.env" ]; then
    echo "❌ 配置文件不存在，正在创建..."
    cp config.example.env config.env
    echo "✅ 已创建 config.env，请编辑后重新运行"
    exit 1
fi

# 检查数据目录
if [ ! -d "data" ]; then
    echo "📁 创建数据目录..."
    mkdir -p data
fi

# 编译并运行
echo "📦 编译中..."
go build -o bin/master cmd/master/main.go

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译完成"
echo "🎯 启动主控Bot..."
./bin/master
