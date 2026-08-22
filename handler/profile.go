package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supermarios-hotel-management-system/authentication/dto"
	"github.com/supermarios-hotel-management-system/authentication/service"
)

// ProfileHandler defines the interface for the profile handler.
type ProfileHandler interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	ChangePassword(c *gin.Context)
}

type profileHandler struct {
	profileService service.ProfileService
}

// NewProfileHandler creates a new instance of ProfileHandler.
func NewProfileHandler(profileService service.ProfileService) ProfileHandler {
	return &profileHandler{
		profileService: profileService,
	}
}

func (h *profileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID.(uint64))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *profileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.profileService.UpdateProfile(c.Request.Context(), userID.(uint64), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPhone) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}

func (h *profileHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.profileService.ChangePassword(c.Request.Context(), userID.(uint64), &req)
	if err != nil {
		if errors.Is(err, service.ErrIncorrectPassword) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}
