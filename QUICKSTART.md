# 快速开始指南

## 一、环境准备

1. **安装 Go 1.21+**
   ```bash
   go version
   ```

2. **克隆项目**
   ```bash
   git clone <repository>
   cd tg_forward_master
   ```

3. **安装依赖**
   ```bash
   go mod download
   ```

## 二、配置

1. **复制配置文件**
   ```bash
   cp config.example.env config.env
   ```

2. **编辑配置文件**
   ```bash
   nano config.env
   ```

   **必填项：**
   - `MASTER_BOT_TOKEN` - 你的主控Bot Token（从 @BotFather 获取）
   - `ENCRYPTION_KEY` - 32字节加密密钥（随机生成）

   **可选项：**
   - `DATABASE_PATH` - 数据库路径（默认：./data/master.db）
   - `SERVER_PORT` - 服务器端口（默认：8080）
   - `DEFAULT_AI_ENDPOINT` - AI API地址（可后续在Bot中配置）
   - `DEFAULT_AI_KEY` - AI API密钥（可后续在Bot中配置）
   - `DEFAULT_AI_MODEL` - AI模型名称（可后续在Bot中配置）

   **生成加密密钥（32字节）：**
   ```bash
   # Linux/Mac
   openssl rand -hex 16

   # 或使用在线工具生成32位随机字符串
   ```

## 三、启动

### 方式1：使用启动脚本（推荐）
```bash
chmod +x start.sh
./start.sh
```

### 方式2：手动启动
```bash
# 编译
go build -o master.exe cmd/master/main.go

# 运行
./master.exe
```

## 四、首次使用

### 1. 设置管理员
在Telegram中向Bot发送以下命令：

```
/start
```

系统会提示你设置管理员，然后发送：

```
/setadmin
```

### 2. 查看管理员命令
```
/admin_help
```

### 3. 基础配置流程

#### (1) 创建套餐
```
/admin_add_plan
```
按提示输入：
- 套餐名称：如 "标准版"
- 价格：如 9.9
- 有效期：如 30（天）
- Bot数量：如 3

#### (2) 配置支付（如果需要用户付费）
```
/admin_payment_config
```
按提示输入易支付配置

#### (3) 配置AI（如果需要AI反垃圾）
```
/admin_ai_config
```
按提示输入AI配置

#### (4) 生成兑换码（用于免费赠送）
```
/admin_generate_code
```
选择套餐ID，输入生成数量

## 五、管理员命令列表

### 套餐管理
- `/admin_add_plan` - 创建新套餐
- `/admin_list_plans` - 查看所有套餐

### 兑换码管理
- `/admin_generate_code` - 批量生成兑换码
- `/admin_code_list` - 查看兑换码列表

### 支付管理
- `/admin_payment_config` - 配置易支付
- `/admin_payment_status` - 查看支付状态

### AI配置
- `/admin_ai_config` - 配置默认AI

### 系统管理
- `/admin_stats` - 系统统计
- `/admin_help` - 查看管理员帮助

## 六、用户使用流程

1. 用户发送 `/start` 查看状态
2. 用户使用 `/buy` 购买订阅（或 `/redeem` 使用兑换码）
3. 用户使用 `/createbot` 创建子Bot
4. 用户在子Bot中配置转发规则

## 七、常见问题

### 1. Bot无响应
- 检查 `MASTER_BOT_TOKEN` 是否正确
- 检查网络连接
- 查看日志输出

### 2. 数据库错误
- 确保 `data` 目录存在且有写入权限
- 检查 `DATABASE_PATH` 配置

### 3. 管理员命令无效
- 确认已使用 `/setadmin` 设置管理员
- 管理员命令必须以 `admin_` 开头（注意下划线）

### 4. 会话超时
- 多步骤命令（如创建套餐）有5分钟超时
- 超时后需要重新开始命令

## 八、项目结构

```
tg_forward_master/
├── cmd/
│   └── master/          # 主控Bot入口
├── internal/
│   ├── config/          # 配置管理
│   ├── database/        # 数据库层
│   ├── master/          # 主控Bot逻辑
│   │   ├── handlers/    # 命令处理器
│   │   └── service/     # 业务服务
│   ├── models/          # 数据模型
│   └── utils/           # 工具函数
├── migrations/          # 数据库迁移
├── data/               # 数据文件（自动创建）
├── config.env          # 配置文件（需自行创建）
├── config.example.env  # 配置示例
├── start.sh            # 启动脚本
└── README.md           # 项目说明
```

## 九、开发计划

- [x] 管理员功能
- [ ] 用户订阅购买
- [ ] 子Bot创建和管理
- [ ] 消息转发功能
- [ ] AI反垃圾检测
- [ ] 统计和监控

## 十、技术栈

- Go 1.21+
- SQLite 3
- telegram-bot-api
- AES-256 加密

## 十一、许可证

MIT License

## 十二、支持

- Issues: https://github.com/ACnoway/tg_forward_master/issues
- 文档: [AGENTS.md](AGENTS.md)
