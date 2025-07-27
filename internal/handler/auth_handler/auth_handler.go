package auth_handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ipxsandbox/internal/entity"
	"github.com/ipxsandbox/internal/usecase/auth_usercase"
	customValidator "github.com/ipxsandbox/internal/validator"
)

type AuthHandler struct {
	authUsecase auth_usercase.AuthUsecaseInterface
}

func NewAuthHandler(auc auth_usercase.AuthUsecaseInterface) *AuthHandler {
	return &AuthHandler{authUsecase: auc}
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	customValidator.RegisterCustomValidators(validate)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var userData entity.User
	err := c.ShouldBind(&userData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	if err := validate.Struct(userData); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": customValidator.TranslateValidationError(err)})
		return
	}

	resp, err := h.authUsecase.Register(userData)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.Login(userData.Email, userData.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.SetCookie("access_token", accessToken, 60*15, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 60*60*24*7, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "login success"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	newAccessToken, err := h.authUsecase.RefreshAccessToken(refreshToken)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	c.SetCookie("access_token", newAccessToken, 60*15, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}
