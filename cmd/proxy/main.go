package main

import (
	"go-reverse-proxy/internal/config"
	"go-reverse-proxy/internal/handler"
	"log"
	"net/http"
	"net/url"
)

func main() {
	cfg, err := config.LoadConfig("config/config.json")

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	targetUrl, err := url.Parse(cfg.TargetURL)

	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}

	proxyHandler := handler.NewProxyHandler(targetUrl)

	log.Printf("Starting reverse proxy on port %s", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, proxyHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
