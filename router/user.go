// Package router provides the router for the user service.
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/supermarios-hotel-management-system/authentication/handler"
	"github.com/supermarios-hotel-management-system/authentication/middleware"
	"github.com/supermarios-hotel-management-system/authentication/repository"
	"github.com/supermarios-hotel-management-system/authentication/service"
)

// RegisterUserRoutes registers the user routes.
func RegisterUserRoutes(r *gin.RouterGroup, userRepo repository.UserRepository, authSessionRepo repository.AuthenticationSessionRepository, redisClient *redis.Client) {
	service := service.NewUserService(userRepo, authSessionRepo, redisClient)
	handler := handler.NewUserHandler(service)
	r.POST("/signup", handler.Signup)
	r.POST("/login", handler.Login)
	r.POST("/logout", middleware.AuthMiddleware(redisClient), handler.Logout)
}
