package router

import (
	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"encoding/json"
	"log"
	"net/http"

	"api-gateway/internal/proxy"
)

func NewRouter(cfg config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	RegisterRoutes(mux, cfg)
	return mux
}

func RegisterRoutes(mux *http.ServeMux, config config.Config) {

	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Generate token
		jwt, err := middleware.GenerateToken(config.JwtSecret)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Prepare response
		response := map[string]interface{}{
			"access_token": jwt,
			"token_type":   "Bearer",
		}

		json.NewEncoder(w).Encode(response)
	})

	// Register custom routes provided in the config.yaml
	for _, route := range config.Routes {
		log.Printf("Registering route: %s", route.Path)
		proxy := proxy.NewProxy(route.Endpoint)
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.Header().Set("Allow", route.Method)
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			proxy.ServeHTTP(w, r)
		})
	}
}
