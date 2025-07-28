package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"
	rdb "github.com/redis/go-redis/v9"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ipxsandbox/internal/entity"
	"github.com/ipxsandbox/internal/pkg/redis"
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

const (
	maxLoginAttempts = 5
	initialBlockTime = 5 * time.Minute
)

func (h *AuthHandler) Register(c *gin.Context) {
	var userData entity.User
	err := c.ShouldBindJSON(&userData)
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

func (h *AuthHandler) isBlocked(c *gin.Context, email string) bool {
	blockKey := fmt.Sprintf("login_blocked:%s", email)
	blockTTL, err := redis.Rdb.TTL(redis.Ctx, blockKey).Result()
	if err != nil {
		log.Println("Redis TTL error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return true
	}
	if blockTTL > 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": fmt.Sprintf("Too many failed attempts. Try again in %v", blockTTL.Round(time.Second)),
		})
		return true
	}
	return false
}

func (h *AuthHandler) handleLoginSuccess(c *gin.Context, email, accessToken, refreshToken string) {
	attemptKey := fmt.Sprintf("login_attempt:%s", email)
	blockKey := fmt.Sprintf("login_blocked:%s", email)

	if err := redis.Rdb.Del(redis.Ctx, attemptKey, blockKey).Err(); err != nil {
		log.Println("Failed to delete Redis keys after login:", err)
	}

	c.SetCookie("access_token", accessToken, 60*15, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 60*60*24*7, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "login success"})
}

func (h *AuthHandler) handleLoginFailure(c *gin.Context, email string) {
	attemptKey := fmt.Sprintf("login_attempt:%s", email)
	blockKey := fmt.Sprintf("login_blocked:%s", email)

	// ตรวจสอบว่าโดน block อยู่ไหม
	blockTTL, err := redis.Rdb.TTL(redis.Ctx, blockKey).Result()
	if err != nil && err != rdb.Nil {
		log.Println("Failed to check TTL of login_blocked:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// หากยังถูก block อยู่ ให้แจ้งกลับ
	if blockTTL > 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": fmt.Sprintf("Too many failed attempts. You are still blocked for %v", blockTTL.Round(time.Second)),
		})
		return
	}

	attempts, err := redis.Rdb.Get(redis.Ctx, attemptKey).Int()
	if err != nil && err != rdb.Nil {
		log.Println("Failed to get login_attempt from Redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	attempts++
	if err := redis.Rdb.Set(redis.Ctx, attemptKey, attempts, 15*time.Minute).Err(); err != nil {
		log.Println("Failed to set login_attempt in Redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if attempts > maxLoginAttempts {
		blockMultiplier := attempts - maxLoginAttempts
		blockTime := initialBlockTime * time.Duration(blockMultiplier)

		if err := redis.Rdb.Set(redis.Ctx, blockKey, "blocked", blockTime).Err(); err != nil {
			log.Println("Failed to set login_blocked in Redis:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": fmt.Sprintf("Too many failed attempts. You are blocked for %v", blockTime),
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var userData struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,password"`
	}
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	if err := validate.Struct(userData); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": customValidator.TranslateValidationError(err)})
		return
	}

	if h.isBlocked(c, userData.Email) {
		return
	}

	accessToken, refreshToken, err := h.authUsecase.Login(userData.Email, userData.Password)
	if err == nil {
		h.handleLoginSuccess(c, userData.Email, accessToken, refreshToken)
		return
	}

	h.handleLoginFailure(c, userData.Email)
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
