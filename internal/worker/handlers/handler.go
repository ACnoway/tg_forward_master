package handlers

import (
	"fmt"
	"log"

	"github.com/acnoway/tg_forward_master/internal/database"
	"github.com/acnoway/tg_forward_master/internal/models"
	"github.com/acnoway/tg_forward_master/internal/worker/forwarder"
	"github.com/acnoway/tg_forward_master/internal/worker/spam"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// WorkerHandler 子Bot处理器
type WorkerHandler struct {
	bot          *tgbotapi.BotAPI
	db           *database.DB
	botID        int64
	ownerID      int64
	customerRepo *database.CustomerRepository
	blockLogRepo *database.BlockLogRepository
	configRepo   *database.BotConfigRepository

	forwarder *forwarder.MessageForwarder
	aiChecker *spam.AISpamChecker
}

// NewWorkerHandler 创建子Bot处理器
func NewWorkerHandler(bot *tgbotapi.BotAPI, db *database.DB, botID, ownerID int64) *WorkerHandler {
	// 初始化仓库
	customerRepo := database.NewCustomerRepository(db)
	blockLogRepo := database.NewBlockLogRepository(db)
	configRepo := database.NewBotConfigRepository(db)

	// 初始化转发器
	forwarder := forwarder.NewMessageForwarder(bot, db, botID, ownerID)

	// 初始化AI检查器
	aiChecker := spam.NewAISpamChecker(db, botID)

	return &WorkerHandler{
		bot:          bot,
		db:           db,
		botID:        botID,
		ownerID:      ownerID,
		customerRepo: customerRepo,
		blockLogRepo: blockLogRepo,
		configRepo:   configRepo,
		forwarder:    forwarder,
		aiChecker:    aiChecker,
	}
}

// HandleMessage 处理消息
func (h *WorkerHandler) HandleMessage(message *tgbotapi.Message) {
	// 检查是否是主人的消息
	if message.From.ID == h.ownerID {
		h.handleOwnerMessage(message)
		return
	}

	// 检查是否是命令
	if message.IsCommand() {
		h.handleCommand(message)
		return
	}

	// 处理普通消息
	h.handleCustomerMessage(message)
}

// handleOwnerMessage 处理主人的消息
func (h *WorkerHandler) handleOwnerMessage(message *tgbotapi.Message) {
	// 处理主人的回复
	if message.ReplyToMessage != nil {
		h.handleOwnerReply(message)
		return
	}

	// 处理主人的命令
	h.sendMessage(message.Chat.ID, "主人可以使用以下命令：")
	h.sendMessage(message.Chat.ID, "/help - 帮助信息")
}

// handleOwnerReply 处理主人的回复消息
func (h *WorkerHandler) handleOwnerReply(message *tgbotapi.Message) {
	// 检查回复的目标消息
	if message.ReplyToMessage.ForwardFrom == nil {
		return
	}

	// 转发回复给客户
	if err := h.forwarder.ReplyToCustomer(message, message.ReplyToMessage.ForwardFrom.ID); err != nil {
		h.sendMessage(message.Chat.ID, "❌ 回复客户失败："+err.Error())
	}
}

// handleCommand 处理命令
func (h *WorkerHandler) handleCommand(message *tgbotapi.Message) {
	command := message.Command()

	switch command {
	case "start":
		h.handleStart(message)
	case "help":
		h.handleHelp(message)
	case "ai_config":
		h.handleAIConfig(message)
	case "ai_enable":
		h.handleAIEnable(message)
	case "ai_disable":
		h.handleAIDisable(message)
	case "ai_threshold":
		h.handleAIThreshold(message)
	case "block":
		h.handleBlock(message)
	case "blacklist":
		h.handleBlacklist(message)
	case "whitelist":
		h.handleWhitelist(message)
	case "trust":
		h.handleTrust(message)
	case "unblock":
		h.handleUnblock(message)
	case "untrust":
		h.handleUntrust(message)
	case "stats":
		h.handleStats(message)
	default:
		h.sendMessage(message.Chat.ID, "未知命令，使用 /help 查看可用命令")
	}
}

// handleCustomerMessage 处理客户消息
func (h *WorkerHandler) handleCustomerMessage(message *tgbotapi.Message) {
	// 检查是否在黑名单中
	customer, err := h.customerRepo.GetOrCreate(
		h.botID,
		message.From.ID,
		message.From.UserName,
		message.From.FirstName,
		message.From.LastName,
	)
	if err != nil {
		log.Printf("获取客户信息失败: %v", err)
		return
	}

	if customer.IsBlacklisted {
		h.blockCustomer(message, customer, "blacklist")
		return
	}

	// 检查是否在白名单中
	if customer.IsWhitelisted {
		// 直接转发
		if err := h.forwarder.ForwardToOwner(message); err != nil {
			log.Printf("转发消息失败: %v", err)
		}
		return
	}

	// 检查是否开启AI检测
	config, err := h.configRepo.GetByBotID(h.botID)
	if err != nil {
		log.Printf("获取配置失败: %v", err)
		return
	}

	if !config.AIEnabled {
		// 未开启AI检测，直接转发
		if err := h.forwarder.ForwardToOwner(message); err != nil {
			log.Printf("转发消息失败: %v", err)
		}
		return
	}

	// 检查AI是否通过
	score, err := h.aiChecker.CheckSpam(message.Text)
	if err != nil {
		log.Printf("AI检测失败: %v", err)
		// 检测失败，允许通过
		if err := h.forwarder.ForwardToOwner(message); err != nil {
			log.Printf("转发消息失败: %v", err)
		}
		return
	}

	// 检查是否是垃圾消息
	if score >= config.WhitelistThreshold {
		h.blockCustomer(message, customer, "ai_spam")
		return
	}

	// 通过AI检测，转发给主人
	if err := h.forwarder.ForwardToOwner(message); err != nil {
		log.Printf("转发消息失败: %v", err)
	}

	// 增加已验证计数
	if err := h.customerRepo.IncrementVerifiedCount(h.botID, message.From.ID, config.WhitelistThreshold); err != nil {
		log.Printf("更新验证计数失败: %v", err)
	}
}

// blockCustomer 阻止客户消息
func (h *WorkerHandler) blockCustomer(message *tgbotapi.Message, customer *models.Customer, reason string) {
	// 记录拦截日志
	logEntry := &models.BlockLog{
		BotID:          h.botID,
		TelegramID:     customer.TelegramID,
		Username:       customer.Username,
		FirstName:      customer.FirstName,
		LastName:       customer.LastName,
		BlockReason:    reason,
		MessageContent: message.Text,
		MessageType:    "text",
	}

	if err := h.blockLogRepo.Create(logEntry); err != nil {
		log.Printf("记录拦截日志失败: %v", err)
	}

	// 通知主人
	config, err := h.configRepo.GetByBotID(h.botID)
	if err != nil {
		log.Printf("获取配置失败: %v", err)
		config = &models.BotConfig{NotifyOnBlock: true}
	}

	if config.NotifyOnBlock {
		blockMsg := tgbotapi.NewMessage(h.ownerID, fmt.Sprintf(`📥 新拦截消息

👤 发件人: %s
原因: %s
内容: %s`,
			customer.GetDisplayName(),
			reason,
			message.Text,
		))
		h.bot.Send(blockMsg)
	}

	log.Printf("阻止客户消息: %d, 原因: %s", customer.TelegramID, reason)
}

// handleStart 处理 /start 命令
func (h *WorkerHandler) handleStart(message *tgbotapi.Message) {
	h.sendMessage(message.Chat.ID, "👋 欢迎使用消息转发Bot！")
	h.sendMessage(message.Chat.ID, "我已将您的消息转发给主人，请稍候等待回复。")
}

// handleHelp 处理 /help 命令
func (h *WorkerHandler) handleHelp(message *tgbotapi.Message) {
	h.sendMessage(message.Chat.ID, `📖 帮助文档

【常用命令】
/start - 开始使用
/help - 显示帮助

【AI管理】
/ai_config - 查看AI配置
/ai_enable - 开启AI检测
/ai_disable - 关闭AI检测
/ai_threshold <数字> - 设置白名单阈值

【黑名单管理】
/block - 拉黑用户（回复消息）
/blacklist - 查看黑名单
/unblock <ID> - 解除拉黑

【白名单管理】
/whitelist - 查看白名单
/trust - 信任用户（回复消息）
/untrust <ID> - 取消信任

【统计】
/stats - 查看统计信息`)
}

// sendMessage 发送消息
func (h *WorkerHandler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}

// HandleCallback 处理回调查询
func (h *WorkerHandler) HandleCallback(callback *tgbotapi.CallbackQuery) {
	// 处理回调
	h.bot.Send(tgbotapi.NewCallback(callback.ID, "功能开发中"))
}

// handleAIConfig 处理 /ai_config 命令
func (h *WorkerHandler) handleAIConfig(message *tgbotapi.Message) {
	config, err := h.configRepo.GetByBotID(h.botID)
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取AI配置失败："+err.Error())
		return
	}

	status := "✅ 已启用"
	if !config.AIEnabled {
		status = "❌ 已禁用"
	}

	mode := "默认AI配置"
	if config.UseCustomAI {
		mode = "自定义AI配置"
	}

	h.sendMessage(message.Chat.ID, fmt.Sprintf(`🤖 AI配置

状态：%s
模式：%s
白名单阈值：%d

使用以下命令管理AI：
/ai_enable - 开启AI检测
/ai_disable - 关闭AI检测
/ai_threshold <数字> - 设置白名单阈值`, status, mode, config.WhitelistThreshold))
}

// handleAIEnable 处理 /ai_enable 命令
func (h *WorkerHandler) handleAIEnable(message *tgbotapi.Message) {
	if err := h.aiChecker.EnableAI(); err != nil {
		h.sendMessage(message.Chat.ID, "❌ 启用AI失败："+err.Error())
		return
	}
	h.sendMessage(message.Chat.ID, "✅ AI检测已启用")
}

// handleAIDisable 处理 /ai_disable 命令
func (h *WorkerHandler) handleAIDisable(message *tgbotapi.Message) {
	if err := h.aiChecker.DisableAI(); err != nil {
		h.sendMessage(message.Chat.ID, "❌ 禁用AI失败："+err.Error())
		return
	}
	h.sendMessage(message.Chat.ID, "✅ AI检测已禁用")
}

// handleAIThreshold 处理 /ai_threshold 命令
func (h *WorkerHandler) handleAIThreshold(message *tgbotapi.Message) {
	h.sendMessage(message.Chat.ID, "请输入白名单阈值（0-100）：")
	// TODO: 实现阈值设置逻辑
}

// handleBlock 处理 /block 命令
func (h *WorkerHandler) handleBlock(message *tgbotapi.Message) {
	if message.ReplyToMessage == nil {
		h.sendMessage(message.Chat.ID, "请回复要拉黑的用户消息")
		return
	}

	if message.ReplyToMessage.ForwardFrom == nil {
		h.sendMessage(message.Chat.ID, "请回复转发的消息")
		return
	}

	if err := h.customerRepo.Blacklist(h.botID, message.ReplyToMessage.ForwardFrom.ID, "手动拉黑"); err != nil {
		h.sendMessage(message.Chat.ID, "❌ 拉黑失败："+err.Error())
		return
	}

	h.sendMessage(message.Chat.ID, "✅ 用户已拉黑")
}

// handleBlacklist 处理 /blacklist 命令
func (h *WorkerHandler) handleBlacklist(message *tgbotapi.Message) {
	customers, err := h.customerRepo.GetBlacklist(h.botID)
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取黑名单失败："+err.Error())
		return
	}

	if len(customers) == 0 {
		h.sendMessage(message.Chat.ID, "黑名单为空")
		return
	}

	var text string
	for _, customer := range customers {
		text += fmt.Sprintf("👤 %s\n", customer.GetDisplayName())
		text += fmt.Sprintf("   ID: %d\n", customer.TelegramID)
		if customer.BlacklistReason != "" {
			text += fmt.Sprintf("   原因: %s\n", customer.BlacklistReason)
		}
		text += "\n"
	}

	h.sendMessage(message.Chat.ID, "黑名单列表：\n\n"+text)
}

// handleWhitelist 处理 /whitelist 命令
func (h *WorkerHandler) handleWhitelist(message *tgbotapi.Message) {
	customers, err := h.customerRepo.GetWhitelist(h.botID)
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取白名单失败："+err.Error())
		return
	}

	if len(customers) == 0 {
		h.sendMessage(message.Chat.ID, "白名单为空")
		return
	}

	var text string
	for _, customer := range customers {
		text += fmt.Sprintf("👤 %s\n", customer.GetDisplayName())
		text += fmt.Sprintf("   ID: %d\n", customer.TelegramID)
		text += "\n"
	}

	h.sendMessage(message.Chat.ID, "白名单列表：\n\n"+text)
}

// handleTrust 处理 /trust 命令
func (h *WorkerHandler) handleTrust(message *tgbotapi.Message) {
	if message.ReplyToMessage == nil {
		h.sendMessage(message.Chat.ID, "请回复要信任的用户消息")
		return
	}

	if message.ReplyToMessage.ForwardFrom == nil {
		h.sendMessage(message.Chat.ID, "请回复转发的消息")
		return
	}

	if err := h.customerRepo.Whitelist(h.botID, message.ReplyToMessage.ForwardFrom.ID); err != nil {
		h.sendMessage(message.Chat.ID, "❌ 信任失败："+err.Error())
		return
	}

	h.sendMessage(message.Chat.ID, "✅ 用户已信任")
}

// handleUnblock 处理 /unblock 命令
func (h *WorkerHandler) handleUnblock(message *tgbotapi.Message) {
	// TODO: 实现解除拉黑逻辑
	h.sendMessage(message.Chat.ID, "功能开发中")
}

// handleUntrust 处理 /untrust 命令
func (h *WorkerHandler) handleUntrust(message *tgbotapi.Message) {
	// TODO: 实现取消信任逻辑
	h.sendMessage(message.Chat.ID, "功能开发中")
}

// handleStats 处理 /stats 命令
func (h *WorkerHandler) handleStats(message *tgbotapi.Message) {
	stats, err := h.configRepo.GetStats(h.botID)
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取统计失败："+err.Error())
		return
	}

	h.sendMessage(message.Chat.ID, fmt.Sprintf(`📊 统计信息

👥 客户统计：
总客户数：%d
白名单：%d
黑名单：%d

🛡️ 拦截统计：
总拦截：%d`,
		stats["total_customers"],
		stats["whitelisted_customers"],
		stats["blacklisted_customers"],
		stats["total_blocks"],
	))
}
