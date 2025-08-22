package router

import (
	"api-gateway/internal/config"
	"fmt"
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
	serviceEndpointMap := map[string]string{}
	for _, route := range config.Routes {
		serviceEndpointMap[route.Path] = route.Endpoint
	}
	proxy := proxy.NewProxy(serviceEndpointMap)
	for _, route := range config.Routes {
		log.Printf("Registering route: %s", route.Path)
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(route.Path)
			proxy.ServeHTTP(w, r)
		})
	}
}
