# 快速入门指南

## 第一步：安装依赖

### 确保已安装：
- Go 1.21 或更高版本
- Git

### 下载项目
```bash
git clone <repository>
cd tg_forward_master
```

## 第二步：配置Bot Token

### 1. 创建Telegram Bot
1. 在Telegram中找到 [@BotFather](https://t.me/BotFather)
2. 发送 `/newbot` 创建新Bot
3. 按提示设置Bot名称和用户名
4. 复制获得的Bot Token

### 2. 配置文件
```bash
# 复制配置示例
cp config.example.env config.env

# 编辑配置文件
nano config.env  # Linux/Mac
notepad config.env  # Windows
```

### 3. 填写必要配置
```env
# 将你的Bot Token填入
MASTER_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz

# 生成32位加密密钥（可以使用随机字符串）
ENCRYPTION_KEY=your_32_character_secret_key!!
```

## 第三步：安装依赖并运行

### 方法1：使用启动脚本（推荐）

**Linux/Mac:**
```bash
chmod +x start.sh
./start.sh
```

**Windows:**
```bash
start.bat
```

### 方法2：手动运行

```bash
# 安装依赖
go mod tidy

# 运行主控Bot
go run cmd/master/main.go
```

## 第四步：初始化系统

### 1. 向你的Bot发送消息
在Telegram中找到你的Bot，发送 `/start`

### 2. 设置管理员
首次运行时，系统会检测到没有管理员。

**手动设置管理员：**
```bash
# 停止Bot（Ctrl+C）

# 使用SQLite命令设置管理员
sqlite3 data/master.db
UPDATE users SET is_admin = 1 WHERE telegram_id = YOUR_TELEGRAM_ID;
.exit

# 重新启动Bot
./start.sh  # 或 start.bat
```

### 3. 测试管理员功能
向Bot发送 `/admin_help` 查看管理员命令

## 第五步：创建第一个子Bot

### 1. 创建子Bot Token
重复"第二步"的步骤1，创建另一个Bot作为转发Bot

### 2. 在主控Bot中创建子Bot
```
/createbot
# 按提示输入子Bot Token
```

## 常用命令

### 普通用户命令
```
/start - 开始使用
/help - 查看帮助
/mybots - 查看我的Bot列表
/createbot - 创建新Bot
```

### 管理员命令
```
/admin_help - 管理员帮助
/admin_stats - 系统统计
/admin_add_plan - 添加套餐
/admin_generate_code - 生成兑换码
```

### 子Bot命令（在子Bot中使用）
```
/help - 查看帮助
/ai_config - AI配置
/ai_enable - 开启AI检测
/ai_disable - 关闭AI检测
/block - 拉黑用户（回复消息）
/blacklist - 查看黑名单
/whitelist - 查看白名单
/stats - 查看统计
```

## 故障排除

### 问题1：Bot无响应
- 检查Bot Token是否正确
- 检查网络连接
- 查看控制台错误日志

### 问题2：数据库错误
```bash
# 删除数据库重新初始化
rm data/master.db
./start.sh
```

### 问题3：找不到命令
- 确保Go已正确安装：`go version`
- 确保在项目根目录运行命令

## 开发中的功能

当前版本实现了基础框架，以下功能正在开发：

- ✅ 用户管理
- ✅ 基础命令系统
- ⏳ 订阅管理
- ⏳ 支付集成
- ⏳ 子Bot创建和管理
- ⏳ 消息转发
- ⏳ AI反垃圾
- ⏳ 黑白名单

## 下一步

1. 查看 [AGENTS.md](AGENTS.md) 了解完整架构
2. 查看 [README.md](README.md) 了解详细功能
3. 关注项目更新获取新功能

## 需要帮助？

- 查看文档：[AGENTS.md](AGENTS.md)
- 提交Issue：<repository>/issues
- 查看日志：控制台输出

---

祝使用愉快！ 🚀
