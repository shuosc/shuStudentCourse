package token

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

func StudentIdForToken(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	studentId := claims["studentId"].(string)
	return studentId
}

func GenerateJWT(studentId string) string {
	result, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"studentId": studentId,
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	return result
}
