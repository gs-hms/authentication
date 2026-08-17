package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/supermarios-hotel-management-system/authentication/database"
)

func main() {
	ctx := context.Background()
	databseUrl := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(databseUrl) == "" {
		panic("DATABASE_URL environment variable is not set")
	}

	db, err := database.NewPostgres(ctx, databseUrl)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	defer db.Pool.Close()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	r.Run()
}
