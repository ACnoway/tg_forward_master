package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/acnoway/tg_forward_master/internal/models"
)

// SubscriptionRepository 订阅数据访问层
type SubscriptionRepository struct {
	db *DB
}

// NewSubscriptionRepository 创建订阅仓库
func NewSubscriptionRepository(db *DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create 创建订阅
func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, plan_id, status, expires_at, max_bots)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, sub.UserID, sub.PlanID, sub.Status, sub.ExpiresAt, sub.MaxBots)
	if err != nil {
		return fmt.Errorf("创建订阅失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取订阅ID失败: %w", err)
	}

	sub.ID = id
	sub.CreatedAt = time.Now()
	return nil
}

// GetActiveByUserID 获取用户的有效订阅
func (r *SubscriptionRepository) GetActiveByUserID(userID int64) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, expires_at, max_bots, created_at
		FROM subscriptions
		WHERE user_id = ? AND status = 'active' AND expires_at > datetime('now')
		ORDER BY expires_at DESC
		LIMIT 1
	`
	sub := &models.Subscription{}
	err := r.db.QueryRow(query, userID).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.PlanID,
		&sub.Status,
		&sub.ExpiresAt,
		&sub.MaxBots,
		&sub.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询订阅失败: %w", err)
	}
	return sub, nil
}

// GetByID 根据ID获取订阅
func (r *SubscriptionRepository) GetByID(id int64) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, expires_at, max_bots, created_at
		FROM subscriptions
		WHERE id = ?
	`
	sub := &models.Subscription{}
	err := r.db.QueryRow(query, id).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.PlanID,
		&sub.Status,
		&sub.ExpiresAt,
		&sub.MaxBots,
		&sub.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询订阅失败: %w", err)
	}
	return sub, nil
}

// UpdateStatus 更新订阅状态
func (r *SubscriptionRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE subscriptions SET status = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("更新订阅状态失败: %w", err)
	}
	return nil
}

// ExtendExpiration 延长订阅
func (r *SubscriptionRepository) ExtendExpiration(id int64, days int) error {
	query := `
		UPDATE subscriptions
		SET expires_at = datetime(expires_at, '+' || ? || ' days')
		WHERE id = ?
	`
	_, err := r.db.Exec(query, days, id)
	if err != nil {
		return fmt.Errorf("延长订阅失败: %w", err)
	}
	return nil
}

// ExpireOldSubscriptions 使过期订阅失效
func (r *SubscriptionRepository) ExpireOldSubscriptions() (int, error) {
	query := `
		UPDATE subscriptions
		SET status = 'expired'
		WHERE status = 'active' AND expires_at < datetime('now')
	`
	result, err := r.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("更新过期订阅失败: %w", err)
	}

	affected, _ := result.RowsAffected()
	return int(affected), nil
}

// GetAllByUserID 获取用户所有订阅记录
func (r *SubscriptionRepository) GetAllByUserID(userID int64) ([]*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, expires_at, max_bots, created_at
		FROM subscriptions
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询订阅列表失败: %w", err)
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		if err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.PlanID,
			&sub.Status,
			&sub.ExpiresAt,
			&sub.MaxBots,
			&sub.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描订阅数据失败: %w", err)
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

// Count 获取总订阅数
func (r *SubscriptionRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM subscriptions`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询订阅数量失败: %w", err)
	}
	return count, nil
}

// CountActive 获取有效订阅数
func (r *SubscriptionRepository) CountActive() (int, error) {
	query := `SELECT COUNT(*) FROM subscriptions WHERE status = 'active' AND expires_at > datetime('now')`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询有效订阅数量失败: %w", err)
	}
	return count, nil
}
