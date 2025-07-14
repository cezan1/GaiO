package model

// AIRequest 定义调用第三方 AI 服务的请求结构体
// 传参方式：json
type AIRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIResponse 定义调用第三方 AI 服务的响应结构体
type AIResponse struct {
	Answer string `json:"answer"`
}

// 第三方 AI 请求结构体
type ThirdPartyAIRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// 第三方 AI 响应结构体
type ThirdPartyAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
