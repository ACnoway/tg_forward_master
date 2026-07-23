package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/acnoway/tg_forward_master/internal/master/service"
	"github.com/acnoway/tg_forward_master/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleAdminAddPlan 处理 /admin_add_plan 命令
func (h *MasterHandler) handleAdminAddPlan(message *tgbotapi.Message) {
	h.sessionManager.StartSession(message.From.ID, "admin_add_plan")
	h.sendMessage(message.Chat.ID, `📦 创建新套餐

请输入套餐名称：`)
}

// handleAdminListPlans 处理 /admin_list_plans 命令
func (h *MasterHandler) handleAdminListPlans(message *tgbotapi.Message) {
	plans, err := h.adminService.GetAllPlans()
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取套餐列表失败："+err.Error())
		return
	}

	if len(plans) == 0 {
		h.sendMessage(message.Chat.ID, "📦 当前没有套餐\n\n使用 /admin_add_plan 创建新套餐")
		return
	}

	var text strings.Builder
	text.WriteString("📦 套餐列表\n\n")

	for i, plan := range plans {
		status := "✅"
		if !plan.IsActive {
			status = "❌"
		}
		text.WriteString(fmt.Sprintf("%s %d. %s\n", status, i+1, plan.Name))
		text.WriteString(fmt.Sprintf("   ID: %d\n", plan.ID))
		text.WriteString(fmt.Sprintf("   价格: %s\n", service.FormatPrice(plan.Price)))
		text.WriteString(fmt.Sprintf("   有效期: %s\n", service.FormatDuration(plan.DurationDays)))
		text.WriteString(fmt.Sprintf("   Bot数量: %d\n", plan.MaxBots))
		text.WriteString(fmt.Sprintf("   创建时间: %s\n\n", service.FormatTime(plan.CreatedAt)))
	}

	text.WriteString("💡 提示：使用 /admin_add_plan 创建新套餐")

	h.sendMessage(message.Chat.ID, text.String())
}

// handleAdminGenerateCode 处理 /admin_generate_code 命令
func (h *MasterHandler) handleAdminGenerateCode(message *tgbotapi.Message) {
	// 先显示套餐列表
	plans, err := h.planRepo.GetActive()
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取套餐列表失败："+err.Error())
		return
	}

	if len(plans) == 0 {
		h.sendMessage(message.Chat.ID, "📦 当前没有可用套餐\n\n请先使用 /admin_add_plan 创建套餐")
		return
	}

	var text strings.Builder
	text.WriteString("🎫 生成兑换码\n\n")
	text.WriteString("请选择套餐（输入套餐ID）：\n\n")

	for _, plan := range plans {
		text.WriteString(fmt.Sprintf("ID: %d - %s (%s, %s, %d个Bot)\n",
			plan.ID,
			plan.Name,
			service.FormatPrice(plan.Price),
			service.FormatDuration(plan.DurationDays),
			plan.MaxBots,
		))
	}

	h.sessionManager.StartSession(message.From.ID, "admin_generate_code")
	h.sendMessage(message.Chat.ID, text.String())
}

// handleAdminCodeList 处理 /admin_code_list 命令
func (h *MasterHandler) handleAdminCodeList(message *tgbotapi.Message) {
	codes, err := h.adminService.GetRedeemCodes(50, 0)
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取兑换码列表失败："+err.Error())
		return
	}

	if len(codes) == 0 {
		h.sendMessage(message.Chat.ID, "🎫 当前没有兑换码\n\n使用 /admin_generate_code 生成兑换码")
		return
	}

	var text strings.Builder
	text.WriteString("🎫 兑换码列表（最近50个）\n\n")

	for i, code := range codes {
		plan, _ := h.planRepo.GetByID(code.PlanID)
		planName := "未知套餐"
		if plan != nil {
			planName = plan.Name
		}

		statusIcon := "🟢"
		statusText := "未使用"
		if code.Status == "used" {
			statusIcon = "🔴"
			statusText = "已使用"
		}

		text.WriteString(fmt.Sprintf("%d. %s `%s`\n", i+1, statusIcon, code.Code))
		text.WriteString(fmt.Sprintf("   套餐: %s\n", planName))
		text.WriteString(fmt.Sprintf("   状态: %s\n", statusText))
		if code.UsedBy != nil {
			text.WriteString(fmt.Sprintf("   使用者ID: %d\n", *code.UsedBy))
		}
		text.WriteString(fmt.Sprintf("   创建时间: %s\n\n", service.FormatTime(code.CreatedAt)))
	}

	h.sendMessage(message.Chat.ID, text.String())
}

// handleAdminPaymentConfig 处理 /admin_payment_config 命令
func (h *MasterHandler) handleAdminPaymentConfig(message *tgbotapi.Message) {
	h.sessionManager.StartSession(message.From.ID, "admin_payment_config")
	h.sendMessage(message.Chat.ID, `💳 配置易支付

请输入易支付API地址：
例如：https://pay.example.com/submit.php`)
}

// handleAdminPaymentStatus 处理 /admin_payment_status 命令
func (h *MasterHandler) handleAdminPaymentStatus(message *tgbotapi.Message) {
	apiURL, merchantID, _, err := h.adminService.GetPaymentConfig()
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取支付配置失败："+err.Error())
		return
	}

	status := "❌ 未配置"
	if apiURL != "" && merchantID != "" {
		status = "✅ 已配置"
	}

	text := fmt.Sprintf(`💳 易支付状态

状态：%s

配置信息：
API地址：%s
商户ID：%s
商户密钥：%s

使用 /admin_payment_config 修改配置`,
		status,
		getConfigDisplay(apiURL),
		getConfigDisplay(merchantID),
		getConfigDisplay("***"),
	)

	h.sendMessage(message.Chat.ID, text)
}

// handleAdminAIConfig 处理 /admin_ai_config 命令
func (h *MasterHandler) handleAdminAIConfig(message *tgbotapi.Message) {
	h.sessionManager.StartSession(message.From.ID, "admin_ai_config")
	h.sendMessage(message.Chat.ID, `🤖 配置默认AI

请输入API地址：
例如：https://api.openai.com/v1/chat/completions`)
}

// getConfigDisplay 获取配置显示文本
func getConfigDisplay(value string) string {
	if value == "" {
		return "未设置"
	}
	return value
}

// handleAdminAITest 处理 /admin_ai_test 命令
func (h *MasterHandler) handleAdminAITest(message *tgbotapi.Message) {
	h.sendMessage(message.Chat.ID, "🤖 正在测试AI连接...\n\n请稍候...")

	// 测试AI连接
	response, err := h.adminService.TestAIConnection()
	if err != nil {
		text := fmt.Sprintf("❌ AI连接测试失败\n\n错误信息：\n%s\n\n请检查：\n1. API地址是否正确\n2. API密钥是否有效\n3. 模型名称是否正确\n4. 网络连接是否正常", err.Error())
		h.sendMessage(message.Chat.ID, text)
		return
	}

	text := fmt.Sprintf("✅ AI连接测试成功！\n\nAI响应：\n%s\n\n连接正常，可以使用AI功能", response)
	h.sendMessage(message.Chat.ID, text)
}

// handleAdminUsers 处理 /admin_users 命令
func (h *MasterHandler) handleAdminUsers(message *tgbotapi.Message) {
	// 获取用户列表（显示前20个）
	users, total, err := h.adminService.GetUserList(20, 0)
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取用户列表失败："+err.Error())
		return
	}

	if total == 0 {
		h.sendMessage(message.Chat.ID, "👥 当前没有用户")
		return
	}

	var text strings.Builder
	text.WriteString(fmt.Sprintf("👥 用户列表（共 %d 人，显示前20个）\n\n", total))

	for i, user := range users {
		roleIcon := "👤"
		if user.IsAdmin {
			roleIcon = "👑"
		}

		name := user.FirstName
		if user.LastName != "" {
			name += " " + user.LastName
		}
		if name == "" {
			name = "未知"
		}

		username := ""
		if user.Username != "" {
			username = "@" + user.Username
		}

		text.WriteString(fmt.Sprintf("%s %d. %s\n", roleIcon, i+1, name))
		if username != "" {
			text.WriteString(fmt.Sprintf("   用户名: %s\n", username))
		}
		text.WriteString(fmt.Sprintf("   ID: %d\n", user.TelegramID))
		text.WriteString(fmt.Sprintf("   注册时间: %s\n\n", service.FormatTime(user.CreatedAt)))
	}

	if total > 20 {
		text.WriteString(fmt.Sprintf("\n💡 还有 %d 个用户未显示", total-20))
	}

	h.sendMessage(message.Chat.ID, text.String())
}

// handleAdminBackup 处理 /admin_backup 命令
func (h *MasterHandler) handleAdminBackup(message *tgbotapi.Message) {
	h.sendMessage(message.Chat.ID, "💾 正在备份数据库...\n\n请稍候...")

	// 执行备份
	backupPath, err := h.adminService.BackupDatabase(h.config.DatabasePath)
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 备份失败："+err.Error())
		return
	}

	// 获取文件信息
	info, err := os.Stat(backupPath)
	if err != nil {
		h.sendMessage(message.Chat.ID, "❌ 获取备份文件信息失败："+err.Error())
		return
	}

	// 格式化文件大小
	size := float64(info.Size())
	unit := "B"
	if size > 1024 {
		size = size / 1024
		unit = "KB"
	}
	if size > 1024 {
		size = size / 1024
		unit = "MB"
	}

	text := fmt.Sprintf(`✅ 数据库备份成功！

📦 备份信息：
文件名: %s
大小: %.2f %s
时间: %s

💡 备份文件保存在与数据库相同的目录下`,
		filepath.Base(backupPath),
		size,
		unit,
		service.FormatTime(info.ModTime()),
	)

	h.sendMessage(message.Chat.ID, text)

	// 显示最近的备份列表
	backups, err := h.adminService.GetBackupList(h.config.DatabasePath)
	if err == nil && len(backups) > 0 {
		var listText strings.Builder
		listText.WriteString("\n\n📋 最近的备份（最多5个）：\n\n")

		count := len(backups)
		if count > 5 {
			count = 5
		}

		for i := 0; i < count; i++ {
			backup := backups[i]
			bSize := float64(backup.Size)
			bUnit := "B"
			if bSize > 1024 {
				bSize = bSize / 1024
				bUnit = "KB"
			}
			if bSize > 1024 {
				bSize = bSize / 1024
				bUnit = "MB"
			}

			listText.WriteString(fmt.Sprintf("%d. %s\n", i+1, filepath.Base(backup.Path)))
			listText.WriteString(fmt.Sprintf("   大小: %.2f %s\n", bSize, bUnit))
			listText.WriteString(fmt.Sprintf("   时间: %s\n\n", service.FormatTime(backup.ModTime)))
		}

		h.sendMessage(message.Chat.ID, listText.String())
	}
}

// handleSessionMessage 处理会话中的消息
func (h *MasterHandler) handleSessionMessage(message *tgbotapi.Message, user *models.User, session *SessionState) {
	switch session.Command {
	case "admin_add_plan":
		h.handleAddPlanSession(message, user, session)
	case "admin_generate_code":
		h.handleGenerateCodeSession(message, user, session)
	case "admin_payment_config":
		h.handlePaymentConfigSession(message, user, session)
	case "admin_ai_config":
		h.handleAIConfigSession(message, user, session)
	default:
		h.sessionManager.EndSession(message.From.ID)
		h.sendMessage(message.Chat.ID, "会话已过期，请重新开始")
	}
}

// handleAddPlanSession 处理创建套餐会话
func (h *MasterHandler) handleAddPlanSession(message *tgbotapi.Message, user *models.User, session *SessionState) {
	switch session.Step {
	case 0: // 等待套餐名称
		name := strings.TrimSpace(message.Text)
		if name == "" {
			h.sendMessage(message.Chat.ID, "❌ 套餐名称不能为空，请重新输入：")
			return
		}
		h.sessionManager.UpdateSession(message.From.ID, 1, map[string]interface{}{"name": name})
		h.sendMessage(message.Chat.ID, "✅ 套餐名称："+name+"\n\n请输入套餐价格（元）：")

	case 1: // 等待价格
		price, err := strconv.ParseFloat(message.Text, 64)
		if err != nil || price < 0 {
			h.sendMessage(message.Chat.ID, "❌ 价格格式不正确，请输入有效的数字：")
			return
		}
		h.sessionManager.UpdateSession(message.From.ID, 2, map[string]interface{}{"price": price})
		h.sendMessage(message.Chat.ID, fmt.Sprintf("✅ 价格：%.2f 元\n\n请输入有效期（天数）：", price))

	case 2: // 等待有效期
		days, err := strconv.Atoi(message.Text)
		if err != nil || days <= 0 {
			h.sendMessage(message.Chat.ID, "❌ 天数格式不正确，请输入大于0的整数：")
			return
		}
		h.sessionManager.UpdateSession(message.From.ID, 3, map[string]interface{}{"duration_days": days})
		h.sendMessage(message.Chat.ID, fmt.Sprintf("✅ 有效期：%d 天\n\n请输入可创建的Bot数量：", days))

	case 3: // 等待Bot数量
		maxBots, err := strconv.Atoi(message.Text)
		if err != nil || maxBots <= 0 {
			h.sendMessage(message.Chat.ID, "❌ Bot数量格式不正确，请输入大于0的整数：")
			return
		}

		// 获取之前输入的数据
		name := session.Data["name"].(string)
		price := session.Data["price"].(float64)
		days := session.Data["duration_days"].(int)

		// 创建套餐
		plan, err := h.adminService.CreatePlan(name, price, days, maxBots)
		if err != nil {
			h.sendMessage(message.Chat.ID, "❌ 创建套餐失败："+err.Error())
			h.sessionManager.EndSession(message.From.ID)
			return
		}

		text := fmt.Sprintf(`✅ 套餐创建成功！

📦 套餐信息：
名称：%s
价格：%s
有效期：%s
Bot数量：%d 个
套餐ID：%d

使用 /admin_list_plans 查看所有套餐`,
			plan.Name,
			service.FormatPrice(plan.Price),
			service.FormatDuration(plan.DurationDays),
			plan.MaxBots,
			plan.ID,
		)

		h.sendMessage(message.Chat.ID, text)
		h.sessionManager.EndSession(message.From.ID)
	}
}

// handleGenerateCodeSession 处理生成兑换码会话
func (h *MasterHandler) handleGenerateCodeSession(message *tgbotapi.Message, user *models.User, session *SessionState) {
	switch session.Step {
	case 0: // 等待套餐ID
		planID, err := strconv.ParseInt(message.Text, 10, 64)
		if err != nil {
			h.sendMessage(message.Chat.ID, "❌ 套餐ID格式不正确，请输入有效的数字：")
			return
		}

		// 检查套餐是否存在
		plan, err := h.planRepo.GetByID(planID)
		if err != nil || plan == nil {
			h.sendMessage(message.Chat.ID, "❌ 套餐不存在，请重新输入：")
			return
		}

		h.sessionManager.UpdateSession(message.From.ID, 1, map[string]interface{}{
			"plan_id":   planID,
			"plan_name": plan.Name,
		})
		h.sendMessage(message.Chat.ID, fmt.Sprintf("✅ 已选择套餐：%s\n\n请输入生成数量（1-100）：", plan.Name))

	case 1: // 等待生成数量
		count, err := strconv.Atoi(message.Text)
		if err != nil || count <= 0 || count > 100 {
			h.sendMessage(message.Chat.ID, "❌ 数量必须在1-100之间，请重新输入：")
			return
		}

		planID := session.Data["plan_id"].(int64)
		planName := session.Data["plan_name"].(string)

		// 生成兑换码
		codes, err := h.adminService.GenerateRedeemCodes(planID, count)
		if err != nil {
			h.sendMessage(message.Chat.ID, "❌ 生成兑换码失败："+err.Error())
			h.sessionManager.EndSession(message.From.ID)
			return
		}

		// 构建兑换码列表
		var codeList strings.Builder
		codeList.WriteString(fmt.Sprintf("✅ 成功生成 %d 个兑换码！\n\n", count))
		codeList.WriteString(fmt.Sprintf("📦 套餐：%s\n\n", planName))
		codeList.WriteString("🎫 兑换码列表：\n")
		for i, code := range codes {
			codeList.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, code.Code))
		}
		codeList.WriteString("\n💡 用户可使用 /redeem 命令兑换")

		h.sendMessage(message.Chat.ID, codeList.String())
		h.sessionManager.EndSession(message.From.ID)
	}
}

// handlePaymentConfigSession 处理支付配置会话
func (h *MasterHandler) handlePaymentConfigSession(message *tgbotapi.Message, user *models.User, session *SessionState) {
	switch session.Step {
	case 0: // 等待API地址
		apiURL := strings.TrimSpace(message.Text)
		if apiURL == "" {
			h.sendMessage(message.Chat.ID, "❌ API地址不能为空，请重新输入：")
			return
		}
		h.sessionManager.UpdateSession(message.From.ID, 1, map[string]interface{}{"api_url": apiURL})
		h.sendMessage(message.Chat.ID, "✅ API地址已设置\n\n请输入商户ID：")

	case 1: // 等待商户ID
		merchantID := strings.TrimSpace(message.Text)
		if merchantID == "" {
			h.sendMessage(message.Chat.ID, "❌ 商户ID不能为空，请重新输入：")
			return
		}
		h.sessionManager.UpdateSession(message.From.ID, 2, map[string]interface{}{"merchant_id": merchantID})
		h.sendMessage(message.Chat.ID, "✅ 商户ID已设置\n\n请输入商户密钥：")

	case 2: // 等待商户密钥
		merchantKey := strings.TrimSpace(message.Text)
		if merchantKey == "" {
			h.sendMessage(message.Chat.ID, "❌ 商户密钥不能为空，请重新输入：")
			return
		}

		apiURL := session.Data["api_url"].(string)
		merchantID := session.Data["merchant_id"].(string)

		// 保存配置
		err := h.adminService.SetPaymentConfig(apiURL, merchantID, merchantKey)
		if err != nil {
			h.sendMessage(message.Chat.ID, "❌ 保存配置失败："+err.Error())
			h.sessionManager.EndSession(message.From.ID)
			return
		}

		text := `✅ 易支付配置已保存！

🔧 配置信息：
API地址：已设置
商户ID：已设置
商户密钥：已设置

现在用户可以通过 /buy 购买订阅了`

		h.sendMessage(message.Chat.ID, text)
		h.sessionManager.EndSession(message.From.ID)
	}
}

// handleAIConfigSession 处理AI配置会话
func (h *MasterHandler) handleAIConfigSession(message *tgbotapi.Message, user *models.User, session *SessionState) {
	switch session.Step {
	case 0: // 等待API地址
		endpoint := strings.TrimSpace(message.Text)
		if endpoint == "" {
			h.sendMessage(message.Chat.ID, "❌ API地址不能为空，请重新输入：")
			return
		}
		h.sessionManager.UpdateSession(message.From.ID, 1, map[string]interface{}{"endpoint": endpoint})
		h.sendMessage(message.Chat.ID, "✅ API地址已设置\n\n请输入API密钥：")

	case 1: // 等待API密钥
		apiKey := strings.TrimSpace(message.Text)
		if apiKey == "" {
			h.sendMessage(message.Chat.ID, "❌ API密钥不能为空，请重新输入：")
			return
		}
		h.sessionManager.UpdateSession(message.From.ID, 2, map[string]interface{}{"api_key": apiKey})
		h.sendMessage(message.Chat.ID, "✅ API密钥已设置\n\n请输入模型名称（如：gpt-3.5-turbo）：")

	case 2: // 等待模型名称
		model := strings.TrimSpace(message.Text)
		if model == "" {
			h.sendMessage(message.Chat.ID, "❌ 模型名称不能为空，请重新输入：")
			return
		}

		endpoint := session.Data["endpoint"].(string)
		apiKey := session.Data["api_key"].(string)

		// 保存配置
		err := h.adminService.SetAIConfig(endpoint, apiKey, model)
		if err != nil {
			h.sendMessage(message.Chat.ID, "❌ 保存配置失败："+err.Error())
			h.sessionManager.EndSession(message.From.ID)
			return
		}

		text := fmt.Sprintf(`✅ 默认AI配置已保存！

🤖 配置信息：
API地址：%s
模型：%s

此配置将作为所有子Bot的默认AI配置
用户可以在子Bot中使用 /ai_custom 切换到自定义配置`, endpoint, model)

		h.sendMessage(message.Chat.ID, text)
		h.sessionManager.EndSession(message.From.ID)
	}
}
