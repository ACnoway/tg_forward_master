# Telegram Bot Master-Worker System

## 项目概述

这是一个基于Go语言开发的Telegram Bot双层架构系统，包含主控Bot（Master Bot）和子Bot（Worker Bot）。主控Bot负责用户管理、订阅管理和子Bot创建，子Bot负责消息转发和AI反垃圾检测。

## 核心功能

### 主控Bot（Master Bot）
- 用户注册与订阅管理
- 易支付集成（订阅制）
- 兑换码生成与管理
- 子Bot实例创建与管理
- 套餐配置
- 系统统计与监控

### 子Bot（Worker Bot）
- 消息转发（客户 → 主人）
- 回复转发（主人 → 客户）
- 黑名单管理
- AI反垃圾检测（智能白名单机制）
- 统计与日志

## 技术架构

### 技术栈
- **语言**: Go 1.21+
- **数据库**: SQLite 3
- **Bot框架**: telegram-bot-api
- **加密**: AES-256
- **AI接口**: OpenAI兼容API

### 项目结构
```
tg_forward_master/
├── cmd/
│   ├── master/          # 主控Bot入口
│   │   └── main.go
│   └── worker/          # 子Bot入口（测试用）
│       └── main.go
├── internal/
│   ├── master/          # 主控Bot业务逻辑
│   │   ├── handlers/    # 命令处理器
│   │   ├── payment/     # 支付集成
│   │   └── service/     # 业务服务
│   ├── worker/          # 子Bot业务逻辑
│   │   ├── handlers/    # 命令处理器
│   │   ├── forwarder/   # 消息转发
│   │   └── spam/        # AI反垃圾
│   ├── models/          # 数据模型
│   ├── database/        # 数据库操作
│   ├── config/          # 配置管理
│   └── utils/           # 工具函数
├── migrations/          # 数据库迁移
├── docs/               # 文档
├── config.example.env  # 配置示例
├── go.mod
├── go.sum
├── README.md
└── AGENTS.md          # 本文件
```

## 数据库设计

### 核心表

#### users - 用户表
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id INTEGER UNIQUE NOT NULL,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    is_admin BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### subscriptions - 订阅表
```sql
CREATE TABLE subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    plan_id INTEGER NOT NULL,
    status TEXT DEFAULT 'active',  -- active, expired, banned
    expires_at DATETIME NOT NULL,
    max_bots INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

#### plans - 套餐表
```sql
CREATE TABLE plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    duration_days INTEGER NOT NULL,
    max_bots INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### redeem_codes - 兑换码表
```sql
CREATE TABLE redeem_codes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT UNIQUE NOT NULL,
    plan_id INTEGER NOT NULL,
    status TEXT DEFAULT 'unused',  -- unused, used, expired
    used_by INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    used_at DATETIME,
    FOREIGN KEY (plan_id) REFERENCES plans(id),
    FOREIGN KEY (used_by) REFERENCES users(id)
);
```

#### worker_bots - 子Bot表
```sql
CREATE TABLE worker_bots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    bot_token TEXT NOT NULL,
    bot_username TEXT NOT NULL,
    bot_nickname TEXT,
    owner_telegram_id INTEGER NOT NULL,
    status TEXT DEFAULT 'stopped',  -- running, stopped
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

#### bot_configs - 子Bot配置表
```sql
CREATE TABLE bot_configs (
    bot_id INTEGER PRIMARY KEY,
    owner_telegram_id INTEGER NOT NULL,
    
    -- AI配置
    ai_enabled BOOLEAN DEFAULT 1,
    use_custom_ai BOOLEAN DEFAULT 0,
    custom_ai_endpoint TEXT,
    custom_ai_key TEXT,
    custom_ai_model TEXT,
    whitelist_threshold INTEGER DEFAULT 3,
    
    -- 通知配置
    notify_on_block BOOLEAN DEFAULT 1,
    notify_on_ai_block BOOLEAN DEFAULT 0,
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bot_id) REFERENCES worker_bots(id)
);
```

#### customers - 客户表
```sql
CREATE TABLE customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bot_id INTEGER NOT NULL,
    telegram_id INTEGER NOT NULL,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    
    verified_count INTEGER DEFAULT 0,
    total_messages INTEGER DEFAULT 0,
    is_whitelisted BOOLEAN DEFAULT 0,
    is_blacklisted BOOLEAN DEFAULT 0,
    
    first_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_message_at DATETIME,
    whitelisted_at DATETIME,
    blacklisted_at DATETIME,
    
    blacklist_reason TEXT,
    is_manual_trust BOOLEAN DEFAULT 0,
    
    UNIQUE(bot_id, telegram_id),
    FOREIGN KEY (bot_id) REFERENCES worker_bots(id)
);
```

#### block_logs - 拦截日志表
```sql
CREATE TABLE block_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bot_id INTEGER NOT NULL,
    telegram_id INTEGER NOT NULL,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    
    block_reason TEXT,  -- 'blacklist' or 'ai_spam'
    ai_confidence REAL,
    message_content TEXT,
    message_type TEXT,
    
    blocked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_false_positive BOOLEAN DEFAULT 0,
    FOREIGN KEY (bot_id) REFERENCES worker_bots(id)
);
```

#### system_configs - 系统配置表
```sql
CREATE TABLE system_configs (
    key TEXT PRIMARY KEY,
    value TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### orders - 订单表
```sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    plan_id INTEGER NOT NULL,
    order_no TEXT UNIQUE NOT NULL,
    amount REAL NOT NULL,
    status TEXT DEFAULT 'pending',  -- pending, paid, failed
    payment_url TEXT,
    paid_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (plan_id) REFERENCES plans(id)
);
```

## 核心功能流程

### 1. 用户订阅流程
```
用户发送 /start
  ↓
检查是否已注册 → 否 → 创建用户记录
  ↓ 是
检查订阅状态
  ↓
显示主菜单（购买订阅/兑换码/我的Bot）
```

### 2. 创建子Bot流程
```
用户发送 /createbot
  ↓
检查订阅状态和Bot数量限制
  ↓
要求输入Bot Token
  ↓
验证Token有效性
  ↓
要求输入Bot昵称
  ↓
要求输入主人Telegram ID
  ↓
创建Bot记录和配置
  ↓
询问是否立即启动
```

### 3. 消息转发流程（子Bot）
```
收到客户消息
  ↓
检查黑名单 → 是 → 拦截 + 记录日志
  ↓ 否
检查AI是否开启 → 否 → 直接转发
  ↓ 是
检查白名单 → 是 → 直接转发
  ↓ 否
调用AI检测
  ↓
是垃圾？ → 是 → 拦截 + 记录日志 + 可选通知主人
  ↓ 否
转发给主人 + verified_count++
  ↓
达到阈值？ → 是 → 加入白名单
```

### 4. AI智能白名单机制
```
每个客户独立计数
  ↓
通过AI验证的消息 verified_count++
  ↓
当 verified_count >= threshold（默认3）
  ↓
自动加入白名单
  ↓
后续消息不再调用AI，直接转发
  ↓
节省API调用成本
```

## 配置管理

### 最小化配置文件
系统采用"Bot交互优先"的设计理念，配置文件仅保留以下必要项：

```env
# config.env
MASTER_BOT_TOKEN=your_master_bot_token
DATABASE_PATH=./data/master.db
ENCRYPTION_KEY=your_32_byte_encryption_key
SERVER_PORT=8080
```

### 其他配置通过Bot命令管理
- 支付配置: `/admin payment_config`
- 套餐管理: `/admin add_plan`, `/admin list_plans`
- AI配置: `/admin ai_config`
- 子Bot配置: 在子Bot中使用 `/ai_config`, `/ai_custom` 等

## 命令列表

### 主控Bot - 用户命令
| 命令 | 说明 |
|------|------|
| /start | 开始使用/查看状态 |
| /buy | 购买订阅套餐 |
| /redeem | 使用兑换码 |
| /myplan | 查看我的订阅信息 |
| /mybots | 查看我的子Bot列表 |
| /createbot | 创建新的子Bot |
| /managebot | 管理已有的子Bot |
| /help | 显示帮助信息 |

### 主控Bot - 管理员命令
| 命令 | 说明 |
|------|------|
| /admin payment_config | 配置支付接口 |
| /admin payment_status | 查看支付状态 |
| /admin add_plan | 创建新套餐 |
| /admin list_plans | 查看/编辑套餐 |
| /admin generate_code | 生成兑换码 |
| /admin code_list | 查看兑换码列表 |
| /admin ai_config | 配置默认AI |
| /admin stats | 系统统计 |
| /admin users | 用户列表 |
| /admin help | 管理员帮助 |

### 子Bot - 主人命令
| 命令 | 说明 |
|------|------|
| /ai_config | AI配置总览 |
| /ai_enable | 开启AI检测 |
| /ai_disable | 关闭AI检测 |
| /ai_custom | 使用自定义AI配置 |
| /ai_default | 切换回默认配置 |
| /ai_threshold <数字> | 设置白名单阈值 |
| /block | 拉黑用户（回复消息） |
| /blacklist | 查看黑名单列表 |
| /unblock <ID> | 解除拉黑 |
| /whitelist | 查看白名单列表 |
| /trust | 信任用户（回复消息） |
| /untrust <ID> | 取消信任 |
| /stats | 查看统计信息 |
| /help | 显示帮助信息 |

## 安全性设计

### 1. 敏感信息保护
- Bot Token使用AES-256加密存储
- 主控默认AI API Key通过代理层调用，不暴露给子Bot主人
- 用户自定义AI配置独立加密存储
- 数据库文件权限控制

### 2. 权限控制
- 只有主人可以管理自己的子Bot
- 管理员特殊权限（查看所有Bot、封禁用户）
- 命令执行前验证用户身份

### 3. 防滥用
- 创建子Bot频率限制
- AI调用频率限制
- 兑换码使用次数限制
- 智能白名单减少AI调用成本

## 部署说明

### 1. 环境准备
```bash
# 安装Go 1.21+
# 克隆项目
git clone <repository>
cd tg_forward_master

# 安装依赖
go mod download
```

### 2. 配置
```bash
# 复制配置示例
cp config.example.env config.env

# 编辑配置文件
nano config.env
```

### 3. 初始化数据库
```bash
# 数据库会在首次运行时自动创建
mkdir -p data
```

### 4. 启动主控Bot
```bash
# 开发环境
go run cmd/master/main.go

# 生产环境
go build -o master cmd/master/main.go
./master
```

### 5. 配置管理员
```
首次启动后，在Telegram中向主控Bot发送消息
系统会提示设置管理员
```

## 扩展性考虑

### 当前架构适用规模
- 用户数: < 5000
- 子Bot数: < 500
- 消息量: < 100万/天

### 扩展方案
如需支持更大规模：
1. 迁移到PostgreSQL
2. 子Bot分布式部署
3. 消息队列处理AI检测
4. Redis缓存

## 开发路线图

### Phase 1: 核心功能 ✅
- [x] 项目架构设计
- [ ] 数据库设计与迁移
- [ ] 主控Bot基础框架
- [ ] 子Bot消息转发

### Phase 2: 商业化
- [ ] 易支付集成
- [ ] 订阅管理
- [ ] 兑换码系统

### Phase 3: 高级功能
- [ ] 黑名单管理
- [ ] AI反垃圾（智能白名单）
- [ ] 统计与日志

### Phase 4: 优化上线
- [ ] 安全加固
- [ ] 性能优化
- [ ] 监控告警
- [ ] 文档完善

## 贡献指南

### 代码规范
- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 添加必要的注释
- 编写单元测试

### 提交规范
- feat: 新功能
- fix: 修复bug
- docs: 文档更新
- refactor: 重构
- test: 测试相关

## 许可证

MIT License

## 联系方式

- Issues: <repository>/issues
- Email: your-email@example.com

---

最后更新: 2026-07-23
