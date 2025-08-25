package proxy

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
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

			incommingPath := req.URL.Path

			var path string
			re := regexp.MustCompile(`^/api/[a-zA-Z]+/(.*)$`)
			matches := re.FindStringSubmatch(req.URL.Path)
			if len(matches) > 1 {
				path = "/" + matches[1]
			} else {
				path = ""
			}

			req.URL.Scheme = endpointURL.Scheme
			req.URL.Host = endpointURL.Host
			req.Host = endpointURL.Host
			req.URL.Path = path

			// Log outgoing request details
			log.Printf("[Forwarding] %s %s --> %s://%s%s", req.Method, incommingPath, req.URL.Scheme, req.URL.Host, req.URL.Path)
		},
		ModifyResponse: func(resp *http.Response) error {
			// 🔎 Remove backend rate-limit headers
			resp.Header.Del("X-RateLimit-Limit")
			resp.Header.Del("X-RateLimit-Remaining")
			resp.Header.Del("X-RateLimit-Reset")
			return nil
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
