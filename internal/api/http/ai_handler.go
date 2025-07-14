package http

import (
	"fmt"
	"net/http"

	"github.com/cezan1/GaiO/internal/application"
	"github.com/cezan1/GaiO/internal/config"
	"github.com/cezan1/GaiO/internal/domain/model"
	helper "github.com/cezan1/GaiO/pkg/helper"
	"github.com/gin-gonic/gin"
)

// AIHandler 定义 AI 相关接口的处理结构体
type AIHandler struct {
	aiAppService *application.AIAppService
}

// NewAIHandler 创建 AIHandler 实例
func NewAIHandler(aiAppService *application.AIAppService) *AIHandler {
	return &AIHandler{
		aiAppService: aiAppService,
	}
}

// RegisterRoutes 注册 AI 相关路由
func (h *AIHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api")
	{
		group.POST("/ai", h.handleAIRequest)
	}
}

// handleAIRequest 处理 AI 请求
func (h *AIHandler) handleAIRequest(c *gin.Context) {
	// 绑定请求参数
	req := &helper.Request{}
	req, err := helper.BindRequest(c, req)
	if err != nil {
		helper.SendResponse(c, http.StatusBadRequest, err.Error(), nil, "")
		return
	}
	//校验参数
	question, ok := req.Params["question"].(string)
	if !ok || question == "" {
		helper.SendResponse(c, http.StatusBadRequest, "question is required", nil, req.Trace.RequestID)
		return
	}
	//获取文本模型响应
	aiReq := model.AIRequest{
		Role:    "user",
		Content: question,
	}

	//写入日志
	msg := fmt.Sprintf("RequestID: %s, question: %s", req.Trace.RequestID, question)
	helper.WriteLog(msg)

	//历史问题
	questionList, err := helper.GetDataByRequestIDDesc(req.Trace.RequestID)
	if err != nil {
		helper.SendResponse(c, http.StatusInternalServerError, err.Error(), nil, req.Trace.RequestID)
		return
	}
	var aiReqList []*model.AIRequest
	if questionList != nil {
		l := len(questionList)
		for i := l - 1; i >= 0; i-- {
			//与当前问题重复去重
			if questionList[i] == question {
				continue
			}
			aiReqList = append(aiReqList, &model.AIRequest{
				Role:    "user",
				Content: questionList[i],
			})
		}
	}
	//最后问题追加
	aiReqList = append(aiReqList, &aiReq)
	//写入关联问题数据
	helper.WriteRequestData(req.Trace.RequestID, question)

	resp, err := h.aiAppService.GetAIAnswer(c, aiReqList, config.BiggiBModelName)
	if err != nil {
		helper.SendResponse(c, http.StatusInternalServerError, err.Error(), nil, req.Trace.RequestID)
		return
	}

	// 统一返回参数
	helper.SendResponse(c, http.StatusOK, "success", resp.Answer, req.Trace.RequestID)
}
