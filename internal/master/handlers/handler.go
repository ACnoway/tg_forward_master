package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/acnoway/tg_forward_master/internal/config"
	"github.com/acnoway/tg_forward_master/internal/database"
	"github.com/acnoway/tg_forward_master/internal/master/service"
	"github.com/acnoway/tg_forward_master/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MasterHandler 主控Bot处理器
type MasterHandler struct {
	bot    *tgbotapi.BotAPI
	db     *database.DB
	config *config.Config

	userRepo         *database.UserRepository
	planRepo         *database.PlanRepository
	subscriptionRepo *database.SubscriptionRepository
	redeemCodeRepo   *database.RedeemCodeRepository
	workerBotRepo    *database.WorkerBotRepository
	configRepo       *database.SystemConfigRepository

	adminService   *service.AdminService
	sessionManager *SessionManager
}

// NewMasterHandler 创建主控Bot处理器
func NewMasterHandler(bot *tgbotapi.BotAPI, db *database.DB, cfg *config.Config) *MasterHandler {
	// 初始化仓库
	userRepo := database.NewUserRepository(db)
	planRepo := database.NewPlanRepository(db)
	subscriptionRepo := database.NewSubscriptionRepository(db)
	redeemCodeRepo := database.NewRedeemCodeRepository(db)
	workerBotRepo := database.NewWorkerBotRepository(db)
	configRepo := database.NewSystemConfigRepository(db)

	// 初始化服务
	adminService := service.NewAdminService(
		userRepo,
		planRepo,
		redeemCodeRepo,
		subscriptionRepo,
		workerBotRepo,
		configRepo,
	)

	return &MasterHandler{
		bot:              bot,
		db:               db,
		config:           cfg,
		userRepo:         userRepo,
		planRepo:         planRepo,
		subscriptionRepo: subscriptionRepo,
		redeemCodeRepo:   redeemCodeRepo,
		workerBotRepo:    workerBotRepo,
		configRepo:       configRepo,
		adminService:     adminService,
		sessionManager:   NewSessionManager(),
	}
}

// HandleMessage 处理消息
func (h *MasterHandler) HandleMessage(message *tgbotapi.Message) {
	// 获取或创建用户
	user, err := h.userRepo.GetOrCreate(
		message.From.ID,
		message.From.UserName,
		message.From.FirstName,
		message.From.LastName,
	)
	if err != nil {
		log.Printf("获取用户失败: %v", err)
		return
	}

	// 检查是否有进行中的会话
	if session, exists := h.sessionManager.GetSession(message.From.ID); exists {
		h.handleSessionMessage(message, user, session)
		return
	}

	// 处理命令
	if message.IsCommand() {
		h.handleCommand(message, user)
		return
	}

	// 非命令消息
	h.sendMessage(message.Chat.ID, "请使用 /help 查看可用命令")
}

// handleCommand 处理命令
func (h *MasterHandler) handleCommand(message *tgbotapi.Message, user *models.User) {
	command := message.Command()
	args := message.CommandArguments()

	// 管理员命令（命令以 admin_ 开头）
	if strings.HasPrefix(command, "admin_") && user.IsAdmin {
		h.handleAdminCommand(message, command, args)
		return
	}

	// 特殊命令：设置管理员（仅在没有管理员时可用）
	if command == "setadmin" {
		h.handleSetAdmin(message, user)
		return
	}

	// 普通用户命令
	switch command {
	case "start":
		h.handleStart(message, user)
	case "help":
		h.handleHelp(message, user)
	case "buy":
		h.handleBuy(message, user)
	case "redeem":
		h.handleRedeem(message, user)
	case "myplan":
		h.handleMyPlan(message, user)
	case "mybots":
		h.handleMyBots(message, user)
	case "createbot":
		h.handleCreateBot(message, user)
	case "managebot":
		h.handleManageBot(message, user)
	default:
		h.sendMessage(message.Chat.ID, "未知命令，使用 /help 查看可用命令")
	}
}

// handleStart 处理 /start 命令
func (h *MasterHandler) handleStart(message *tgbotapi.Message, user *models.User) {
	// 检查是否有管理员
	admins, _ := h.userRepo.GetAllAdmins()
	if len(admins) == 0 {
		// 没有管理员，提示设置
		text := fmt.Sprintf(`👋 欢迎使用Telegram消息转发Bot系统！

🆔 您的信息：
用户名：@%s
姓名：%s %s
Telegram ID：%d

⚠️ 系统提示：
当前系统还没有管理员，请使用以下命令设置您为管理员：

/setadmin

设置完成后，使用 /help 查看可用命令`,
			user.Username,
			user.FirstName,
			user.LastName,
			user.TelegramID,
		)
		h.sendMessage(message.Chat.ID, text)
		return
	}

	text := fmt.Sprintf(`👋 欢迎使用Telegram消息转发Bot系统！

🆔 您的信息：
用户名：@%s
姓名：%s %s
Telegram ID：%d

📋 快速开始：
1. /buy - 购买订阅套餐
2. /createbot - 创建您的转发Bot
3. /help - 查看完整帮助

💡 提示：首次使用请先购买订阅套餐`,
		user.Username,
		user.FirstName,
		user.LastName,
		user.TelegramID,
	)

	h.sendMessage(message.Chat.ID, text)
}

// handleSetAdmin 处理 /setadmin 命令
func (h *MasterHandler) handleSetAdmin(message *tgbotapi.Message, user *models.User) {
	// 检查是否已有管理员
	admins, err := h.userRepo.GetAllAdmins()
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 查询管理员失败："+err.Error())
		return
	}

	if len(admins) > 0 {
		h.sendMessage(message.Chat.ID, "❌ 系统已有管理员，无法设置\n\n如需更改管理员，请联系现有管理员使用管理功能")
		return
	}

	// 设置为管理员
	if err := h.userRepo.SetAdmin(user.TelegramID, true); err != nil {
		h.sendMessage(message.Chat.ID, "❌ 设置管理员失败："+err.Error())
		return
	}

	text := fmt.Sprintf(`✅ 管理员设置成功！

👤 管理员信息：
用户名：@%s
姓名：%s %s
Telegram ID：%d

🔧 管理员专用命令：
使用 /admin_help 查看所有管理员命令

📋 快速开始：
1. /admin_add_plan - 创建订阅套餐
2. /admin_payment_config - 配置支付接口
3. /admin_ai_config - 配置AI反垃圾
4. /admin_stats - 查看系统统计`,
		user.Username,
		user.FirstName,
		user.LastName,
		user.TelegramID,
	)

	h.sendMessage(message.Chat.ID, text)
}

// handleHelp 处理 /help 命令
func (h *MasterHandler) handleHelp(message *tgbotapi.Message, user *models.User) {
	helpText := `📖 帮助文档

【订阅管理】
/start - 开始使用/查看状态
/buy - 购买订阅套餐
/redeem - 使用兑换码
/myplan - 查看我的订阅信息

【子Bot管理】
/mybots - 查看我的子Bot列表
/createbot - 创建新的子Bot
/managebot - 管理已有的子Bot

【其他】
/help - 显示此帮助信息

💡 提示：创建子Bot后，请在子Bot中发送 /help 查看子Bot专用指令`

	if user.IsAdmin {
		helpText += `

【管理员命令】
使用 /admin_help 查看管理员专用命令`
	}

	h.sendMessage(message.Chat.ID, helpText)
}

// handleBuy 处理 /buy 命令
func (h *MasterHandler) handleBuy(message *tgbotapi.Message, user *models.User) {
	h.sendMessage(message.Chat.ID, "💳 购买订阅功能开发中，请稍后...")
}

// handleRedeem 处理 /redeem 命令
func (h *MasterHandler) handleRedeem(message *tgbotapi.Message, user *models.User) {
	h.sendMessage(message.Chat.ID, "🎫 兑换码功能开发中，请稍后...")
}

// handleMyPlan 处理 /myplan 命令
func (h *MasterHandler) handleMyPlan(message *tgbotapi.Message, user *models.User) {
	h.sendMessage(message.Chat.ID, "📋 订阅信息功能开发中，请稍后...")
}

// handleMyBots 处理 /mybots 命令
func (h *MasterHandler) handleMyBots(message *tgbotapi.Message, user *models.User) {
	h.sendMessage(message.Chat.ID, "🤖 子Bot列表功能开发中，请稍后...")
}

// handleCreateBot 处理 /createbot 命令
func (h *MasterHandler) handleCreateBot(message *tgbotapi.Message, user *models.User) {
	h.sendMessage(message.Chat.ID, "➕ 创建子Bot功能开发中，请稍后...")
}

// handleManageBot 处理 /managebot 命令
func (h *MasterHandler) handleManageBot(message *tgbotapi.Message, user *models.User) {
	h.sendMessage(message.Chat.ID, "🔧 管理子Bot功能开发中，请稍后...")
}

// handleAdminCommand 处理管理员命令
func (h *MasterHandler) handleAdminCommand(message *tgbotapi.Message, command, args string) {
	switch command {
	case "admin_help":
		h.handleAdminHelp(message)
	case "admin_stats":
		h.handleAdminStats(message)
	case "admin_add_plan":
		h.handleAdminAddPlan(message)
	case "admin_list_plans":
		h.handleAdminListPlans(message)
	case "admin_generate_code":
		h.handleAdminGenerateCode(message)
	case "admin_code_list":
		h.handleAdminCodeList(message)
	case "admin_payment_config":
		h.handleAdminPaymentConfig(message)
	case "admin_payment_status":
		h.handleAdminPaymentStatus(message)
	case "admin_ai_config":
		h.handleAdminAIConfig(message)
	case "admin_ai_test":
		h.handleAdminAITest(message)
	case "admin_users":
		h.handleAdminUsers(message)
	case "admin_backup":
		h.handleAdminBackup(message)
	default:
		h.sendMessage(message.Chat.ID, "未知管理员命令，使用 /admin_help 查看可用命令")
	}
}

// handleAdminHelp 处理 /admin_help 命令
func (h *MasterHandler) handleAdminHelp(message *tgbotapi.Message) {
	helpText := `🔧 管理员指令

【支付管理】
/admin_payment_config - 配置支付接口
/admin_payment_status - 查看支付状态

【套餐管理】
/admin_add_plan - 创建新套餐
/admin_list_plans - 查看/编辑套餐

【兑换码管理】
/admin_generate_code - 生成兑换码
/admin_code_list - 查看兑换码列表

【AI配置】
/admin_ai_config - 配置默认AI
/admin_ai_test - 测试AI连接

【用户管理】
/admin_stats - 系统统计
/admin_users - 用户列表

【系统】
/admin_backup - 立即备份数据库`

	h.sendMessage(message.Chat.ID, helpText)
}

// handleAdminStats 处理 /admin_stats 命令
func (h *MasterHandler) handleAdminStats(message *tgbotapi.Message) {
	stats, err := h.adminService.GetSystemStats()
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取统计信息失败："+err.Error())
		return
	}

	text := fmt.Sprintf(`📊 系统统计

👥 用户统计：
总用户数：%d
管理员数：%d

📦 订阅统计：
总订阅数：%d
有效订阅：%d

🤖 Bot统计：
总Bot数：%d
运行中：%d

💰 套餐统计：
有效套餐：%d

🎫 兑换码统计：
总兑换码：%d
未使用：%d`,
		stats["user_count"],
		stats["admin_count"],
		stats["subscription_count"],
		stats["active_subscription_count"],
		stats["bot_count"],
		stats["running_bot_count"],
		stats["plan_count"],
		stats["redeem_code_count"],
		stats["unused_code_count"],
	)

	h.sendMessage(message.Chat.ID, text)
}

// HandleCallback 处理回调查询
func (h *MasterHandler) HandleCallback(callback *tgbotapi.CallbackQuery) {
	// 回调功能待实现
	h.bot.Send(tgbotapi.NewCallback(callback.ID, "功能开发中"))
}

// sendMessage 发送消息
func (h *MasterHandler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}
