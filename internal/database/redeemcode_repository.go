package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/acnoway/tg_forward_master/internal/models"
)

// RedeemCodeRepository 兑换码数据访问层
type RedeemCodeRepository struct {
	db *DB
}

// NewRedeemCodeRepository 创建兑换码仓库
func NewRedeemCodeRepository(db *DB) *RedeemCodeRepository {
	return &RedeemCodeRepository{db: db}
}

// Create 创建兑换码
func (r *RedeemCodeRepository) Create(code *models.RedeemCode) error {
	query := `
		INSERT INTO redeem_codes (code, plan_id, status)
		VALUES (?, ?, ?)
	`
	result, err := r.db.Exec(query, code.Code, code.PlanID, code.Status)
	if err != nil {
		return fmt.Errorf("创建兑换码失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取兑换码ID失败: %w", err)
	}

	code.ID = id
	code.CreatedAt = time.Now()
	return nil
}

// GetByCode 根据兑换码获取记录
func (r *RedeemCodeRepository) GetByCode(codeStr string) (*models.RedeemCode, error) {
	query := `
		SELECT id, code, plan_id, status, used_by, created_at, used_at
		FROM redeem_codes
		WHERE code = ?
	`
	code := &models.RedeemCode{}
	err := r.db.QueryRow(query, codeStr).Scan(
		&code.ID,
		&code.Code,
		&code.PlanID,
		&code.Status,
		&code.UsedBy,
		&code.CreatedAt,
		&code.UsedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询兑换码失败: %w", err)
	}
	return code, nil
}

// UseCode 使用兑换码
func (r *RedeemCodeRepository) UseCode(id, userID int64) error {
	now := time.Now()
	query := `
		UPDATE redeem_codes
		SET status = 'used', used_by = ?, used_at = ?
		WHERE id = ? AND status = 'unused'
	`
	result, err := r.db.Exec(query, userID, now, id)
	if err != nil {
		return fmt.Errorf("使用兑换码失败: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("兑换码已被使用或不存在")
	}

	return nil
}

// GetUnused 获取所有未使用的兑换码
func (r *RedeemCodeRepository) GetUnused() ([]*models.RedeemCode, error) {
	query := `
		SELECT id, code, plan_id, status, used_by, created_at, used_at
		FROM redeem_codes
		WHERE status = 'unused'
		ORDER BY created_at DESC
	`
	return r.queryList(query)
}

// GetByPlanID 获取某个套餐的所有兑换码
func (r *RedeemCodeRepository) GetByPlanID(planID int64) ([]*models.RedeemCode, error) {
	query := `
		SELECT id, code, plan_id, status, used_by, created_at, used_at
		FROM redeem_codes
		WHERE plan_id = ?
		ORDER BY created_at DESC
	`
	return r.queryList(query, planID)
}

// GetAll 获取所有兑换码
func (r *RedeemCodeRepository) GetAll(limit, offset int) ([]*models.RedeemCode, error) {
	query := `
		SELECT id, code, plan_id, status, used_by, created_at, used_at
		FROM redeem_codes
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	return r.queryList(query, limit, offset)
}

// queryList 查询列表辅助函数
func (r *RedeemCodeRepository) queryList(query string, args ...interface{}) ([]*models.RedeemCode, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询兑换码列表失败: %w", err)
	}
	defer rows.Close()

	var codes []*models.RedeemCode
	for rows.Next() {
		code := &models.RedeemCode{}
		if err := rows.Scan(
			&code.ID,
			&code.Code,
			&code.PlanID,
			&code.Status,
			&code.UsedBy,
			&code.CreatedAt,
			&code.UsedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描兑换码数据失败: %w", err)
		}
		codes = append(codes, code)
	}
	return codes, nil
}

// Count 获取兑换码总数
func (r *RedeemCodeRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM redeem_codes`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询兑换码数量失败: %w", err)
	}
	return count, nil
}

// CountUnused 获取未使用的兑换码数量
func (r *RedeemCodeRepository) CountUnused() (int, error) {
	query := `SELECT COUNT(*) FROM redeem_codes WHERE status = 'unused'`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询未使用兑换码数量失败: %w", err)
	}
	return count, nil
}

// Delete 删除兑换码
func (r *RedeemCodeRepository) Delete(id int64) error {
	query := `DELETE FROM redeem_codes WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("删除兑换码失败: %w", err)
	}
	return nil
}
