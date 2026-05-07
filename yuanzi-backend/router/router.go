package router

import (
	"yuanzi-backend/config"
	"yuanzi-backend/handler"
	"yuanzi-backend/handler/admin"
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

	// Admin routes - authenticated with admin middleware
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.POST("/login", admin.AdminLogin)

	adminAuth := r.Group("/api/v1/admin")
	adminAuth.Use(middleware.JWT(), middleware.AdminAuth())

	// User management
	adminAuth.GET("/users", admin.GetUsers)
	adminAuth.GET("/users/:id", admin.GetUser)
	adminAuth.PUT("/users/:id/status", admin.UpdateUserStatus)
	adminAuth.DELETE("/users/:id", admin.DeleteUser)

	// Family management
	adminAuth.GET("/families", admin.GetFamilies)
	adminAuth.GET("/families/:id", admin.GetFamily)
	adminAuth.DELETE("/families/:id", admin.DeleteFamily)

	// Baby management
	adminAuth.GET("/babies", admin.GetBabies)
	adminAuth.GET("/babies/:id", admin.GetBaby)
	adminAuth.DELETE("/babies/:id", admin.DeleteBaby)

	// Photo management
	adminAuth.GET("/photos", admin.GetPhotos)
	adminAuth.GET("/photos/:id", admin.GetPhoto)
	adminAuth.DELETE("/photos/:id", admin.DeletePhoto)

	// Record management
	adminAuth.GET("/records", admin.GetRecords)
	adminAuth.GET("/records/:id", admin.GetRecord)
	adminAuth.DELETE("/records/:id", admin.DeleteRecord)

	// Statistics
	adminAuth.GET("/stats/overview", admin.GetStatsOverview)
	adminAuth.GET("/stats/daily", admin.GetDailyStats)

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
