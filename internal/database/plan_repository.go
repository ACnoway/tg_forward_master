package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/acnoway/tg_forward_master/internal/models"
)

// PlanRepository 套餐数据访问层
type PlanRepository struct {
	db *DB
}

// NewPlanRepository 创建套餐仓库
func NewPlanRepository(db *DB) *PlanRepository {
	return &PlanRepository{db: db}
}

// Create 创建套餐
func (r *PlanRepository) Create(plan *models.Plan) error {
	query := `
		INSERT INTO plans (name, price, duration_days, max_bots, is_active)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, plan.Name, plan.Price, plan.DurationDays, plan.MaxBots, plan.IsActive)
	if err != nil {
		return fmt.Errorf("创建套餐失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取套餐ID失败: %w", err)
	}

	plan.ID = id
	plan.CreatedAt = time.Now()
	return nil
}

// GetByID 根据ID获取套餐
func (r *PlanRepository) GetByID(id int64) (*models.Plan, error) {
	query := `
		SELECT id, name, price, duration_days, max_bots, is_active, created_at
		FROM plans
		WHERE id = ?
	`
	plan := &models.Plan{}
	err := r.db.QueryRow(query, id).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Price,
		&plan.DurationDays,
		&plan.MaxBots,
		&plan.IsActive,
		&plan.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询套餐失败: %w", err)
	}
	return plan, nil
}

// GetAll 获取所有套餐
func (r *PlanRepository) GetAll() ([]*models.Plan, error) {
	query := `
		SELECT id, name, price, duration_days, max_bots, is_active, created_at
		FROM plans
		ORDER BY price ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询套餐列表失败: %w", err)
	}
	defer rows.Close()

	var plans []*models.Plan
	for rows.Next() {
		plan := &models.Plan{}
		if err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Price,
			&plan.DurationDays,
			&plan.MaxBots,
			&plan.IsActive,
			&plan.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描套餐数据失败: %w", err)
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

// GetActive 获取所有有效套餐
func (r *PlanRepository) GetActive() ([]*models.Plan, error) {
	query := `
		SELECT id, name, price, duration_days, max_bots, is_active, created_at
		FROM plans
		WHERE is_active = 1
		ORDER BY price ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询有效套餐失败: %w", err)
	}
	defer rows.Close()

	var plans []*models.Plan
	for rows.Next() {
		plan := &models.Plan{}
		if err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Price,
			&plan.DurationDays,
			&plan.MaxBots,
			&plan.IsActive,
			&plan.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描套餐数据失败: %w", err)
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

// Update 更新套餐
func (r *PlanRepository) Update(plan *models.Plan) error {
	query := `
		UPDATE plans
		SET name = ?, price = ?, duration_days = ?, max_bots = ?, is_active = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, plan.Name, plan.Price, plan.DurationDays, plan.MaxBots, plan.IsActive, plan.ID)
	if err != nil {
		return fmt.Errorf("更新套餐失败: %w", err)
	}
	return nil
}

// Delete 删除套餐（实际是标记为不可用）
func (r *PlanRepository) Delete(id int64) error {
	query := `UPDATE plans SET is_active = 0 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("删除套餐失败: %w", err)
	}
	return nil
}

// Count 获取套餐总数
func (r *PlanRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM plans WHERE is_active = 1`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询套餐数量失败: %w", err)
	}
	return count, nil
}
