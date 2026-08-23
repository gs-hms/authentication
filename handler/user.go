// Package handler provides the HTTP handlers for the user service.
package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gs-hms/authentication/dto"
	"github.com/gs-hms/authentication/service"
)

// UserHandler defines the interface for the user handler.
type UserHandler interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new instance of the UserHandler.
func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.userService.Signup(c, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserWithEmailExists) || errors.Is(err, service.ErrUserWithPhoneExists) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrInvalidEmail) || errors.Is(err, service.ErrInvalidPhone) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Signup successful",
	})
}

func (h *userHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userService.Login(c, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrInvalidCredentials) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrInactiveUser) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *userHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jtiVal, exists := c.Get("jti")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "jti not found in context"})
		return
	}
	jti, ok := jtiVal.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid jti format"})
		return
	}

	expVal, exists := c.Get("exp")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "exp not found in context"})
		return
	}
	exp, ok := expVal.(time.Time)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid exp format"})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
		return
	}
	userID, ok := userIDVal.(uint64)
	if !ok {
		// Float64 check since JWT might parse numbers as float64 depending on how it's handled,
		// though claims.UserID is uint64. But let's safely cast it or try float64.
		userIDFloat, okFloat := userIDVal.(float64)
		if okFloat {
			userID = uint64(userIDFloat)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id format"})
			return
		}
	}

	if err := h.userService.Logout(c, userID, jti, exp, &req); err != nil {
		if err.Error() == "unauthorized session revocation" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
