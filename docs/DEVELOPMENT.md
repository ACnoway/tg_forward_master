# 开发指南

## 开发环境设置

### 必需工具
- Go 1.21+
- Git
- SQLite 3
- 文本编辑器/IDE（推荐VS Code）

### VS Code 推荐插件
- Go（官方）
- SQLite Viewer
- GitLens

### 环境变量
```bash
# 设置GOPATH（如果未设置）
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

## 项目结构说明

### cmd/
程序入口点
- `master/main.go` - 主控Bot入口
- `worker/main.go` - 子Bot入口（待创建）

### internal/
内部包，不对外暴露

#### internal/config/
配置管理
- 从环境变量和.env文件加载配置
- 配置验证

#### internal/database/
数据库操作
- `database.go` - 数据库连接管理
- `*_repository.go` - 各数据表的CRUD操作

#### internal/models/
数据模型定义
- 对应数据库表的结构体
- 业务逻辑方法

#### internal/master/
主控Bot业务逻辑
- `handlers/` - 命令处理器
- `service/` - 业务服务层（待创建）
- `payment/` - 支付集成（待创建）

#### internal/worker/
子Bot业务逻辑
- `handlers/` - 命令处理器（待创建）
- `forwarder/` - 消息转发（部分完成）
- `spam/` - AI反垃圾（待创建）

#### internal/utils/
工具函数
- 加密/解密
- 字符串处理
- 格式化工具

### migrations/
数据库迁移脚本
- 按序号命名：`001_init.sql`, `002_xxx.sql`

### docs/
文档目录
- `QUICKSTART.md` - 快速入门
- `STATUS.md` - 项目状态
- `DEVELOPMENT.md` - 本文件

## 开发工作流

### 1. 获取代码
```bash
git clone <repository>
cd tg_forward_master
```

### 2. 安装依赖
```bash
go mod download
```

### 3. 配置环境
```bash
cp config.example.env config.env
# 编辑 config.env，填入你的Bot Token
```

### 4. 运行开发环境
```bash
go run cmd/master/main.go
```

### 5. 测试更改
```bash
# 编译
go build -o bin/master cmd/master/main.go

# 运行
./bin/master
```

## 代码规范

### 命名规范
- 包名：小写，单个单词
- 文件名：小写，下划线分隔（user_repository.go）
- 类型名：大写开头（PascalCase）
- 函数名：
  - 导出函数：大写开头（PascalCase）
  - 内部函数：小写开头（camelCase）
- 常量：大写开头（PascalCase）或全大写（CONSTANT_NAME）

### 注释规范
```go
// Package handlers 提供主控Bot的命令处理功能
package handlers

// UserRepository 用户数据访问层
// 提供用户的CRUD操作
type UserRepository struct {
    db *DB
}

// Create 创建新用户
// 返回创建的用户ID，如果失败返回error
func (r *UserRepository) Create(user *User) error {
    // 实现...
}
```

### 错误处理
```go
// ✅ 好的做法
result, err := someFunction()
if err != nil {
    return fmt.Errorf("操作失败: %w", err)
}

// ❌ 避免这样
result, _ := someFunction() // 忽略错误
```

### 日志规范
```go
import "log"

// 信息日志
log.Printf("✅ 用户 %d 登录成功", userID)

// 错误日志
log.Printf("❌ 数据库连接失败: %v", err)

// 警告日志
log.Printf("⚠️ 用户 %d 尝试非法操作", userID)
```

## 添加新功能

### 1. 添加新的命令（主控Bot）

#### 步骤1：在handler.go中添加命令处理
```go
// handleCommand 中添加case
case "newcommand":
    h.handleNewCommand(message, user)
```

#### 步骤2：实现处理函数
```go
// handleNewCommand 处理 /newcommand 命令
func (h *MasterHandler) handleNewCommand(message *tgbotapi.Message, user *models.User) {
    // 实现逻辑
    h.sendMessage(message.Chat.ID, "命令执行成功")
}
```

#### 步骤3：更新help命令
在`handleHelp`中添加新命令说明

### 2. 添加新的数据表

#### 步骤1：在models.go中定义模型
```go
type NewModel struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### 步骤2：创建迁移脚本
```sql
-- migrations/002_add_new_table.sql
CREATE TABLE IF NOT EXISTS new_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### 步骤3：创建Repository
```go
// internal/database/new_repository.go
type NewRepository struct {
    db *DB
}

func NewNewRepository(db *DB) *NewRepository {
    return &NewRepository{db: db}
}

func (r *NewRepository) Create(model *models.NewModel) error {
    // 实现CRUD操作
}
```

### 3. 添加新的配置项

#### 步骤1：在config.go中添加字段
```go
type Config struct {
    // 现有字段...
    NewConfigItem string
}
```

#### 步骤2：在Load函数中加载
```go
config := &Config{
    // 现有配置...
    NewConfigItem: os.Getenv("NEW_CONFIG_ITEM"),
}
```

#### 步骤3：在config.example.env中添加
```env
NEW_CONFIG_ITEM=default_value
```

## 数据库操作

### 查看数据库
```bash
sqlite3 data/master.db

# SQLite命令
.tables                 # 查看所有表
.schema users           # 查看表结构
SELECT * FROM users;    # 查询数据
.quit                   # 退出
```

### 常用SQL
```sql
-- 查看用户
SELECT * FROM users;

-- 设置管理员
UPDATE users SET is_admin = 1 WHERE telegram_id = YOUR_ID;

-- 查看订阅
SELECT u.username, s.* 
FROM users u 
JOIN subscriptions s ON u.id = s.user_id;

-- 清空数据（慎用！）
DELETE FROM users;
```

### 备份数据库
```bash
# 备份
cp data/master.db data/master.db.backup

# 恢复
cp data/master.db.backup data/master.db
```

## 调试技巧

### 1. 启用Debug模式
```go
// cmd/master/main.go
bot.Debug = true  // 显示详细的API调用日志
```

### 2. 添加调试日志
```go
log.Printf("🐛 DEBUG: user=%+v", user)
```

### 3. 使用Delve调试器
```bash
# 安装
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试
dlv debug cmd/master/main.go
```

## 测试

### 单元测试
```go
// internal/utils/crypto_test.go
package utils

import "testing"

func TestEncryptDecrypt(t *testing.T) {
    key := "12345678901234567890123456789012"
    plaintext := "hello world"
    
    encrypted, err := Encrypt(plaintext, key)
    if err != nil {
        t.Fatalf("加密失败: %v", err)
    }
    
    decrypted, err := Decrypt(encrypted, key)
    if err != nil {
        t.Fatalf("解密失败: %v", err)
    }
    
    if decrypted != plaintext {
        t.Errorf("期望 %s, 得到 %s", plaintext, decrypted)
    }
}
```

### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/utils

# 带覆盖率
go test -cover ./...
```

## 常见问题

### Q: 如何添加新的管理员？
A: 
```bash
sqlite3 data/master.db
UPDATE users SET is_admin = 1 WHERE telegram_id = YOUR_TELEGRAM_ID;
.quit
```

### Q: 如何重置数据库？
A:
```bash
rm data/master.db
go run cmd/master/main.go  # 会自动重新创建
```

### Q: 编译错误：找不到包
A:
```bash
go mod tidy
go mod download
```

### Q: Bot不响应消息
A: 检查：
1. Bot Token是否正确
2. 网络连接是否正常
3. 查看控制台错误日志

## Git工作流

### 分支管理
```bash
# 创建功能分支
git checkout -b feature/new-feature

# 提交更改
git add .
git commit -m "feat: 添加新功能"

# 推送到远程
git push origin feature/new-feature

# 合并到main（通过PR）
```

### 提交信息规范
```
feat: 新功能
fix: 修复bug
docs: 文档更新
refactor: 重构
test: 测试相关
chore: 构建/工具相关
```

## 性能优化建议

### 1. 数据库查询
- 使用索引
- 避免N+1查询
- 使用事务批量操作

### 2. 并发处理
- 使用goroutine处理独立任务
- 使用channel通信
- 注意并发安全

### 3. 内存管理
- 及时释放资源
- 避免内存泄漏
- 使用对象池（如需要）

## 安全注意事项

### 1. 敏感信息
- ❌ 不要将Bot Token提交到Git
- ✅ 使用环境变量或.env文件
- ✅ .gitignore中排除config.env

### 2. 数据加密
- ✅ 使用提供的加密工具
- ✅ 密钥长度必须32字节
- ❌ 不要硬编码密钥

### 3. 权限检查
- ✅ 验证用户身份
- ✅ 检查管理员权限
- ✅ 验证订阅状态

## 参考资源

### Go语言
- [Go官方文档](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)

### Telegram Bot API
- [官方文档](https://core.telegram.org/bots/api)
- [Go Bot API库](https://github.com/go-telegram-bot-api/telegram-bot-api)

### SQLite
- [SQLite文档](https://www.sqlite.org/docs.html)
- [Go SQLite驱动](https://github.com/mattn/go-sqlite3)

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 获取帮助

- 查看文档：docs/目录
- 提交Issue：<repository>/issues
- 查看代码注释

---

Happy Coding! 🚀
