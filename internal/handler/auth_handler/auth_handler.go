package auth_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ipxsandbox/internal/entity"
	"github.com/ipxsandbox/internal/usecase/auth_usercase"
)

type AuthHandler struct {
	authUsecase auth_usercase.AuthUsecaseInterface
}

func NewAuthHandler(auc auth_usercase.AuthUsecaseInterface) *AuthHandler {
	return &AuthHandler{authUsecase: auc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var userData entity.User
	err := c.ShouldBindBodyWithJSON(&userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	resp, err := h.authUsecase.Register(userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var userData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.Login(userData.Email, userData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}