package handlers

import (
	"sync"
	"time"
)

// SessionState 会话状态
type SessionState struct {
	Command    string                 // 当前命令
	Step       int                    // 当前步骤
	Data       map[string]interface{} // 临时数据
	CreatedAt  time.Time              // 创建时间
	ExpireAt   time.Time              // 过期时间
}

// SessionManager 会话管理器
type SessionManager struct {
	sessions map[int64]*SessionState
	mu       sync.RWMutex
}

// NewSessionManager 创建会话管理器
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[int64]*SessionState),
	}
	// 启动清理goroutine
	go sm.cleanupExpired()
	return sm
}

// StartSession 开始会话
func (sm *SessionManager) StartSession(userID int64, command string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	sm.sessions[userID] = &SessionState{
		Command:   command,
		Step:      0,
		Data:      make(map[string]interface{}),
		CreatedAt: now,
		ExpireAt:  now.Add(5 * time.Minute), // 5分钟过期
	}
}

// GetSession 获取会话
func (sm *SessionManager) GetSession(userID int64) (*SessionState, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	state, exists := sm.sessions[userID]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(state.ExpireAt) {
		return nil, false
	}

	return state, true
}

// UpdateSession 更新会话
func (sm *SessionManager) UpdateSession(userID int64, step int, data map[string]interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if state, exists := sm.sessions[userID]; exists {
		state.Step = step
		for k, v := range data {
			state.Data[k] = v
		}
		state.ExpireAt = time.Now().Add(5 * time.Minute) // 重置过期时间
	}
}

// EndSession 结束会话
func (sm *SessionManager) EndSession(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, userID)
}

// cleanupExpired 清理过期会话
func (sm *SessionManager) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()
		for userID, state := range sm.sessions {
			if now.After(state.ExpireAt) {
				delete(sm.sessions, userID)
			}
		}
		sm.mu.Unlock()
	}
}
