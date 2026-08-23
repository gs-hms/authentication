// Authentication service for Hotel Management System.
package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gs-hms/authentication/database"
	"github.com/gs-hms/authentication/kafka"
	"github.com/gs-hms/authentication/repository"
	"github.com/gs-hms/authentication/router"
	"github.com/gs-hms/event-handler/events"
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

	redisConn, err := database.ConnectRedis(ctx)
	if err != nil {
		log.Printf("Failed to connect to redis: %v", err)
		redisConn = &database.Redis{Client: nil}
	}

	kafkaHost := os.Getenv("KAFKA_BROKERS")
	if strings.TrimSpace(kafkaHost) == "" {
		log.Fatal("KAFKA_BROKERS environment variable is not set")
	}

	kafkaProducer := kafka.NewProducer(
		strings.Split(kafkaHost, ","),
		events.KafkaTopics,
	)
	defer kafkaProducer.Close()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	authSessionRepo := repository.NewAuthenticationSessionRepository(db)
	userRepo := repository.NewUserRepository(db)

	v1 := r.Group("/v1/authentication")
	userRouter := v1.Group("/user")
	router.RegisterUserRoutes(userRouter, userRepo, authSessionRepo, redisConn.Client, kafkaProducer)
	profileRouter := v1.Group("/profile")
	router.RegisterProfileRoutes(profileRouter, userRepo, redisConn.Client)

	err = r.Run()
	if err != nil {
		log.Fatalf("start server: %v", err)
	}
}
