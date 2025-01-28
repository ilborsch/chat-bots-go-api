package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"net/http"
)

func WithJWTAuth(log *slog.Logger, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := getRequestToken(c)
		if tokenString == "" {
			log.Error("auth token is empty")
			permissionDenied(c)
		}

		token, err := validateJWT(tokenString, jwtSecret)
		if err != nil {
			log.Error("invalid auth token")
			permissionDenied(c)
			return
		}

		if !token.Valid {
			log.Error("invalid auth token")
			permissionDenied(c)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["uid"].(float64)

		c.Set("uid", userID)
		c.Next()
	}
}

func getRequestToken(c *gin.Context) string {
	tokenQuery := c.Query("token")
	if tokenQuery != "" {
		return tokenQuery
	}

	tokenHeader, err := c.Cookie("Authorization")
	if err != nil {
		return ""
	}
	return tokenHeader
}

func validateJWT(tokenString string, jwtSecret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
}

func permissionDenied(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"Error": "Permission denied",
	})
}
