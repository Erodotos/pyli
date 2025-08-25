package router

import (
	"api-gateway/internal/config"
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
	for _, route := range config.Routes {
		log.Printf("Registering route: %s", route.Path)
		proxy := proxy.NewProxy(route.Endpoint)
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
	}
}
