package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/acnoway/tg_forward_master/internal/models"
)

// SystemConfigRepository 系统配置数据访问层
type SystemConfigRepository struct {
	db *DB
}

// NewSystemConfigRepository 创建系统配置仓库
func NewSystemConfigRepository(db *DB) *SystemConfigRepository {
	return &SystemConfigRepository{db: db}
}

// Set 设置配置项
func (r *SystemConfigRepository) Set(key, value string) error {
	query := `
		INSERT INTO system_configs (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = ?
	`
	now := time.Now()
	_, err := r.db.Exec(query, key, value, now, value, now)
	if err != nil {
		return fmt.Errorf("设置配置失败: %w", err)
	}
	return nil
}

// Get 获取配置项
func (r *SystemConfigRepository) Get(key string) (string, error) {
	query := `SELECT value FROM system_configs WHERE key = ?`
	var value string
	err := r.db.QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("获取配置失败: %w", err)
	}
	return value, nil
}

// GetWithDefault 获取配置项，如果不存在则返回默认值
func (r *SystemConfigRepository) GetWithDefault(key, defaultValue string) string {
	value, err := r.Get(key)
	if err != nil || value == "" {
		return defaultValue
	}
	return value
}

// GetAll 获取所有配置
func (r *SystemConfigRepository) GetAll() ([]*models.SystemConfig, error) {
	query := `SELECT key, value, updated_at FROM system_configs ORDER BY key`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询配置列表失败: %w", err)
	}
	defer rows.Close()

	var configs []*models.SystemConfig
	for rows.Next() {
		config := &models.SystemConfig{}
		if err := rows.Scan(&config.Key, &config.Value, &config.UpdatedAt); err != nil {
			return nil, fmt.Errorf("扫描配置数据失败: %w", err)
		}
		configs = append(configs, config)
	}
	return configs, nil
}

// Delete 删除配置项
func (r *SystemConfigRepository) Delete(key string) error {
	query := `DELETE FROM system_configs WHERE key = ?`
	_, err := r.db.Exec(query, key)
	if err != nil {
		return fmt.Errorf("删除配置失败: %w", err)
	}
	return nil
}

// GetByPrefix 根据前缀获取配置
func (r *SystemConfigRepository) GetByPrefix(prefix string) ([]*models.SystemConfig, error) {
	query := `SELECT key, value, updated_at FROM system_configs WHERE key LIKE ? ORDER BY key`
	rows, err := r.db.Query(query, prefix+"%")
	if err != nil {
		return nil, fmt.Errorf("查询配置列表失败: %w", err)
	}
	defer rows.Close()

	var configs []*models.SystemConfig
	for rows.Next() {
		config := &models.SystemConfig{}
		if err := rows.Scan(&config.Key, &config.Value, &config.UpdatedAt); err != nil {
			return nil, fmt.Errorf("扫描配置数据失败: %w", err)
		}
		configs = append(configs, config)
	}
	return configs, nil
}
