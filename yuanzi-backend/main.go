package main

import (
	"fmt"
	"net/http"
	"time"
	"yuanzi-backend/config"
	_ "yuanzi-backend/docs"
	"yuanzi-backend/logger"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/gredis"
	"yuanzi-backend/router"
)

// @title Yuanzi API
// @version 1.0
// @description 小园子母婴记录应用后端 API (基于 PRD-V1.1)
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	config.Setup()
	logger.Setup()
	mysql.Setup()
	gredis.Setup()

	r := router.SetupRouter()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.GlobalConfig.Server.HttpPort),
		Handler:        r,
		ReadTimeout:    time.Duration(config.GlobalConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.GlobalConfig.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Info("Server starting",
		logger.String("addr", server.Addr),
		logger.String("mode", config.GlobalConfig.Server.RunMode),
	)

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Server start failed", logger.Err(err))
	}
}
