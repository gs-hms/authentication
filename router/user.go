// Package router provides the router for the user service.
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gs-hms/authentication/handler"
	"github.com/gs-hms/authentication/kafka"
	"github.com/gs-hms/authentication/middleware"
	"github.com/gs-hms/authentication/repository"
	"github.com/gs-hms/authentication/service"
	"github.com/redis/go-redis/v9"
)

// RegisterUserRoutes registers the user routes.
func RegisterUserRoutes(r *gin.RouterGroup, userRepo repository.UserRepository, authSessionRepo repository.AuthenticationSessionRepository, redisClient *redis.Client, kafkaProducer *kafka.Producer) {
	service := service.NewUserService(userRepo, authSessionRepo, redisClient, kafkaProducer)
	handler := handler.NewUserHandler(service)
	r.POST("/signup", handler.Signup)
	r.POST("/login", handler.Login)
	r.POST("/logout", middleware.AuthMiddleware(redisClient), handler.Logout)
}
