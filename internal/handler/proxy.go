package handler

import (
	"go-reverse-proxy/internal/loadbalancer"
	"io"
	"net/http"
)

func NewProxyHandler(lb *loadbalancer.LoadBalancer) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		targetURL := lb.NextServer() + request.RequestURI
		proxyRequest, err := http.NewRequest(request.Method, targetURL, request.Body)

		if err != nil {
			http.Error(responseWriter, "Could not proxy to backend", http.StatusInternalServerError)
		}

		proxyRequest.Header = request.Header

		response, err := http.DefaultClient.Do(proxyRequest)

		if err != nil {
			http.Error(responseWriter, "Could not proxy to backend", http.StatusInternalServerError)
		}

		defer request.Body.Close()

		for name, values := range response.Header {
			for _, value := range values {
				responseWriter.Header().Add(name, value)
			}
		}
		responseWriter.WriteHeader(response.StatusCode)
		_, err = io.Copy(responseWriter, response.Body)

		if err != nil {
			http.Error(responseWriter, "Error copying response to client", http.StatusInternalServerError)
		}
	}
}
