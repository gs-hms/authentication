// Package handler provides the HTTP handlers for the user service.
package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supermarios-hotel-management-system/authentication/dto"
	"github.com/supermarios-hotel-management-system/authentication/service"
)

// UserHandler defines the interface for the user handler.
type UserHandler interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
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
