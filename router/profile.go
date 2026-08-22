package router

import (
	"github.com/gin-gonic/gin"
	"github.com/supermarios-hotel-management-system/authentication/handler"
	"github.com/supermarios-hotel-management-system/authentication/middleware"
	"github.com/supermarios-hotel-management-system/authentication/repository"
	"github.com/supermarios-hotel-management-system/authentication/service"
)

// RegisterProfileRoutes registers the profile routes.
func RegisterProfileRoutes(r *gin.RouterGroup, userRepo repository.UserRepository) {
	profileService := service.NewProfileService(userRepo)
	profileHandler := handler.NewProfileHandler(profileService)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/", profileHandler.GetProfile)
		protected.PUT("/", profileHandler.UpdateProfile)
		protected.PUT("/password", profileHandler.ChangePassword)
	}
}
