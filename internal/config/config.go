package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config 系统配置
type Config struct {
	// 主控Bot配置
	MasterBotToken string

	// 数据库配置
	DatabasePath string

	// 加密密钥
	EncryptionKey string

	// 服务器配置
	ServerPort string

	// 默认AI配置
	DefaultAIEndpoint string
	DefaultAIKey      string
	DefaultAIModel    string

	// 易支付配置
	EpayAPIURL      string
	EpayMerchantID  string
	EpayMerchantKey string
}

// Load 加载配置
func Load(envPath string) (*Config, error) {
	// 加载.env文件
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %w", err)
	}

	config := &Config{
		MasterBotToken:    os.Getenv("MASTER_BOT_TOKEN"),
		DatabasePath:      os.Getenv("DATABASE_PATH"),
		EncryptionKey:     os.Getenv("ENCRYPTION_KEY"),
		ServerPort:        os.Getenv("SERVER_PORT"),
		DefaultAIEndpoint: os.Getenv("DEFAULT_AI_ENDPOINT"),
		DefaultAIKey:      os.Getenv("DEFAULT_AI_KEY"),
		DefaultAIModel:    os.Getenv("DEFAULT_AI_MODEL"),
		EpayAPIURL:        os.Getenv("EPAY_API_URL"),
		EpayMerchantID:    os.Getenv("EPAY_MERCHANT_ID"),
		EpayMerchantKey:   os.Getenv("EPAY_MERCHANT_KEY"),
	}

	// 验证必要配置
	if config.MasterBotToken == "" {
		return nil, fmt.Errorf("MASTER_BOT_TOKEN 未配置")
	}
	if config.DatabasePath == "" {
		config.DatabasePath = "./data/master.db"
	}
	if config.EncryptionKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY 未配置")
	}
	if len(config.EncryptionKey) != 32 {
		return nil, fmt.Errorf("ENCRYPTION_KEY 必须是32字节")
	}
	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}

	return config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.MasterBotToken == "" {
		return fmt.Errorf("MasterBotToken 不能为空")
	}
	if c.DatabasePath == "" {
		return fmt.Errorf("DatabasePath 不能为空")
	}
	if len(c.EncryptionKey) != 32 {
		return fmt.Errorf("EncryptionKey 必须是32字节")
	}
	return nil
}
