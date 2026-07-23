package forwarder

import (
	"fmt"
	"log"

	"github.com/acnoway/tg_forward_master/internal/database"
	"github.com/acnoway/tg_forward_master/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageForwarder 消息转发器
type MessageForwarder struct {
	bot          *tgbotapi.BotAPI
	db           *database.DB
	botID        int64
	ownerID      int64
	customerRepo *database.CustomerRepository
}

// NewMessageForwarder 创建消息转发器
func NewMessageForwarder(bot *tgbotapi.BotAPI, db *database.DB, botID, ownerID int64) *MessageForwarder {
	return &MessageForwarder{
		bot:          bot,
		db:           db,
		botID:        botID,
		ownerID:      ownerID,
		customerRepo: database.NewCustomerRepository(db),
	}
}

// ForwardToOwner 转发消息给主人
func (f *MessageForwarder) ForwardToOwner(message *tgbotapi.Message) error {
	// 获取或创建客户信息
	customer, err := f.customerRepo.GetOrCreate(
		f.botID,
		message.From.ID,
		message.From.UserName,
		message.From.FirstName,
		message.From.LastName,
	)
	if err != nil {
		return fmt.Errorf("获取客户信息失败: %w", err)
	}

	// 构建消息头
	header := f.buildMessageHeader(customer)

	// 转发原始消息
	forwardMsg := tgbotapi.NewForward(f.ownerID, message.Chat.ID, message.MessageID)
	sentMsg, err := f.bot.Send(forwardMsg)
	if err != nil {
		return fmt.Errorf("转发消息失败: %w", err)
	}

	// 发送消息头和快捷操作
	headerMsg := tgbotapi.NewMessage(f.ownerID, header)
	headerMsg.ReplyToMessageID = sentMsg.MessageID
	headerMsg.ReplyMarkup = f.buildQuickActions(customer.TelegramID)

	if _, err := f.bot.Send(headerMsg); err != nil {
		log.Printf("发送消息头失败: %v", err)
	}

	// 更新消息计数
	if err := f.customerRepo.IncrementMessageCount(f.botID, message.From.ID); err != nil {
		log.Printf("更新消息计数失败: %v", err)
	}

	return nil
}

// buildMessageHeader 构建消息头
func (f *MessageForwarder) buildMessageHeader(customer *models.Customer) string {
	status := "🟢 新用户"
	if customer.IsWhitelisted {
		status = fmt.Sprintf("✅ 已验证（白名单，%d条消息）", customer.TotalMessages)
	} else if customer.VerifiedCount > 0 {
		status = fmt.Sprintf("🟡 验证中（%d/%d条）", customer.VerifiedCount, 3)
	}

	return fmt.Sprintf(`【收到新消息】

👤 发件人信息
姓名：%s
用户名：%s
ID：%d
状态：%s

---
⚡ 快捷操作：
回复此消息 = 回复客户
/block_%d = 拉黑此人
/trust_%d = 信任此人`,
		customer.GetDisplayName(),
		func() string {
			if customer.Username != "" {
				return "@" + customer.Username
			}
			return "无"
		}(),
		customer.TelegramID,
		status,
		customer.TelegramID,
		customer.TelegramID,
	)
}

// buildQuickActions 构建快捷操作按钮
func (f *MessageForwarder) buildQuickActions(customerID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚫 拉黑", fmt.Sprintf("block_%d", customerID)),
			tgbotapi.NewInlineKeyboardButtonData("✅ 信任", fmt.Sprintf("trust_%d", customerID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 查看详情", fmt.Sprintf("info_%d", customerID)),
		),
	)
}

// ReplyToCustomer 回复客户
func (f *MessageForwarder) ReplyToCustomer(ownerMessage *tgbotapi.Message, customerID int64) error {
	// 检查是否是回复消息
	if ownerMessage.ReplyToMessage == nil {
		return fmt.Errorf("请回复客户的消息")
	}

	// 复制消息到客户
	copyMsg := tgbotapi.NewCopyMessage(customerID, ownerMessage.Chat.ID, ownerMessage.MessageID)
	if _, err := f.bot.Send(copyMsg); err != nil {
		return fmt.Errorf("回复客户失败: %w", err)
	}

	// 通知主人发送成功
	confirmMsg := tgbotapi.NewMessage(f.ownerID, "✅ 已发送给客户")
	confirmMsg.ReplyToMessageID = ownerMessage.MessageID
	f.bot.Send(confirmMsg)

	return nil
}
