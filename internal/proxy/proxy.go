package proxy

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func NewProxy(serviceEndpointMap map[string]string) *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			trimmed := strings.TrimPrefix(req.URL.Path, "/api/")
			parts := strings.SplitN(trimmed, "/", 2)
			path := ""
			if len(parts) > 1 {
				path = "/" + parts[1]
			}

			req.URL.Scheme = "https"
			req.URL.Host = serviceEndpointMap["/api/"+parts[0]+"/"]
			req.Host = serviceEndpointMap["/api/"+parts[0]+"/"]
			req.URL.Path = path

			// 🔎 Log outgoing request details
			log.Printf("[Forwarding] %s %s://%s%s", req.Method, req.URL.Scheme, req.URL.Host, req.URL.Path)
			for name, values := range req.Header {
				for _, v := range values {
					log.Printf("  %s: %s", name, v)
				}
			}
		},
		ErrorHandler: func(w http.ResponseWriter, req *http.Request, err error) {
			log.Printf("[Error] %v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)

			respObj := map[string]interface{}{
				"error": "Service unavailable",
				"data":  nil,
			}
			respBytes, _ := json.MarshalIndent(respObj, "", "  ")
			resp := string(respBytes)
			w.Write([]byte(resp))
		},
	}
	return proxy
}
