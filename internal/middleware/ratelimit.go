package middleware

import (
	"api-gateway/internal/config"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var visitors = map[string]*rate.Limiter{}
var mu sync.Mutex

func getClientLimiter(ip string, cfg config.Config) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if visitor, exists := visitors[ip]; exists {
		return visitor
	}

	// Create a new limiter.
	limiter := rate.NewLimiter(rate.Every(time.Minute), cfg.Server.RpmLimit)
	visitors[ip] = limiter
	return limiter
}

func RateLimit(next http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getClientLimiter(ip, cfg)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(int(limiter.Tokens())))

		// Check if the request is allowed
		if !limiter.Allow() {
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			respObj := map[string]interface{}{
				"error": "Too Many Requests",
				"data":  nil,
			}
			respBytes, _ := json.MarshalIndent(respObj, "", "  ")
			resp := string(respBytes)
			w.Write([]byte(resp))
			return
		}

		next.ServeHTTP(w, r)

	})
}
