package proxy

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewProxy(serviceEndpoint string) *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// Parse the serviceEndpoint into scheme, host, and path
			endpointURL, err := url.Parse(serviceEndpoint)
			if err != nil {
				log.Printf("Error parsing service endpoint %s: %v. Fallback to http://localhost", serviceEndpoint, err)
				endpointURL = &url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   "",
				}
			}

			req.URL.Scheme = endpointURL.Scheme
			req.URL.Host = endpointURL.Host
			req.Host = endpointURL.Host

			// Remove "/api" prefix from incoming request path, then join with basePath
			path := strings.TrimPrefix(req.URL.Path, "/api")
			// Ensure no double slashes when joining paths
			if !strings.HasSuffix(endpointURL.Path, "/") && !strings.HasPrefix(path, "/") {
				req.URL.Path = endpointURL.Path + "/" + path
			} else {
				req.URL.Path = endpointURL.Path + path
			}

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
