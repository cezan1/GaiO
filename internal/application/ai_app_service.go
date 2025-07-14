package application

import (
	"context"
	"fmt"
	"log"

	"github.com/cezan1/GaiO/internal/domain/model"
	"github.com/cezan1/GaiO/internal/domain/service"
)

// AIAppService 定义应用层服务接口
type AIAppService struct {
	aiService service.AIService
}

// NewAIAppService 创建 AIAppService 实例
func NewAIAppService(aiService service.AIService) *AIAppService {
	return &AIAppService{
		aiService: aiService,
	}
}

// GetAIAnswer 调用领域服务获取 AI 答案
func (s *AIAppService) GetAIAnswer(ctx context.Context, req []*model.AIRequest, modelName string) (model.AIResponse, error) {
	ans, err := s.aiService.GetAIAnswer(ctx, req, modelName)
	if err != nil || ans.Answer == "" {
		log.Printf("获取AI答案失败: %v", err)
		return model.AIResponse{}, fmt.Errorf("获取AI答案失败: %w", err)
	}
	return ans, nil
}
