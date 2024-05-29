package util

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
)

var JWTSecret = "111"

type UserClaims struct {
	ID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT TODO: 这里肯定不可能只以 user_id 签发 token 的，之后再改
func GenerateJWT(userID uint) string {
	claims := UserClaims{
		ID:               userID,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		slog.Warn("GenerateJWT error: ", err)
		return ""
	}
	return str
}

func ParseJWT(str string) (id uint, err error) {
	token, err := jwt.ParseWithClaims(str, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims.ID, nil
	} else {
		return 0, err
	}
}
