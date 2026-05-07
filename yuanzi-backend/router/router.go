package router

import (
	"yuanzi-backend/config"
	"yuanzi-backend/handler"
	"yuanzi-backend/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	gin.SetMode(config.GlobalConfig.Server.RunMode)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", handler.HealthCheck)
	r.GET("/ping", handler.Ping)

	registerAPIRoutes(r, "/api")
	registerAPIRoutes(r, "/api/v1")

	return r
}

func registerAPIRoutes(r *gin.Engine, base string) {
	api := r.Group(base)

	api.POST("/photo/callback", handler.PhotoUploadCallback)

	auth := api.Group("/auth")
	auth.POST("/send-code", handler.SendVerificationCode)
	auth.POST("/sms", handler.SendVerificationCode)
	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.RefreshToken)
	auth.POST("/logout", handler.Logout)

	authorized := api.Group("")
	authorized.Use(middleware.JWT())
	authorized.Use(middleware.RequireDB())

	user := authorized.Group("/user")
	user.GET("/profile", handler.GetUserProfile)
	user.PUT("/profile", handler.UpdateUserProfile)

	family := authorized.Group("/family")
	family.POST("", handler.CreateFamily)
	family.GET("/:id", handler.GetFamily)
	family.POST("/:id/invite", handler.InviteFamilyMember)
	family.GET("/:id/members", handler.GetFamilyMembers)
	family.DELETE("/:id/members/:userId", handler.RemoveFamilyMember)

	baby := authorized.Group("/baby")
	baby.POST("", handler.CreateBaby)
	baby.GET("", handler.ListBabies)
	baby.GET("/:id", handler.GetBaby)
	baby.PUT("/:id", handler.UpdateBaby)
	baby.DELETE("/:id", handler.DeleteBaby)

	photo := authorized.Group("/photo")
	photo.POST("/upload-url", handler.GetPhotoUploadURL)
	photo.GET("", handler.ListPhotos)
	photo.DELETE("/:id", handler.DeletePhoto)

	record := authorized.Group("/record")
	record.POST("", handler.CreateRecord)
	record.GET("", handler.ListRecords)
	record.GET("/:id", handler.GetRecord)
	record.PUT("/:id", handler.UpdateRecord)
	record.DELETE("/:id", handler.DeleteRecord)

	sync := authorized.Group("/sync")
	sync.GET("/stream", handler.SSEStream)

	device := authorized.Group("/device")
	device.POST("/register", handler.RegisterDevice)

	stats := authorized.Group("/stats")
	stats.GET("/daily", handler.GetDailyStats)
	stats.GET("/weekly", handler.GetWeeklyStats)

	ai := authorized.Group("/ai")
	ai.POST("/chat", handler.AIChat)
	ai.POST("/speech/recognize", handler.SpeechRecognize)
	ai.GET("/quota", handler.GetAIQuota)
}
