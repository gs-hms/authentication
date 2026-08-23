package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gs-hms/authentication/handler"
	"github.com/gs-hms/authentication/middleware"
	"github.com/gs-hms/authentication/repository"
	"github.com/gs-hms/authentication/service"
	"github.com/redis/go-redis/v9"
)

// RegisterProfileRoutes registers the profile routes.
func RegisterProfileRoutes(r *gin.RouterGroup, userRepo repository.UserRepository, redisClient *redis.Client) {
	profileService := service.NewProfileService(userRepo)
	profileHandler := handler.NewProfileHandler(profileService)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient))
	{
		protected.GET("/", profileHandler.GetProfile)
		protected.PUT("/", profileHandler.UpdateProfile)
		protected.PUT("/password", profileHandler.ChangePassword)
	}
}
