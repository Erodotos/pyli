package middleware

// LoginService to provide user login with JWT token support
import (
	"api-gateway/internal/config"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": "public-access",
			"exp":      time.Now().Add(time.Hour).Unix(),
		})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(next http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// This middleware does not apply on /login route
		// hence we skip
		if r.URL.Path == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := r.Header.Get("X-Access-Token")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorised", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

	})
}
