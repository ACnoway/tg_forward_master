package models

import "time"

// User 用户模型
type User struct {
	ID          int64     `json:"id"`
	TelegramID  int64     `json:"telegram_id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	IsAdmin     bool      `json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
}

// Subscription 订阅模型
type Subscription struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PlanID    int64     `json:"plan_id"`
	Status    string    `json:"status"` // active, expired, banned
	ExpiresAt time.Time `json:"expires_at"`
	MaxBots   int       `json:"max_bots"`
	CreatedAt time.Time `json:"created_at"`
}

// Plan 套餐模型
type Plan struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	DurationDays int       `json:"duration_days"`
	MaxBots      int       `json:"max_bots"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// RedeemCode 兑换码模型
type RedeemCode struct {
	ID        int64      `json:"id"`
	Code      string     `json:"code"`
	PlanID    int64      `json:"plan_id"`
	Status    string     `json:"status"` // unused, used, expired
	UsedBy    *int64     `json:"used_by"`
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at"`
}

// WorkerBot 子Bot模型
type WorkerBot struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	BotToken        string    `json:"bot_token"`
	BotUsername     string    `json:"bot_username"`
	BotNickname     string    `json:"bot_nickname"`
	OwnerTelegramID int64     `json:"owner_telegram_id"`
	Status          string    `json:"status"` // running, stopped
	CreatedAt       time.Time `json:"created_at"`
}

// BotConfig 子Bot配置模型
type BotConfig struct {
	BotID             int64     `json:"bot_id"`
	OwnerTelegramID   int64     `json:"owner_telegram_id"`
	AIEnabled         bool      `json:"ai_enabled"`
	UseCustomAI       bool      `json:"use_custom_ai"`
	CustomAIEndpoint  string    `json:"custom_ai_endpoint"`
	CustomAIKey       string    `json:"custom_ai_key"`
	CustomAIModel     string    `json:"custom_ai_model"`
	WhitelistThreshold int      `json:"whitelist_threshold"`
	NotifyOnBlock     bool      `json:"notify_on_block"`
	NotifyOnAIBlock   bool      `json:"notify_on_ai_block"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Customer 客户模型
type Customer struct {
	ID              int64      `json:"id"`
	BotID           int64      `json:"bot_id"`
	TelegramID      int64      `json:"telegram_id"`
	Username        string     `json:"username"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	VerifiedCount   int        `json:"verified_count"`
	TotalMessages   int        `json:"total_messages"`
	IsWhitelisted   bool       `json:"is_whitelisted"`
	IsBlacklisted   bool       `json:"is_blacklisted"`
	FirstSeenAt     time.Time  `json:"first_seen_at"`
	LastMessageAt   *time.Time `json:"last_message_at"`
	WhitelistedAt   *time.Time `json:"whitelisted_at"`
	BlacklistedAt   *time.Time `json:"blacklisted_at"`
	BlacklistReason string     `json:"blacklist_reason"`
	IsManualTrust   bool       `json:"is_manual_trust"`
}

// GetDisplayName 获取客户显示名称
func (c *Customer) GetDisplayName() string {
	if c.FirstName != "" {
		if c.LastName != "" {
			return c.FirstName + " " + c.LastName
		}
		return c.FirstName
	}
	if c.Username != "" {
		return "@" + c.Username
	}
	return "User" + string(rune(c.TelegramID))
}

// GetFullIdentifier 获取完整标识
func (c *Customer) GetFullIdentifier() string {
	parts := ""

	// 昵称
	if c.FirstName != "" {
		parts = c.FirstName
		if c.LastName != "" {
			parts += " " + c.LastName
		}
		parts += " | "
	}

	// 用户名
	if c.Username != "" {
		parts += "@" + c.Username + " | "
	}

	// ID
	parts += "ID:" + string(rune(c.TelegramID))

	return parts
}

// BlockLog 拦截日志模型
type BlockLog struct {
	ID            int64     `json:"id"`
	BotID         int64     `json:"bot_id"`
	TelegramID    int64     `json:"telegram_id"`
	Username      string    `json:"username"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	BlockReason   string    `json:"block_reason"` // blacklist, ai_spam
	AIConfidence  *float64  `json:"ai_confidence"`
	MessageContent string   `json:"message_content"`
	MessageType   string    `json:"message_type"`
	BlockedAt     time.Time `json:"blocked_at"`
	IsFalsePositive bool    `json:"is_false_positive"`
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Order 订单模型
type Order struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	PlanID     int64      `json:"plan_id"`
	OrderNo    string     `json:"order_no"`
	Amount     float64    `json:"amount"`
	Status     string     `json:"status"` // pending, paid, failed
	PaymentURL string     `json:"payment_url"`
	PaidAt     *time.Time `json:"paid_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
