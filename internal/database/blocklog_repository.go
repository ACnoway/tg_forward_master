package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/acnoway/tg_forward_master/internal/models"
)

// BlockLogRepository 拦截日志数据访问层
type BlockLogRepository struct {
	db *DB
}

// NewBlockLogRepository 创建拦截日志仓库
func NewBlockLogRepository(db *DB) *BlockLogRepository {
	return &BlockLogRepository{db: db}
}

// Create 创建拦截日志
func (r *BlockLogRepository) Create(log *models.BlockLog) error {
	query := `
		INSERT INTO block_logs (bot_id, telegram_id, username, first_name, last_name,
		                        block_reason, ai_confidence, message_content, message_type)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		log.BotID,
		log.TelegramID,
		log.Username,
		log.FirstName,
		log.LastName,
		log.BlockReason,
		log.AIConfidence,
		log.MessageContent,
		log.MessageType,
	)

	if err != nil {
		return fmt.Errorf("创建拦截日志失败: %w", err)
	}

	log.BlockedAt = time.Now()
	return nil
}

// GetByID 根据ID获取拦截日志
func (r *BlockLogRepository) GetByID(id int64) (*models.BlockLog, error) {
	query := `
		SELECT id, bot_id, telegram_id, username, first_name, last_name,
		       block_reason, ai_confidence, message_content, message_type,
		       blocked_at, is_false_positive
		FROM block_logs
		WHERE id = ?
	`

	log := &models.BlockLog{}
	err := r.db.QueryRow(query, id).Scan(
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
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("查询拦截日志失败: %w", err)
	}

	return log, nil
}

// GetByBotID 获取Bot的所有拦截日志
func (r *BlockLogRepository) GetByBotID(botID int64) ([]*models.BlockLog, error) {
	query := `
		SELECT id, bot_id, telegram_id, username, first_name, last_name,
		       block_reason, ai_confidence, message_content, message_type,
		       blocked_at, is_false_positive
		FROM block_logs
		WHERE bot_id = ?
		ORDER BY blocked_at DESC
	`

	rows, err := r.db.Query(query, botID)
	if err != nil {
		return nil, fmt.Errorf("查询拦截日志列表失败: %w", err)
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

// GetRecent 获取最近的拦截日志
func (r *BlockLogRepository) GetRecent(botID int64, limit int) ([]*models.BlockLog, error) {
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
		return nil, fmt.Errorf("查询最近拦截日志失败: %w", err)
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

// GetByTelegramID 获取特定用户的拦截日志
func (r *BlockLogRepository) GetByTelegramID(botID, telegramID int64) ([]*models.BlockLog, error) {
	query := `
		SELECT id, bot_id, telegram_id, username, first_name, last_name,
		       block_reason, ai_confidence, message_content, message_type,
		       blocked_at, is_false_positive
		FROM block_logs
		WHERE bot_id = ? AND telegram_id = ?
		ORDER BY blocked_at DESC
	`

	rows, err := r.db.Query(query, botID, telegramID)
	if err != nil {
		return nil, fmt.Errorf("查询用户拦截日志失败: %w", err)
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

// MarkAsFalsePositive 标记为误判
func (r *BlockLogRepository) MarkAsFalsePositive(id int64) error {
	query := `UPDATE block_logs SET is_false_positive = 1 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("标记为误判失败: %w", err)
	}
	return nil
}

// Count 获取拦截日志总数
func (r *BlockLogRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM block_logs`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询拦截日志数量失败: %w", err)
	}
	return count, nil
}

// CountByBotID 获取Bot的拦截日志数量
func (r *BlockLogRepository) CountByBotID(botID int64) (int, error) {
	query := `SELECT COUNT(*) FROM block_logs WHERE bot_id = ?`
	var count int
	err := r.db.QueryRow(query, botID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询Bot拦截日志数量失败: %w", err)
	}
	return count, nil
}

// CountByReason 按原因统计拦截日志
func (r *BlockLogRepository) CountByReason(botID int64) (map[string]int, error) {
	query := `
		SELECT block_reason, COUNT(*) as count
		FROM block_logs
		WHERE bot_id = ?
		GROUP BY block_reason
	`

	rows, err := r.db.Query(query, botID)
	if err != nil {
		return nil, fmt.Errorf("按原因统计拦截日志失败: %w", err)
	}
	defer rows.Close()

	reasons := make(map[string]int)
	for rows.Next() {
		var reason string
		var count int
		if err := rows.Scan(&reason, &count); err != nil {
			return nil, fmt.Errorf("扫描统计结果失败: %w", err)
		}
		reasons[reason] = count
	}

	return reasons, rows.Err()
}

// DeleteOldLogs 删除旧的拦截日志
func (r *BlockLogRepository) DeleteOldLogs(days int) (int, error) {
	query := `DELETE FROM block_logs WHERE blocked_at < datetime('now', '-' || ? || ' days')`
	result, err := r.db.Exec(query, days)
	if err != nil {
		return 0, fmt.Errorf("删除旧拦截日志失败: %w", err)
	}

	affected, _ := result.RowsAffected()
	return int(affected), nil
}