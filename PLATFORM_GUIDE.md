# 平台使用指南

## 快速参考

| 操作系统 | 启动脚本 | 编译命令 | 可执行文件 |
|---------|---------|---------|-----------|
| Linux | `./start.sh` | `go build -o master cmd/master/main.go` | `./master` |
| Mac | `./start.sh` | `go build -o master cmd/master/main.go` | `./master` |
| Windows | `start.bat` | `go build -o master.exe cmd/master/main.go` | `master.exe` |

## Linux/Mac 使用

### 1. 首次设置
```bash
# 赋予执行权限
chmod +x start.sh

# 复制配置文件
cp config.example.env config.env

# 编辑配置
nano config.env
```

### 2. 启动
```bash
./start.sh
```

### 3. 手动编译（可选）
```bash
go build -o master cmd/master/main.go
./master
```

### 4. 后台运行（可选）
```bash
# 使用 nohup
nohup ./master > output.log 2>&1 &

# 或使用 screen
screen -S tgbot
./master
# Ctrl+A, D 分离会话

# 或使用 tmux
tmux new -s tgbot
./master
# Ctrl+B, D 分离会话
```

### 5. 查看日志
```bash
tail -f output.log
```

## Windows 使用

### 1. 首次设置
```cmd
# 复制配置文件
copy config.example.env config.env

# 编辑配置（使用记事本或其他编辑器）
notepad config.env
```

### 2. 启动
```cmd
start.bat
```

或直接双击 `start.bat`

### 3. 手动编译（可选）
```cmd
go build -o master.exe cmd/master/main.go
master.exe
```

### 4. 后台运行（可选）
```cmd
# 创建一个VBS脚本 run_hidden.vbs
Set WshShell = CreateObject("WScript.Shell")
WshShell.Run "master.exe", 0
Set WshShell = Nothing

# 然后双击 run_hidden.vbs 启动
```

### 5. 使用Windows服务（高级）
可以使用 NSSM (Non-Sucking Service Manager) 将程序注册为Windows服务：
```cmd
# 下载 NSSM: https://nssm.cc/download
nssm install TelegramBot "C:\path\to\master.exe"
nssm start TelegramBot
```

## 配置文件说明

### 必填项
```env
# 从 @BotFather 获取
MASTER_BOT_TOKEN=your_bot_token_here

# 32字节加密密钥（生成方法见下文）
ENCRYPTION_KEY=your_32_byte_encryption_key_here!!
```

### 可选项
```env
# 数据库路径（默认：./data/master.db）
DATABASE_PATH=./data/master.db

# 服务器端口（默认：8080）
SERVER_PORT=8080

# 以下配置可以在Bot中通过命令设置，无需在此配置
# DEFAULT_AI_ENDPOINT=https://api.openai.com/v1/chat/completions
# DEFAULT_AI_KEY=
# DEFAULT_AI_MODEL=gpt-3.5-turbo
```

### 生成加密密钥

**Linux/Mac:**
```bash
openssl rand -hex 16
# 或
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1
```

**Windows (PowerShell):**
```powershell
-join ((65..90) + (97..122) + (48..57) | Get-Random -Count 32 | ForEach-Object {[char]$_})
```

**在线生成:**
- https://www.random.org/strings/
- 设置：长度32，数字+字母

## 开发模式

### 直接运行（不编译）
```bash
# Linux/Mac/Windows 通用
go run cmd/master/main.go
```

### 热重载（需要安装 air）
```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 运行
air
```

## 生产环境部署

### 使用 systemd (Linux)

1. 创建服务文件 `/etc/systemd/system/tgbot.service`:
```ini
[Unit]
Description=Telegram Bot Master
After=network.target

[Service]
Type=simple
User=your_user
WorkingDirectory=/path/to/tg_forward_master
ExecStart=/path/to/tg_forward_master/master
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

2. 启动服务:
```bash
sudo systemctl daemon-reload
sudo systemctl enable tgbot
sudo systemctl start tgbot

# 查看状态
sudo systemctl status tgbot

# 查看日志
sudo journalctl -u tgbot -f
```

### 使用 Docker (跨平台)

1. 创建 `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o master cmd/master/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/master .
COPY config.env .
CMD ["./master"]
```

2. 构建和运行:
```bash
docker build -t tgbot .
docker run -d --name tgbot -v $(pwd)/data:/root/data tgbot
```

## 常见问题

### Linux/Mac: Permission denied
```bash
chmod +x start.sh
chmod +x master
```

### Windows: 无法识别命令
确保 Go 已添加到 PATH 环境变量

### 所有平台: 端口被占用
修改 `config.env` 中的 `SERVER_PORT`

### 所有平台: 数据库权限错误
```bash
# Linux/Mac
chmod 755 data
chmod 644 data/master.db

# Windows
检查文件夹权限设置
```

## 更新程序

### 拉取最新代码
```bash
git pull origin master
```

### 重新编译
```bash
# Linux/Mac
./start.sh

# Windows
start.bat
```

### 查看更新日志
```bash
git log --oneline -10
```

## 性能优化

### 编译优化
```bash
# Linux/Mac
go build -ldflags="-s -w" -o master cmd/master/main.go

# Windows
go build -ldflags="-s -w" -o master.exe cmd/master/main.go
```
- `-s`: 去除符号表
- `-w`: 去除调试信息
- 可减小约 30% 文件大小

### 运行时优化
```bash
# 设置 Go 运行时参数
GOMAXPROCS=4 ./master
```

## 获取帮助

- 文档: [README.md](README.md)
- 快速开始: [QUICKSTART.md](QUICKSTART.md)
- 测试清单: [TEST_CHECKLIST.md](TEST_CHECKLIST.md)
- 开发文档: [AGENTS.md](AGENTS.md)
- Issues: https://github.com/ACnoway/tg_forward_master/issues
