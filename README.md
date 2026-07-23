# Telegram Bot 主控-子Bot系统

一个基于Go语言开发的Telegram Bot双层架构系统，支持订阅制商业化运营。

## 特性

### 🎯 核心功能
- **双层架构**: 主控Bot管理 + 多个子Bot实例
- **消息转发**: 客户消息自动转发给主人
- **智能反垃圾**: AI检测 + 智能白名单机制
- **黑名单管理**: 灵活的用户拦截功能
- **完整用户信息**: 显示昵称、用户名、ID

### 💰 商业化功能
- **订阅制管理**: 支持多种套餐
- **易支付集成**: 便捷的支付接口
- **兑换码系统**: 支持批量生成和管理
- **子Bot数量限制**: 根据订阅等级限制

### 🤖 AI反垃圾
- **智能检测**: 自动识别垃圾消息
- **自动白名单**: 3条消息验证后自动信任
- **成本优化**: 白名单用户不再消耗AI调用
- **灵活配置**: 可使用默认配置或自定义配置
- **一键开关**: 随时启用/禁用AI检测

### 📊 统计与监控
- **实时统计**: 消息转发、拦截统计
- **详细日志**: 完整的操作和拦截记录
- **系统监控**: 运行状态、用户数、Bot数

## 快速开始

### 1. 环境要求
- Go 1.21 或更高版本
- SQLite 3

### 2. 安装
```bash
git clone <repository>
cd tg_forward_master
go mod download
```

### 3. 配置
```bash
cp config.example.env config.env
# 编辑 config.env，填入你的Bot Token
```

### 4. 运行
```bash
go run cmd/master/main.go
```

### 5. 初始化
在Telegram中向主控Bot发送消息，按提示完成初始化。

## 使用指南

### 普通用户

#### 订阅管理
```
/start      - 开始使用
/buy        - 购买订阅
/redeem     - 使用兑换码
/myplan     - 查看订阅状态
```

#### 子Bot管理
```
/mybots     - 查看我的子Bot
/createbot  - 创建新子Bot
/managebot  - 管理子Bot
```

### 子Bot主人

#### AI配置
```
/ai_config     - 查看AI配置
/ai_enable     - 开启AI检测
/ai_disable    - 关闭AI检测
/ai_custom     - 使用自定义AI
/ai_threshold  - 设置白名单阈值
```

#### 黑白名单
```
/block      - 拉黑用户（回复消息）
/blacklist  - 查看黑名单
/unblock    - 解除拉黑
/whitelist  - 查看白名单
/trust      - 信任用户
```

#### 统计
```
/stats         - 今日统计
/stats_week    - 本周统计
/recent_blocks - 最近拦截
```

### 管理员

#### 支付管理
```
/admin payment_config - 配置支付
/admin payment_status - 支付状态
```

#### 套餐管理
```
/admin add_plan   - 创建套餐
/admin list_plans - 查看套餐
```

#### 兑换码
```
/admin generate_code - 生成兑换码
/admin code_list     - 兑换码列表
```

#### 系统管理
```
/admin stats - 系统统计
/admin users - 用户管理
```

## 架构说明

### 主控Bot
- 用户注册与订阅管理
- 支付集成
- 子Bot实例管理
- 系统配置

### 子Bot
- 消息转发
- AI反垃圾检测
- 黑白名单管理
- 统计日志

详细架构请查看 [AGENTS.md](AGENTS.md)

## 配置说明

### 最小化配置
系统采用"Bot交互优先"设计，配置文件仅需：

```env
MASTER_BOT_TOKEN=your_token
DATABASE_PATH=./data/master.db
ENCRYPTION_KEY=your_32_byte_key
SERVER_PORT=8080
```

其他配置均可通过Bot命令完成。

## AI反垃圾机制

### 工作流程
1. 新客户发送消息
2. AI检测是否为垃圾
3. 通过检测，转发给主人，计数+1
4. 累计3条（可配置）通过验证
5. 自动加入白名单
6. 后续消息直接转发，不再消耗AI调用

### 成本优化
- 每个客户最多消耗3次AI调用
- 验证后永久信任（除非手动移除）
- 可随时关闭AI检测
- 支持自定义API（不消耗主控配额）

## 数据库

使用SQLite 3，适合中小规模部署（< 5000用户）。

数据库文件默认位置: `./data/master.db`

### 核心表
- users: 用户信息
- subscriptions: 订阅记录
- worker_bots: 子Bot实例
- customers: 客户信息（昵称、用户名、ID）
- block_logs: 拦截日志

详细设计请查看 [AGENTS.md](AGENTS.md)

## 安全性

- ✅ Bot Token AES-256加密存储
- ✅ 默认AI配置对子Bot主人隐藏
- ✅ 用户自定义配置独立加密
- ✅ 完整的权限控制
- ✅ 操作审计日志

## 扩展性

### 当前架构支持
- 用户数: < 5000
- 子Bot数: < 500  
- 消息量: < 100万/天

### 扩展方案
需要更大规模时，可考虑：
- 迁移到PostgreSQL
- 分布式部署
- 消息队列
- Redis缓存

## 开发计划

- [x] 架构设计
- [ ] 核心功能实现
- [ ] 支付集成
- [ ] AI反垃圾
- [ ] 统计监控
- [ ] 文档完善

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License

## 支持

- 文档: [AGENTS.md](AGENTS.md)
- Issues: <repository>/issues
