package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/gs-hms/authentication/dto"
	"github.com/gs-hms/authentication/middleware"
)

func generateToken(t *testing.T, jti string, exp time.Time) string {
	claims := dto.Claims{
		UserID: 1,
		Email:  "test@example.com",
		Name:   "Test User",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    "authentication-service",
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-30 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-30 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	assert.NoError(t, err)
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET_STRING", "secret")
	gin.SetMode(gin.TestMode)

	t.Run("valid token", func(t *testing.T) {
		client, mock := redismock.NewClientMock()
		mock.ExpectGet("jti-valid").RedisNil()

		r := gin.New()
		r.Use(middleware.AuthMiddleware(client))
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		token := generateToken(t, "jti-valid", time.Now().Add(15*time.Minute))
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("expired token", func(t *testing.T) {
		client, _ := redismock.NewClientMock()

		r := gin.New()
		r.Use(middleware.AuthMiddleware(client))
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		token := generateToken(t, "jti-expired", time.Now().Add(-15*time.Minute))
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "invalid token")
	})

	t.Run("token after logout (blacklisted)", func(t *testing.T) {
		client, mock := redismock.NewClientMock()
		mock.ExpectGet("jti-blacklisted").SetVal("blacklisted")

		r := gin.New()
		r.Use(middleware.AuthMiddleware(client))
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		token := generateToken(t, "jti-blacklisted", time.Now().Add(15*time.Minute))
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "token has been revoked")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("missing authorization header", func(t *testing.T) {
		client, _ := redismock.NewClientMock()

		r := gin.New()
		r.Use(middleware.AuthMiddleware(client))
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "authorization header is required")
	})

	t.Run("invalid authorization header format", func(t *testing.T) {
		client, _ := redismock.NewClientMock()

		r := gin.New()
		r.Use(middleware.AuthMiddleware(client))
		r.GET("/test", func(c *gin.Context) { c.Status(200) })

		token := generateToken(t, "jti-valid", time.Now().Add(15*time.Minute))
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", token) // Missing "Bearer " prefix
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "authorization header must be bearer token")
	})
}
