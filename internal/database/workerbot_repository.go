package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/acnoway/tg_forward_master/internal/models"
)

// WorkerBotRepository 子Bot数据访问层
type WorkerBotRepository struct {
	db *DB
}

// NewWorkerBotRepository 创建子Bot仓库
func NewWorkerBotRepository(db *DB) *WorkerBotRepository {
	return &WorkerBotRepository{db: db}
}

// Create 创建子Bot
func (r *WorkerBotRepository) Create(bot *models.WorkerBot) error {
	query := `
		INSERT INTO worker_bots (user_id, bot_token, bot_username, bot_nickname, owner_telegram_id, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, bot.UserID, bot.BotToken, bot.BotUsername, bot.BotNickname, bot.OwnerTelegramID, bot.Status)
	if err != nil {
		return fmt.Errorf("创建子Bot失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取子Bot ID失败: %w", err)
	}

	bot.ID = id
	bot.CreatedAt = time.Now()
	return nil
}

// GetByID 根据ID获取子Bot
func (r *WorkerBotRepository) GetByID(id int64) (*models.WorkerBot, error) {
	query := `
		SELECT id, user_id, bot_token, bot_username, bot_nickname, owner_telegram_id, status, created_at
		FROM worker_bots
		WHERE id = ?
	`
	bot := &models.WorkerBot{}
	err := r.db.QueryRow(query, id).Scan(
		&bot.ID,
		&bot.UserID,
		&bot.BotToken,
		&bot.BotUsername,
		&bot.BotNickname,
		&bot.OwnerTelegramID,
		&bot.Status,
		&bot.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询子Bot失败: %w", err)
	}
	return bot, nil
}

// GetByUserID 获取用户的所有子Bot
func (r *WorkerBotRepository) GetByUserID(userID int64) ([]*models.WorkerBot, error) {
	query := `
		SELECT id, user_id, bot_token, bot_username, bot_nickname, owner_telegram_id, status, created_at
		FROM worker_bots
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询子Bot列表失败: %w", err)
	}
	defer rows.Close()

	var bots []*models.WorkerBot
	for rows.Next() {
		bot := &models.WorkerBot{}
		if err := rows.Scan(
			&bot.ID,
			&bot.UserID,
			&bot.BotToken,
			&bot.BotUsername,
			&bot.BotNickname,
			&bot.OwnerTelegramID,
			&bot.Status,
			&bot.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描子Bot数据失败: %w", err)
		}
		bots = append(bots, bot)
	}
	return bots, nil
}

// Update 更新子Bot信息
func (r *WorkerBotRepository) Update(bot *models.WorkerBot) error {
	query := `
		UPDATE worker_bots
		SET bot_token = ?, bot_username = ?, bot_nickname = ?, owner_telegram_id = ?, status = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, bot.BotToken, bot.BotUsername, bot.BotNickname, bot.OwnerTelegramID, bot.Status, bot.ID)
	if err != nil {
		return fmt.Errorf("更新子Bot失败: %w", err)
	}
	return nil
}

// UpdateStatus 更新子Bot状态
func (r *WorkerBotRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE worker_bots SET status = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("更新子Bot状态失败: %w", err)
	}
	return nil
}

// Delete 删除子Bot
func (r *WorkerBotRepository) Delete(id int64) error {
	// 先删除相关配置
	query := `DELETE FROM bot_configs WHERE bot_id = ?`
	if _, err := r.db.Exec(query, id); err != nil {
		return fmt.Errorf("删除子Bot配置失败: %w", err)
	}

	// 再删除Bot记录
	query = `DELETE FROM worker_bots WHERE id = ?`
	if _, err := r.db.Exec(query, id); err != nil {
		return fmt.Errorf("删除子Bot失败: %w", err)
	}

	return nil
}

// CountByUserID 获取用户的子Bot数量
func (r *WorkerBotRepository) CountByUserID(userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM worker_bots WHERE user_id = ?`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询子Bot数量失败: %w", err)
	}
	return count, nil
}

// Count 获取所有子Bot数量
func (r *WorkerBotRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM worker_bots`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询子Bot总数失败: %w", err)
	}
	return count, nil
}

// CountRunning 获取运行中的子Bot数量
func (r *WorkerBotRepository) CountRunning() (int, error) {
	query := `SELECT COUNT(*) FROM worker_bots WHERE status = 'running'`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询运行中子Bot数量失败: %w", err)
	}
	return count, nil
}
