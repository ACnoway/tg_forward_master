package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB 数据库连接封装
type DB struct {
	*sql.DB
}

// New 创建新的数据库连接
func New(dbPath string) (*DB, error) {
	// 确保数据库目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %w", err)
	}

	// 打开数据库
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return &DB{db}, nil
}

// RunMigrations 运行数据库迁移
func (db *DB) RunMigrations(migrationPath string) error {
	// 读取迁移文件
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("读取迁移文件失败: %w", err)
	}

	// 执行迁移
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("执行迁移失败: %w", err)
	}

	return nil
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	return db.DB.Close()
}
