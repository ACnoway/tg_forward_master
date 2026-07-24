package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/acnoway/tg_forward_master/internal/models"
)

// BotConfigRepository Bot配置数据访问层
type BotConfigRepository struct {
	db *DB
}

// NewBotConfigRepository 创建Bot配置仓库
func NewBotConfigRepository(db *DB) *BotConfigRepository {
	return &BotConfigRepository{db: db}
}

// GetByBotID 根据Bot ID获取配置
func (r *BotConfigRepository) GetByBotID(botID int64) (*models.BotConfig, error) {
	query := `
		SELECT bot_id, owner_telegram_id, ai_enabled, use_custom_ai,
		       custom_ai_endpoint, custom_ai_key, custom_ai_model,
		       whitelist_threshold, notify_on_block, notify_on_ai_block,
		       created_at, updated_at
		FROM bot_configs
		WHERE bot_id = ?
	`

	config := &models.BotConfig{}
	err := r.db.QueryRow(query, botID).Scan(
		&config.BotID,
		&config.OwnerTelegramID,
		&config.AIEnabled,
		&config.UseCustomAI,
		&config.CustomAIEndpoint,
		&config.CustomAIKey,
		&config.CustomAIModel,
		&config.WhitelistThreshold,
		&config.NotifyOnBlock,
		&config.NotifyOnAIBlock,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// 如果不存在，返回默认配置
		return &models.BotConfig{
			BotID:              botID,
			OwnerTelegramID:    0,
			AIEnabled:          true,
			UseCustomAI:        false,
			CustomAIEndpoint:   "",
			CustomAIKey:        "",
			CustomAIModel:      "",
			WhitelistThreshold: 3,
			NotifyOnBlock:      true,
			NotifyOnAIBlock:    false,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("查询Bot配置失败: %w", err)
	}

	return config, nil
}

// Create 创建Bot配置
func (r *BotConfigRepository) Create(config *models.BotConfig) error {
	query := `
		INSERT INTO bot_configs (bot_id, owner_telegram_id, ai_enabled, use_custom_ai,
		                         custom_ai_endpoint, custom_ai_key, custom_ai_model,
		                         whitelist_threshold, notify_on_block, notify_on_ai_block)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		config.BotID,
		config.OwnerTelegramID,
		config.AIEnabled,
		config.UseCustomAI,
		config.CustomAIEndpoint,
		config.CustomAIKey,
		config.CustomAIModel,
		config.WhitelistThreshold,
		config.NotifyOnBlock,
		config.NotifyOnAIBlock,
	)

	if err != nil {
		return fmt.Errorf("创建Bot配置失败: %w", err)
	}

	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	return nil
}

// Update 更新Bot配置
func (r *BotConfigRepository) Update(config *models.BotConfig) error {
	query := `
		UPDATE bot_configs
		SET ai_enabled = ?, use_custom_ai = ?,
		    custom_ai_endpoint = ?, custom_ai_key = ?, custom_ai_model = ?,
		    whitelist_threshold = ?, notify_on_block = ?, notify_on_ai_block = ?,
		    updated_at = ?
		WHERE bot_id = ?
	`

	_, err := r.db.Exec(query,
		config.AIEnabled,
		config.UseCustomAI,
		config.CustomAIEndpoint,
		config.CustomAIKey,
		config.CustomAIModel,
		config.WhitelistThreshold,
		config.NotifyOnBlock,
		config.NotifyOnAIBlock,
		time.Now(),
		config.BotID,
	)

	if err != nil {
		return fmt.Errorf("更新Bot配置失败: %w", err)
	}

	config.UpdatedAt = time.Now()
	return nil
}

// SetAIEnabled 设置AI启用状态
func (r *BotConfigRepository) SetAIEnabled(botID int64, enabled bool) error {
	query := `UPDATE bot_configs SET ai_enabled = ?, updated_at = ? WHERE bot_id = ?`
	_, err := r.db.Exec(query, enabled, time.Now(), botID)
	if err != nil {
		return fmt.Errorf("更新AI启用状态失败: %w", err)
	}
	return nil
}

// SetThreshold 设置白名单阈值
func (r *BotConfigRepository) SetThreshold(botID int64, threshold int) error {
	query := `UPDATE bot_configs SET whitelist_threshold = ?, updated_at = ? WHERE bot_id = ?`
	_, err := r.db.Exec(query, threshold, time.Now(), botID)
	if err != nil {
		return fmt.Errorf("更新白名单阈值失败: %w", err)
	}
	return nil
}

// SetCustomAI 设置自定义AI配置
func (r *BotConfigRepository) SetCustomAI(botID int64, endpoint, key, model string) error {
	query := `
		UPDATE bot_configs
		SET use_custom_ai = ?, custom_ai_endpoint = ?, custom_ai_key = ?, custom_ai_model = ?, updated_at = ?
		WHERE bot_id = ?
	`
	_, err := r.db.Exec(query, true, endpoint, key, model, time.Now(), botID)
	if err != nil {
		return fmt.Errorf("更新自定义AI配置失败: %w", err)
	}
	return nil
}

// SetDefaultAI 设置使用默认AI配置
func (r *BotConfigRepository) SetDefaultAI(botID int64) error {
	query := `UPDATE bot_configs SET use_custom_ai = ?, updated_at = ? WHERE bot_id = ?`
	_, err := r.db.Exec(query, false, time.Now(), botID)
	if err != nil {
		return fmt.Errorf("更新AI配置模式失败: %w", err)
	}
	return nil
}

// GetBlacklist 获取黑名单列表
func (r *BotConfigRepository) GetBlacklist(botID int64) ([]*models.Customer, error) {
	query := `
		SELECT id, bot_id, telegram_id, username, first_name, last_name,
		       verified_count, total_messages, is_whitelisted, is_blacklisted,
		       first_seen_at, last_message_at, whitelisted_at, blacklisted_at,
		       blacklist_reason, is_manual_trust
		FROM customers
		WHERE bot_id = ? AND is_blacklisted = 1
		ORDER BY blacklisted_at DESC
	`

	rows, err := r.db.Query(query, botID)
	if err != nil {
		return nil, fmt.Errorf("查询黑名单列表失败: %w", err)
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		customer := &models.Customer{}
		if err := rows.Scan(
			&customer.ID,
			&customer.BotID,
			&customer.TelegramID,
			&customer.Username,
			&customer.FirstName,
			&customer.LastName,
			&customer.VerifiedCount,
			&customer.TotalMessages,
			&customer.IsWhitelisted,
			&customer.IsBlacklisted,
			&customer.FirstSeenAt,
			&customer.LastMessageAt,
			&customer.WhitelistedAt,
			&customer.BlacklistedAt,
			&customer.BlacklistReason,
			&customer.IsManualTrust,
		); err != nil {
			return nil, fmt.Errorf("扫描黑名单数据失败: %w", err)
		}
		customers = append(customers, customer)
	}

	return customers, rows.Err()
}

// GetWhitelist 获取白名单列表
func (r *BotConfigRepository) GetWhitelist(botID int64) ([]*models.Customer, error) {
	query := `
		SELECT id, bot_id, telegram_id, username, first_name, last_name,
		       verified_count, total_messages, is_whitelisted, is_blacklisted,
		       first_seen_at, last_message_at, whitelisted_at, blacklisted_at,
		       blacklist_reason, is_manual_trust
		FROM customers
		WHERE bot_id = ? AND is_whitelisted = 1
		ORDER BY whitelisted_at DESC
	`

	rows, err := r.db.Query(query, botID)
	if err != nil {
		return nil, fmt.Errorf("查询白名单列表失败: %w", err)
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		customer := &models.Customer{}
		if err := rows.Scan(
			&customer.ID,
			&customer.BotID,
			&customer.TelegramID,
			&customer.Username,
			&customer.FirstName,
			&customer.LastName,
			&customer.VerifiedCount,
			&customer.TotalMessages,
			&customer.IsWhitelisted,
			&customer.IsBlacklisted,
			&customer.FirstSeenAt,
			&customer.LastMessageAt,
			&customer.WhitelistedAt,
			&customer.BlacklistedAt,
			&customer.BlacklistReason,
			&customer.IsManualTrust,
		); err != nil {
			return nil, fmt.Errorf("扫描白名单数据失败: %w", err)
		}
		customers = append(customers, customer)
	}

	return customers, rows.Err()
}

// GetBlockLogs 获取拦截日志
func (r *BotConfigRepository) GetBlockLogs(botID int64, limit int) ([]*models.BlockLog, error) {
	query := `
		SELECT id, bot_id, telegram_id, username, first_name, last_name,
		       block_reason, ai_confidence, message_content, message_type,
		       blocked_at, is_false_positive
		FROM block_logs
		WHERE bot_id = ?
		ORDER BY blocked_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, botID, limit)
	if err != nil {
		return nil, fmt.Errorf("查询拦截日志失败: %w", err)
	}
	defer rows.Close()

	var logs []*models.BlockLog
	for rows.Next() {
		log := &models.BlockLog{}
		if err := rows.Scan(
			&log.ID,
			&log.BotID,
			&log.TelegramID,
			&log.Username,
			&log.FirstName,
			&log.LastName,
			&log.BlockReason,
			&log.AIConfidence,
			&log.MessageContent,
			&log.MessageType,
			&log.BlockedAt,
			&log.IsFalsePositive,
		); err != nil {
			return nil, fmt.Errorf("扫描拦截日志数据失败: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// CountBlockLogs 统计拦截日志
func (r *BotConfigRepository) CountBlockLogs(botID int64) (int, error) {
	query := `SELECT COUNT(*) FROM block_logs WHERE bot_id = ?`
	var count int
	err := r.db.QueryRow(query, botID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询拦截日志数量失败: %w", err)
	}
	return count, nil
}

// CountCustomers 统计客户数量
func (r *BotConfigRepository) CountCustomers(botID int64) (int, error) {
	query := `SELECT COUNT(*) FROM customers WHERE bot_id = ?`
	var count int
	err := r.db.QueryRow(query, botID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询客户数量失败: %w", err)
	}
	return count, nil
}

// CountWhitelisted 统计白名单客户数量
func (r *BotConfigRepository) CountWhitelisted(botID int64) (int, error) {
	query := `SELECT COUNT(*) FROM customers WHERE bot_id = ? AND is_whitelisted = 1`
	var count int
	err := r.db.QueryRow(query, botID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询白名单客户数量失败: %w", err)
	}
	return count, nil
}

// CountBlacklisted 统计黑名单客户数量
func (r *BotConfigRepository) CountBlacklisted(botID int64) (int, error) {
	query := `SELECT COUNT(*) FROM customers WHERE bot_id = ? AND is_blacklisted = 1`
	var count int
	err := r.db.QueryRow(query, botID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询黑名单客户数量失败: %w", err)
	}
	return count, nil
}

// GetStats 获取统计信息
func (r *BotConfigRepository) GetStats(botID int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 客户统计
	totalCustomers, err := r.CountCustomers(botID)
	if err != nil {
		return nil, err
	}
	stats["total_customers"] = totalCustomers

	whitelisted, err := r.CountWhitelisted(botID)
	if err != nil {
		return nil, err
	}
	stats["whitelisted_customers"] = whitelisted

	blacklisted, err := r.CountBlacklisted(botID)
	if err != nil {
		return nil, err
	}
	stats["blacklisted_customers"] = blacklisted

	// 拦截统计
	blockLogs, err := r.CountBlockLogs(botID)
	if err != nil {
		return nil, err
	}
	stats["total_blocks"] = blockLogs

	return stats, nil
}