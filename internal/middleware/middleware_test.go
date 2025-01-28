package middleware

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"log/slog"
)

func setupRouter(log *slog.Logger, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.Use(WithJWTAuth(log, jwtSecret))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	return r
}

func TestWithJWTAuth(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	jwtSecret := "my_secret_key"

	// Helper function to create a JWT token string
	createToken := func(claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(jwtSecret))
		return tokenString
	}

	t.Run("No Token", func(t *testing.T) {
		router := setupRouter(log, jwtSecret)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		router := setupRouter(log, jwtSecret)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer invalid_token")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Valid Token", func(t *testing.T) {
		router := setupRouter(log, jwtSecret)
		w := httptest.NewRecorder()

		tokenString := createToken(jwt.MapClaims{"uid": "12345"})
		req, _ := http.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{Name: "Authorization", Value: tokenString})
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("Token Query Param Valid", func(t *testing.T) {
		router := setupRouter(log, jwtSecret)
		w := httptest.NewRecorder()

		tokenString := createToken(jwt.MapClaims{"uid": "12345"})
		req, _ := http.NewRequest("GET", "/test?token="+tokenString, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("Invalid Signing Method", func(t *testing.T) {
		router := setupRouter(log, jwtSecret)
		w := httptest.NewRecorder()

		// Create a token with a different signing method
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"uid": "12345"})
		tokenString, _ := token.SignedString([]byte(jwtSecret))
		req, _ := http.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{Name: "Authorization", Value: tokenString})
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
