package jwtutil

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateTokens(userID uint) (accessToken string, refreshToken string, err error) {
	accessTokenClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessToken, err = access.SignedString(jwtSecret)
	if err != nil {
		return
	}

	refreshTokenClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
    refreshToken, err = refresh.SignedString(jwtSecret)
	if err != nil {
        return
    }
    return
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
}