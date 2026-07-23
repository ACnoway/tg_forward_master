package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/acnoway/tg_forward_master/internal/database"
	"github.com/acnoway/tg_forward_master/internal/models"
)

// AdminService 管理员服务
type AdminService struct {
	userRepo         *database.UserRepository
	planRepo         *database.PlanRepository
	redeemCodeRepo   *database.RedeemCodeRepository
	subscriptionRepo *database.SubscriptionRepository
	workerBotRepo    *database.WorkerBotRepository
	configRepo       *database.SystemConfigRepository
}

// NewAdminService 创建管理员服务
func NewAdminService(
	userRepo *database.UserRepository,
	planRepo *database.PlanRepository,
	redeemCodeRepo *database.RedeemCodeRepository,
	subscriptionRepo *database.SubscriptionRepository,
	workerBotRepo *database.WorkerBotRepository,
	configRepo *database.SystemConfigRepository,
) *AdminService {
	return &AdminService{
		userRepo:         userRepo,
		planRepo:         planRepo,
		redeemCodeRepo:   redeemCodeRepo,
		subscriptionRepo: subscriptionRepo,
		workerBotRepo:    workerBotRepo,
		configRepo:       configRepo,
	}
}

// CreatePlan 创建套餐
func (s *AdminService) CreatePlan(name string, price float64, durationDays, maxBots int) (*models.Plan, error) {
	if name == "" {
		return nil, fmt.Errorf("套餐名称不能为空")
	}
	if price < 0 {
		return nil, fmt.Errorf("价格不能为负数")
	}
	if durationDays <= 0 {
		return nil, fmt.Errorf("有效期必须大于0天")
	}
	if maxBots <= 0 {
		return nil, fmt.Errorf("Bot数量必须大于0")
	}

	plan := &models.Plan{
		Name:         name,
		Price:        price,
		DurationDays: durationDays,
		MaxBots:      maxBots,
		IsActive:     true,
	}

	if err := s.planRepo.Create(plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// UpdatePlan 更新套餐
func (s *AdminService) UpdatePlan(planID int64, name string, price float64, durationDays, maxBots int, isActive bool) error {
	plan, err := s.planRepo.GetByID(planID)
	if err != nil {
		return err
	}
	if plan == nil {
		return fmt.Errorf("套餐不存在")
	}

	plan.Name = name
	plan.Price = price
	plan.DurationDays = durationDays
	plan.MaxBots = maxBots
	plan.IsActive = isActive

	return s.planRepo.Update(plan)
}

// DeletePlan 删除套餐
func (s *AdminService) DeletePlan(planID int64) error {
	return s.planRepo.Delete(planID)
}

// GetAllPlans 获取所有套餐
func (s *AdminService) GetAllPlans() ([]*models.Plan, error) {
	return s.planRepo.GetAll()
}

// GenerateRedeemCodes 批量生成兑换码
func (s *AdminService) GenerateRedeemCodes(planID int64, count int) ([]*models.RedeemCode, error) {
	if count <= 0 || count > 100 {
		return nil, fmt.Errorf("生成数量必须在1-100之间")
	}

	// 检查套餐是否存在
	plan, err := s.planRepo.GetByID(planID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, fmt.Errorf("套餐不存在")
	}

	codes := make([]*models.RedeemCode, 0, count)
	for i := 0; i < count; i++ {
		code := &models.RedeemCode{
			Code:   s.generateRandomCode(),
			PlanID: planID,
			Status: "unused",
		}

		if err := s.redeemCodeRepo.Create(code); err != nil {
			return nil, fmt.Errorf("生成第%d个兑换码失败: %w", i+1, err)
		}

		codes = append(codes, code)
	}

	return codes, nil
}

// generateRandomCode 生成随机兑换码
func (s *AdminService) generateRandomCode() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetRedeemCodes 获取兑换码列表
func (s *AdminService) GetRedeemCodes(limit, offset int) ([]*models.RedeemCode, error) {
	return s.redeemCodeRepo.GetAll(limit, offset)
}

// GetUnusedRedeemCodes 获取未使用的兑换码
func (s *AdminService) GetUnusedRedeemCodes() ([]*models.RedeemCode, error) {
	return s.redeemCodeRepo.GetUnused()
}

// DeleteRedeemCode 删除兑换码
func (s *AdminService) DeleteRedeemCode(codeID int64) error {
	return s.redeemCodeRepo.Delete(codeID)
}

// SetPaymentConfig 设置支付配置
func (s *AdminService) SetPaymentConfig(apiURL, merchantID, merchantKey string) error {
	if apiURL == "" {
		return fmt.Errorf("API地址不能为空")
	}
	if merchantID == "" {
		return fmt.Errorf("商户ID不能为空")
	}
	if merchantKey == "" {
		return fmt.Errorf("商户密钥不能为空")
	}

	if err := s.configRepo.Set("epay_api_url", apiURL); err != nil {
		return err
	}
	if err := s.configRepo.Set("epay_merchant_id", merchantID); err != nil {
		return err
	}
	if err := s.configRepo.Set("epay_merchant_key", merchantKey); err != nil {
		return err
	}

	return nil
}

// GetPaymentConfig 获取支付配置
func (s *AdminService) GetPaymentConfig() (apiURL, merchantID, merchantKey string, err error) {
	apiURL, err = s.configRepo.Get("epay_api_url")
	if err != nil {
		return "", "", "", err
	}
	merchantID, err = s.configRepo.Get("epay_merchant_id")
	if err != nil {
		return "", "", "", err
	}
	merchantKey, err = s.configRepo.Get("epay_merchant_key")
	if err != nil {
		return "", "", "", err
	}
	return apiURL, merchantID, merchantKey, nil
}

// SetAIConfig 设置AI配置
func (s *AdminService) SetAIConfig(endpoint, apiKey, model string) error {
	if endpoint == "" {
		return fmt.Errorf("API地址不能为空")
	}
	if apiKey == "" {
		return fmt.Errorf("API密钥不能为空")
	}
	if model == "" {
		return fmt.Errorf("模型名称不能为空")
	}

	if err := s.configRepo.Set("default_ai_endpoint", endpoint); err != nil {
		return err
	}
	if err := s.configRepo.Set("default_ai_key", apiKey); err != nil {
		return err
	}
	if err := s.configRepo.Set("default_ai_model", model); err != nil {
		return err
	}

	return nil
}

// GetAIConfig 获取AI配置
func (s *AdminService) GetAIConfig() (endpoint, apiKey, model string, err error) {
	endpoint, err = s.configRepo.Get("default_ai_endpoint")
	if err != nil {
		return "", "", "", err
	}
	apiKey, err = s.configRepo.Get("default_ai_key")
	if err != nil {
		return "", "", "", err
	}
	model, err = s.configRepo.Get("default_ai_model")
	if err != nil {
		return "", "", "", err
	}
	return endpoint, apiKey, model, nil
}

// GetSystemStats 获取系统统计
func (s *AdminService) GetSystemStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 用户统计
	userCount, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}
	stats["user_count"] = userCount

	// 管理员数量
	admins, err := s.userRepo.GetAllAdmins()
	if err != nil {
		return nil, err
	}
	stats["admin_count"] = len(admins)

	// 订阅统计
	subCount, err := s.subscriptionRepo.Count()
	if err != nil {
		return nil, err
	}
	stats["subscription_count"] = subCount

	activeSubCount, err := s.subscriptionRepo.CountActive()
	if err != nil {
		return nil, err
	}
	stats["active_subscription_count"] = activeSubCount

	// 套餐统计
	planCount, err := s.planRepo.Count()
	if err != nil {
		return nil, err
	}
	stats["plan_count"] = planCount

	// Bot统计
	botCount, err := s.workerBotRepo.Count()
	if err != nil {
		return nil, err
	}
	stats["bot_count"] = botCount

	runningBotCount, err := s.workerBotRepo.CountRunning()
	if err != nil {
		return nil, err
	}
	stats["running_bot_count"] = runningBotCount

	// 兑换码统计
	codeCount, err := s.redeemCodeRepo.Count()
	if err != nil {
		return nil, err
	}
	stats["redeem_code_count"] = codeCount

	unusedCodeCount, err := s.redeemCodeRepo.CountUnused()
	if err != nil {
		return nil, err
	}
	stats["unused_code_count"] = unusedCodeCount

	return stats, nil
}

// SetAdmin 设置管理员
func (s *AdminService) SetAdmin(telegramID int64, isAdmin bool) error {
	return s.userRepo.SetAdmin(telegramID, isAdmin)
}

// GetAllUsers 获取所有用户（分页）
func (s *AdminService) GetAllUsers() (int, error) {
	return s.userRepo.Count()
}

// ExpireOldSubscriptions 清理过期订阅
func (s *AdminService) ExpireOldSubscriptions() (int, error) {
	return s.subscriptionRepo.ExpireOldSubscriptions()
}

// FormatDuration 格式化时长
func FormatDuration(days int) string {
	if days >= 365 {
		years := days / 365
		return strconv.Itoa(years) + "年"
	}
	if days >= 30 {
		months := days / 30
		return strconv.Itoa(months) + "个月"
	}
	return strconv.Itoa(days) + "天"
}

// FormatPrice 格式化价格
func FormatPrice(price float64) string {
	if price == 0 {
		return "免费"
	}
	return fmt.Sprintf("¥%.2f", price)
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
