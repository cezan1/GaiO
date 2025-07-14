package main

import (
	"os"
	"strings"

	"github.com/cezan1/GaiO/internal/api/http"
	"github.com/cezan1/GaiO/internal/application"
	"github.com/cezan1/GaiO/internal/domain/service"
	"github.com/gin-gonic/gin"
)

func main() {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	router := gin.Default()
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		router.SetTrustedProxies(strings.Split(trustedProxies, ","))
	}
	//初始化应用服务
	aiService := service.NewAIService()
	aiAppService := application.NewAIAppService(aiService)
	//注册路由
	aiHandler := http.NewAIHandler(aiAppService)
	aiHandler.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
