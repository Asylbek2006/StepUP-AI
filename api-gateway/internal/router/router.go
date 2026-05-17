package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/stepup-ai/api-gateway/internal/handler"
	"github.com/stepup-ai/api-gateway/internal/middleware"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	universityHandler *handler.UniversityHandler,
	aiHandler *handler.AIHandler,
	jwtSecretKey string,
) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authMiddleware := middleware.JWTAuthMiddleware(jwtSecretKey)

	auth := router.Group("/auth")
	{
		auth.POST("/register", userHandler.RegisterUser)
		auth.POST("/login", userHandler.LoginUser)
		auth.POST("/refresh", userHandler.RefreshToken)
		auth.POST("/forgot-password", userHandler.SendPasswordResetEmail)
		auth.POST("/reset-password", userHandler.ResetPassword)
		auth.POST("/logout", authMiddleware, userHandler.LogoutUser)
	}

	profile := router.Group("/profile", authMiddleware)
	{
		profile.GET("", userHandler.GetUserProfile)
		profile.PUT("", userHandler.UpdateUserProfile)
	}

	universities := router.Group("/universities", authMiddleware)
	{
		universities.GET("", universityHandler.SearchUniversities)
		universities.GET("/:id", universityHandler.GetUniversityDetails)
		universities.POST("/save", universityHandler.SaveUniversity)
		universities.GET("/saved", universityHandler.GetSavedUniversities)
	}

	grants := router.Group("/grants", authMiddleware)
	{
		grants.GET("", universityHandler.SearchGrants)
		grants.POST("/save", universityHandler.SaveGrant)
		grants.GET("/saved", universityHandler.GetSavedGrants)
	}

	ai := router.Group("/ai", authMiddleware)
	{
		ai.POST("/analyze", aiHandler.AnalyzeAdmissionChances)
		ai.POST("/roadmap", aiHandler.GenerateRoadmap)
		ai.POST("/essay", aiHandler.ReviewEssay)
		ai.POST("/grants/match", aiHandler.MatchGrants)
		ai.GET("/history", aiHandler.GetAnalysisHistory)
	}

	return router
}
