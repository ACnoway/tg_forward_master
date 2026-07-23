package database

import (
	"database/sql"
	"time"

	"github.com/acnoway/tg_forward_master/internal/models"
)

// CustomerRepository 客户数据仓库
type CustomerRepository struct {
	db *DB
}

// NewCustomerRepository 创建客户仓库
func NewCustomerRepository(db *DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// GetOrCreate 获取或创建客户记录
func (r *CustomerRepository) GetOrCreate(botID, telegramID int64, username, firstName, lastName string) (*models.Customer, error) {
	// 先尝试获取
	customer, err := r.GetByTelegramID(botID, telegramID)
	if err == nil {
		return customer, nil
	}

	// 如果不存在，创建新记录
	if err == sql.ErrNoRows {
		return r.Create(botID, telegramID, username, firstName, lastName)
	}

	return nil, err
}

// GetByTelegramID 通过Telegram ID获取客户
func (r *CustomerRepository) GetByTelegramID(botID, telegramID int64) (*models.Customer, error) {
	query := `
		SELECT id, bot_id, telegram_id, username, first_name, last_name,
		       verified_count, total_messages, is_whitelisted, is_blacklisted,
		       first_seen_at, last_message_at, whitelisted_at, blacklisted_at,
		       blacklist_reason, is_manual_trust
		FROM customers
		WHERE bot_id = ? AND telegram_id = ?
	`

	customer := &models.Customer{}
	err := r.db.QueryRow(query, botID, telegramID).Scan(
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
	)

	if err != nil {
		return nil, err
	}

	return customer, nil
}

// Create 创建新客户
func (r *CustomerRepository) Create(botID, telegramID int64, username, firstName, lastName string) (*models.Customer, error) {
	query := `
		INSERT INTO customers (bot_id, telegram_id, username, first_name, last_name, first_seen_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(query, botID, telegramID, username, firstName, lastName, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Customer{
		ID:         id,
		BotID:      botID,
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		FirstSeenAt: now,
	}, nil
}

// Update 更新客户信息
func (r *CustomerRepository) Update(customer *models.Customer) error {
	query := `
		UPDATE customers
		SET username = ?, first_name = ?, last_name = ?,
		    verified_count = ?, total_messages = ?,
		    is_whitelisted = ?, is_blacklisted = ?,
		    last_message_at = ?, whitelisted_at = ?,
		    blacklisted_at = ?, blacklist_reason = ?,
		    is_manual_trust = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		customer.Username,
		customer.FirstName,
		customer.LastName,
		customer.VerifiedCount,
		customer.TotalMessages,
		customer.IsWhitelisted,
		customer.IsBlacklisted,
		customer.LastMessageAt,
		customer.WhitelistedAt,
		customer.BlacklistedAt,
		customer.BlacklistReason,
		customer.IsManualTrust,
		customer.ID,
	)

	return err
}

// IncrementMessageCount 增加消息计数
func (r *CustomerRepository) IncrementMessageCount(botID, telegramID int64) error {
	query := `
		UPDATE customers
		SET total_messages = total_messages + 1,
		    last_message_at = ?
		WHERE bot_id = ? AND telegram_id = ?
	`

	_, err := r.db.Exec(query, time.Now(), botID, telegramID)
	return err
}

// IncrementVerifiedCount 增加已验证消息计数
func (r *CustomerRepository) IncrementVerifiedCount(botID, telegramID int64, threshold int) error {
	query := `
		UPDATE customers
		SET verified_count = verified_count + 1,
		    is_whitelisted = CASE
		        WHEN verified_count + 1 >= ? THEN 1
		        ELSE is_whitelisted
		    END,
		    whitelisted_at = CASE
		        WHEN verified_count + 1 >= ? AND whitelisted_at IS NULL THEN ?
		        ELSE whitelisted_at
		    END
		WHERE bot_id = ? AND telegram_id = ?
	`

	now := time.Now()
	_, err := r.db.Exec(query, threshold, threshold, now, botID, telegramID)
	return err
}

// Blacklist 拉黑客户
func (r *CustomerRepository) Blacklist(botID, telegramID int64, reason string) error {
	query := `
		UPDATE customers
		SET is_blacklisted = 1,
		    blacklisted_at = ?,
		    blacklist_reason = ?
		WHERE bot_id = ? AND telegram_id = ?
	`

	_, err := r.db.Exec(query, time.Now(), reason, botID, telegramID)
	return err
}

// Whitelist 信任客户（手动白名单）
func (r *CustomerRepository) Whitelist(botID, telegramID int64) error {
	query := `
		UPDATE customers
		SET is_whitelisted = 1,
		    is_manual_trust = 1,
		    whitelisted_at = ?
		WHERE bot_id = ? AND telegram_id = ?
	`

	_, err := r.db.Exec(query, time.Now(), botID, telegramID)
	return err
}

// Unblacklist 解除拉黑
func (r *CustomerRepository) Unblacklist(botID, telegramID int64) error {
	query := `
		UPDATE customers
		SET is_blacklisted = 0,
		    blacklisted_at = NULL,
		    blacklist_reason = ''
		WHERE bot_id = ? AND telegram_id = ?
	`

	_, err := r.db.Exec(query, botID, telegramID)
	return err
}

// GetBlacklist 获取黑名单列表
func (r *CustomerRepository) GetBlacklist(botID int64) ([]*models.Customer, error) {
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
		return nil, err
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		customer := &models.Customer{}
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, rows.Err()
}

// GetWhitelist 获取白名单列表
func (r *CustomerRepository) GetWhitelist(botID int64) ([]*models.Customer, error) {
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
		return nil, err
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		customer := &models.Customer{}
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, rows.Err()
}
