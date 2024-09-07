package handler

import (
	"io"
	"net/http"
	"net/url"
)

func NewProxyHandler(targetURL *url.URL) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		proxyRequest, err := http.NewRequest(request.Method, targetURL.String()+request.RequestURI, request.Body)

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
