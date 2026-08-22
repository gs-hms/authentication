// Authentication service for Hotel Management System.
package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/supermarios-hotel-management-system/authentication/database"
	"github.com/supermarios-hotel-management-system/authentication/repository"
	"github.com/supermarios-hotel-management-system/authentication/router"
)

func main() {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(databaseURL) == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := database.NewPostgres(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer db.Pool.Close()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	authSessionRepo := repository.NewAuthenticationSessionRepository(db)
	userRepo := repository.NewUserRepository(db)

	v1 := r.Group("/v1")
	userRouter := v1.Group("/user")
	router.RegisterUserRoutes(userRouter, userRepo, authSessionRepo)

	err = r.Run()
	if err != nil {
		log.Fatalf("start server: %v", err)
	}
}
