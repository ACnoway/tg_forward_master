package spam

import (
	"fmt"
	"log"

	"github.com/acnoway/tg_forward_master/internal/config"
	"github.com/acnoway/tg_forward_master/internal/database"
)

// AISpamChecker AI垃圾消息检查器
type AISpamChecker struct {
	db     *database.DB
	botID  int64
	config *config.Config
}

// NewAISpamChecker 创建AI垃圾消息检查器
func NewAISpamChecker(db *database.DB, botID int64) *AISpamChecker {
	// 加载配置
	cfg, err := config.Load("config.env")
	if err != nil {
		log.Printf("加载配置失败: %v", err)
		cfg = &config.Config{}
	}

	return &AISpamChecker{
		db:     db,
		botID:  botID,
		config: cfg,
	}
}

// CheckSpam 检查消息是否为垃圾
func (s *AISpamChecker) CheckSpam(message string) (int, error) {
	// 获取Bot配置
	configRepo := database.NewBotConfigRepository(s.db)
	botConfig, err := configRepo.GetByBotID(s.botID)
	if err != nil {
		return 0, fmt.Errorf("获取配置失败: %w", err)
	}

	// 如果未启用AI检测，返回通过
	if !botConfig.AIEnabled {
		return 0, nil
	}

	// 获取AI配置
	var endpoint, apiKey, model string
	if botConfig.UseCustomAI && botConfig.CustomAIEndpoint != "" {
		endpoint = botConfig.CustomAIEndpoint
		apiKey = botConfig.CustomAIKey
		model = botConfig.CustomAIModel
	} else {
		// 使用默认配置
		endpoint = s.config.DefaultAIEndpoint
		apiKey = s.config.DefaultAIKey
		model = s.config.DefaultAIModel
	}

	if endpoint == "" || apiKey == "" || model == "" {
		return 0, fmt.Errorf("AI配置不完整")
	}

	// 创建AI服务
	aiService := NewAIService(endpoint, apiKey, model)

	// 调用AI检测
	score, err := aiService.CheckSpam(message)
	if err != nil {
		return 0, fmt.Errorf("AI检测失败: %w", err)
	}

	return score, nil
}

// AIService AI服务
type AIService struct {
	endpoint string
	apiKey   string
	model    string
}

// NewAIService 创建AI服务
func NewAIService(endpoint, apiKey, model string) *AIService {
	return &AIService{
		endpoint: endpoint,
		apiKey:   apiKey,
		model:    model,
	}
}

// CheckSpam 检查消息是否为垃圾
// 返回值：0-100，分数越高越可能是垃圾消息
func (s *AIService) CheckSpam(message string) (int, error) {
	// 这里实现具体的AI检测逻辑
	// 简化版本：使用关键词匹配
	score := s.checkKeywords(message)
	return score, nil
}

// checkKeywords 关键词检查
func (s *AIService) checkKeywords(message string) int {
	spamKeywords := []string{
		"WIN!", "CONGRATULATIONS", "FREE", "CLICK HERE",
		"URGENT", "WARNING", "UNSUBSCRIBE", "OPT OUT",
		"VIAGRA", "Cialis", "LOTTERY", "PRIZE",
		"CREDIT CARD", "DEBIT CARD", "FINANCE",
	}

	score := 0
	messageUpper := message
	for _, keyword := range spamKeywords {
		if contains(messageUpper, keyword) {
			score += 10
		}
	}

	return score
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsIgnoreCase(s, substr))
}

func containsIgnoreCase(s, substr string) bool {
	sLower := toLower(s)
	substrLower := toLower(substr)
	return len(sLower) >= len(substrLower) && findSubstring(sLower, substrLower)
}

func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + ('a' - 'A'))
		} else {
			result += string(r)
		}
	}
	return result
}

// GetConfig 获取AI配置
func (s *AISpamChecker) GetConfig() (bool, int, error) {
	configRepo := database.NewBotConfigRepository(s.db)
	botConfig, err := configRepo.GetByBotID(s.botID)
	if err != nil {
		return false, 0, err
	}

	return botConfig.AIEnabled, botConfig.WhitelistThreshold, nil
}

// EnableAI 启用AI检测
func (s *AISpamChecker) EnableAI() error {
	configRepo := database.NewBotConfigRepository(s.db)
	return configRepo.SetAIEnabled(s.botID, true)
}

// DisableAI 禁用AI检测
func (s *AISpamChecker) DisableAI() error {
	configRepo := database.NewBotConfigRepository(s.db)
	return configRepo.SetAIEnabled(s.botID, false)
}

// SetThreshold 设置白名单阈值
func (s *AISpamChecker) SetThreshold(threshold int) error {
	configRepo := database.NewBotConfigRepository(s.db)
	return configRepo.SetThreshold(s.botID, threshold)
}