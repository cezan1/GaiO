package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 统一请求结构
// 示例json
//
//	{
//		"params": {
//			 "param1": "111"
//			 "param2": "222"
//		},
//		"trace": {
//			"request_id": "123456"
//		}
//	}
type Request struct {
	Params map[string]interface{} `json:"params"` // 业务数据
	Trace  Trace                  `json:"trace"`  // 响应唯一标识
}

// 统一响应结构
// 示例json
//
//	{
//		"code": 0,
//		"message": "success",
//		"data": {
//			"answer": "..."
//		},
//		"trace": {
//			"request_id": "123456"
//		}
//	}
type Response struct {
	Code    int         `json:"code"`    // 响应码，0 表示成功
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 业务数据
	Trace   Trace       `json:"trace"`   // 响应唯一标识
}

type Trace struct {
	RequestID string `json:"request_id"` // 请求唯一标识
}

// BindRequest 绑定统一请求结构
func BindRequest(c *gin.Context, req *Request) (*Request, error) {
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	//校验参数
	if req.Params == nil {
		return nil, errors.New("params is required")
	}

	//校验trace
	if req.Trace == (Trace{}) {
		req.Trace = Trace{
			RequestID: GenerateRequestID(),
		}
	}
	if req.Trace.RequestID == "" {
		//生成request_id
		req.Trace.RequestID = GenerateRequestID()
	}

	return req, nil
}

// SendResponse 发送统一响应结构
func SendResponse(c *gin.Context, code int, message string, data interface{}, requestID string) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	resp.Trace.RequestID = requestID
	c.JSON(code, resp)
}

// GenerateRequestID 生成请求唯一标识
func GenerateRequestID() string {
	return uuid.New().String()
}
