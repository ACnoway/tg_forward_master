# 📚 文档导航

欢迎来到 Telegram Bot 主控-子Bot系统！

## 🚀 快速开始

如果你是第一次使用，请从这里开始：

1. **[快速入门指南](docs/QUICKSTART.md)** - 5分钟搭建运行环境
2. **[README.md](README.md)** - 项目介绍和功能概览

## 📖 核心文档

### 架构与设计
- **[AGENTS.md](AGENTS.md)** - 完整的系统架构文档
  - 数据库设计
  - 核心功能流程
  - API和命令列表
  - 安全性设计

### 开发指南
- **[开发指南](docs/DEVELOPMENT.md)** - 开发者必读
  - 开发环境设置
  - 代码规范
  - 如何添加新功能
  - 调试技巧

### 项目状态
- **[项目状态](docs/STATUS.md)** - 当前进度追踪
  - 已完成功能
  - 开发中功能
  - 待实现功能
  - 开发计划

## 🎯 按需查阅

### 我想...

#### 开始使用
→ [快速入门指南](docs/QUICKSTART.md)

#### 了解功能
→ [README.md](README.md) - 功能特性  
→ [AGENTS.md](AGENTS.md) - 详细设计

#### 参与开发
→ [开发指南](docs/DEVELOPMENT.md) - 开发规范  
→ [项目状态](docs/STATUS.md) - 查看待做任务

#### 部署上线
→ [快速入门](docs/QUICKSTART.md) - 基础部署  
→ [AGENTS.md](AGENTS.md) - 生产环境建议

#### 解决问题
→ [开发指南](docs/DEVELOPMENT.md) - 常见问题  
→ 查看GitHub Issues

## 📁 项目结构

```
tg_forward_master/
├── 📄 README.md              # 项目介绍
├── 📄 AGENTS.md              # 架构文档 ⭐
├── 📄 LICENSE                # 开源协议
├── 📄 DOCS.md                # 本文件
├── 📄 .gitignore             # Git忽略配置
├── 📄 go.mod                 # Go模块
├── 📄 config.example.env     # 配置示例
├── 📄 start.sh               # Linux/Mac启动
├── 📄 start.bat              # Windows启动
│
├── 📁 docs/                  # 文档目录
│   ├── QUICKSTART.md         # 快速入门 ⭐
│   ├── DEVELOPMENT.md        # 开发指南 ⭐
│   └── STATUS.md             # 项目状态 ⭐
│
├── 📁 cmd/                   # 程序入口
│   ├── master/main.go        # 主控Bot
│   └── worker/main.go        # 子Bot（待创建）
│
├── 📁 internal/              # 内部代码
│   ├── config/               # 配置管理
│   ├── database/             # 数据库层
│   ├── models/               # 数据模型
│   ├── utils/                # 工具函数
│   ├── master/               # 主控Bot逻辑
│   └── worker/               # 子Bot逻辑
│
└── 📁 migrations/            # 数据库迁移
    └── 001_init.sql          # 初始化脚本
```

## 🔥 核心概念

### 双层架构
- **主控Bot**: 管理平台，处理订阅、创建子Bot
- **子Bot**: 消息转发实例，每个用户可拥有多个

### 消息转发流程
```
客户 → 子Bot → [AI检测] → [黑名单过滤] → 主人
主人 → 子Bot → 客户
```

### AI智能白名单
- 新用户前3条消息经AI检测
- 通过后自动加入白名单
- 白名单用户消息不再消耗AI调用
- 大幅降低运营成本

### 商业化功能
- 订阅制管理
- 易支付集成
- 兑换码系统
- 子Bot数量限制

## 🎓 学习路径

### 初学者
1. 阅读 [README.md](README.md) 了解项目
2. 跟随 [快速入门](docs/QUICKSTART.md) 运行起来
3. 浏览 [AGENTS.md](AGENTS.md) 理解架构

### 开发者
1. 完成初学者路径
2. 精读 [开发指南](docs/DEVELOPMENT.md)
3. 查看 [项目状态](docs/STATUS.md) 选择任务
4. 开始贡献代码

### 运营者
1. 完成初学者路径
2. 了解 [AGENTS.md](AGENTS.md) 中的商业化功能
3. 配置支付接口
4. 设置套餐和兑换码

## 💡 最佳实践

### 开发环境
```bash
# 1. 克隆项目
git clone <repository>
cd tg_forward_master

# 2. 配置环境
cp config.example.env config.env
# 编辑 config.env

# 3. 运行
go run cmd/master/main.go
```

### 生产环境
```bash
# 1. 编译
go build -o master cmd/master/main.go

# 2. 使用systemd管理
# 参考 AGENTS.md 部署章节

# 3. 定期备份数据库
cp data/master.db backups/master_$(date +%Y%m%d).db
```

## 🆘 获取帮助

### 文档没有解答？
1. 查看 [常见问题](docs/DEVELOPMENT.md#常见问题)
2. 搜索 GitHub Issues
3. 提交新 Issue

### 发现Bug？
1. 确认是否已知问题
2. 准备复现步骤
3. 提交详细的 Bug Report

### 想要新功能？
1. 查看 [项目状态](docs/STATUS.md) 确认未在规划中
2. 提交 Feature Request
3. 或直接贡献代码（推荐！）

## 📊 项目统计

### 当前状态
- ✅ 核心框架完成
- ✅ 数据库设计完成
- ✅ 主控Bot基础功能完成
- ⏳ 子Bot功能开发中
- ⏳ 商业化功能开发中

### 代码量（当前）
- Go代码：~2000+ 行
- SQL：~200 行
- 文档：~5000+ 行

### 完成度
- 整体进度：~30%
- 核心框架：80%
- 主控Bot：40%
- 子Bot：10%
- 商业化：5%

## 🗺️ 路线图

### 短期目标（1-2周）
- [ ] 完成子Bot核心转发功能
- [ ] 实现黑名单管理
- [ ] 完成订阅系统基础功能

### 中期目标（3-4周）
- [ ] AI反垃圾检测
- [ ] 支付集成
- [ ] 兑换码系统

### 长期目标（5-8周）
- [ ] 完整统计系统
- [ ] 性能优化
- [ ] 生产环境部署
- [ ] 用户文档完善

## 📝 更新日志

### v0.1.0 (2026-07-23)
- ✅ 项目初始化
- ✅ 核心架构设计
- ✅ 数据库设计
- ✅ 主控Bot基础功能
- ✅ 完整文档系统

## 🤝 贡献

欢迎贡献！请查看：
- [开发指南](docs/DEVELOPMENT.md) - 了解代码规范
- [项目状态](docs/STATUS.md) - 选择任务
- GitHub Issues - 报告问题或领取任务

## 📞 联系方式

- 项目地址：<repository>
- Issue跟踪：<repository>/issues
- 文档反馈：提交Issue标注 `documentation`

---

**文档版本**: v1.0  
**最后更新**: 2026-07-23  
**维护者**: Project Team

开始你的旅程 → [快速入门](docs/QUICKSTART.md) 🚀
