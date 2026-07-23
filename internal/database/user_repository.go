package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/acnoway/tg_forward_master/internal/models"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, is_admin)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, user.TelegramID, user.Username, user.FirstName, user.LastName, user.IsAdmin)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取用户ID失败: %w", err)
	}

	user.ID = id
	user.CreatedAt = time.Now()
	return nil
}

// GetByTelegramID 根据Telegram ID获取用户
func (r *UserRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, is_admin, created_at
		FROM users
		WHERE telegram_id = ?
	`
	user := &models.User{}
	err := r.db.QueryRow(query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.IsAdmin,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return user, nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, is_admin, created_at
		FROM users
		WHERE id = ?
	`
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.IsAdmin,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET username = ?, first_name = ?, last_name = ?, is_admin = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, user.Username, user.FirstName, user.LastName, user.IsAdmin, user.ID)
	if err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	return nil
}

// GetOrCreate 获取或创建用户
func (r *UserRepository) GetOrCreate(telegramID int64, username, firstName, lastName string) (*models.User, error) {
	user, err := r.GetByTelegramID(telegramID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		// 更新用户信息
		user.Username = username
		user.FirstName = firstName
		user.LastName = lastName
		if err := r.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// 创建新用户
	user = &models.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		IsAdmin:    false,
	}
	if err := r.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// SetAdmin 设置管理员
func (r *UserRepository) SetAdmin(telegramID int64, isAdmin bool) error {
	query := `UPDATE users SET is_admin = ? WHERE telegram_id = ?`
	_, err := r.db.Exec(query, isAdmin, telegramID)
	if err != nil {
		return fmt.Errorf("设置管理员失败: %w", err)
	}
	return nil
}

// GetAllAdmins 获取所有管理员
func (r *UserRepository) GetAllAdmins() ([]*models.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, is_admin, created_at
		FROM users
		WHERE is_admin = 1
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询管理员列表失败: %w", err)
	}
	defer rows.Close()

	var admins []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID,
			&user.TelegramID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.IsAdmin,
			&user.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描管理员数据失败: %w", err)
		}
		admins = append(admins, user)
	}
	return admins, nil
}

// Count 获取用户总数
func (r *UserRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询用户数量失败: %w", err)
	}
	return count, nil
}

// GetAll 获取所有用户（分页）
func (r *UserRepository) GetAll(limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, is_admin, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID,
			&user.TelegramID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.IsAdmin,
			&user.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描用户数据失败: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
