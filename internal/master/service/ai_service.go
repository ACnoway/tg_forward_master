package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AIService AI服务
type AIService struct {
	endpoint string
	apiKey   string
	model    string
	client   *http.Client
}

// NewAIService 创建AI服务
func NewAIService(endpoint, apiKey, model string) *AIService {
	return &AIService{
		endpoint: endpoint,
		apiKey:   apiKey,
		model:    model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AIRequest AI请求
type AIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
}

// AIMessage AI消息
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIResponse AI响应
type AIResponse struct {
	Choices []struct {
		Message AIMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// TestConnection 测试AI连接
func (s *AIService) TestConnection() (string, error) {
	// 构建测试请求
	request := AIRequest{
		Model: s.model,
		Messages: []AIMessage{
			{
				Role:    "user",
				Content: "Hello, this is a connection test. Please reply with 'OK'.",
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", s.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查错误
	if aiResp.Error != nil {
		return "", fmt.Errorf("AI错误: %s (%s)", aiResp.Error.Message, aiResp.Error.Type)
	}

	// 检查响应
	if len(aiResp.Choices) == 0 {
		return "", fmt.Errorf("无响应内容")
	}

	return aiResp.Choices[0].Message.Content, nil
}

// CheckSpam 检测垃圾消息
func (s *AIService) CheckSpam(message string) (bool, float64, error) {
	request := AIRequest{
		Model: s.model,
		Messages: []AIMessage{
			{
				Role:    "system",
				Content: "You are a spam detection assistant. Analyze the message and determine if it's spam. Reply with JSON format: {\"is_spam\": true/false, \"confidence\": 0.0-1.0, \"reason\": \"explanation\"}",
			},
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return false, 0, err
	}

	req, err := http.NewRequest("POST", s.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, err
	}

	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return false, 0, err
	}

	if aiResp.Error != nil {
		return false, 0, fmt.Errorf("%s", aiResp.Error.Message)
	}

	if len(aiResp.Choices) == 0 {
		return false, 0, fmt.Errorf("无响应内容")
	}

	// 解析AI返回的JSON
	var result struct {
		IsSpam     bool    `json:"is_spam"`
		Confidence float64 `json:"confidence"`
		Reason     string  `json:"reason"`
	}

	content := aiResp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// 如果无法解析JSON，尝试简单判断
		// 这里可以根据关键词判断
		return false, 0.5, nil
	}

	return result.IsSpam, result.Confidence, nil
}
