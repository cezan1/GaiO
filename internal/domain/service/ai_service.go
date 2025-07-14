package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cezan1/GaiO/internal/config"
	"github.com/cezan1/GaiO/internal/domain/model"
	"github.com/cezan1/GaiO/pkg/helper"
)

// AIService 定义领域服务接口
type AIService interface {
	GetAIAnswer(ctx context.Context, req []*model.AIRequest, modelName string) (model.AIResponse, error)
}

// AIServiceImpl 实现 AIService 接口
type AIServiceImpl struct{}

// NewAIService 创建 AIService 实例
func NewAIService() AIService {
	return &AIServiceImpl{}
}

var QuestionList []*AIRequest

type AIRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GetAIAnswer 调用第三方 AI 服务获取答案
func (s *AIServiceImpl) GetAIAnswer(ctx context.Context, req []*model.AIRequest, modelName string) (model.AIResponse, error) {
	thirdPartyReq := model.ThirdPartyAIRequest{
		Model: modelName,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{},
	}
	for _, q := range req {
		thirdPartyReq.Messages = append(thirdPartyReq.Messages, struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			Role:    q.Role,
			Content: q.Content,
		})
	}

	jsonData, err := json.Marshal(thirdPartyReq)
	if err != nil {
		return model.AIResponse{}, fmt.Errorf("marshal request failed: %w", err)
	}
	//请求ai参数写入日志
	helper.WriteLogToFile("GetAIAnswer.request : " + string(jsonData))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", config.ThirdPartyAIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return model.AIResponse{}, fmt.Errorf("create http request failed: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+config.DASHSCOPE_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return model.AIResponse{}, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	var thirdPartyResp model.ThirdPartyAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&thirdPartyResp); err != nil {
		return model.AIResponse{}, fmt.Errorf("decode response failed: %w", err)
	}

	var aiResponse model.AIResponse
	if len(thirdPartyResp.Choices) > 0 {
		aiResponse.Answer = thirdPartyResp.Choices[0].Message.Content
	}
	return aiResponse, nil
}
