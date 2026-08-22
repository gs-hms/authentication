package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/supermarios-hotel-management-system/authentication/dto"
)

// AuthMiddleware is a gin middleware that authenticates requests using JWT.
func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be bearer token"})
			return
		}

		tokenString := parts[1]
		secretString := os.Getenv("JWT_SECRET_STRING")
		if secretString == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "jwt secret is not configured"})
			return
		}

		claims := &dto.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretString), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Check if token is blacklisted in Redis
		if redisClient != nil {
			err := redisClient.Get(context.Background(), claims.ID).Err()
			if err != redis.Nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.Name)
		c.Set("email", claims.Email)
		c.Set("jti", claims.ID)
		if claims.ExpiresAt != nil {
			c.Set("exp", claims.ExpiresAt.Time)
		}
		c.Next()
	}
}
