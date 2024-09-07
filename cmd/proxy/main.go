package main

import (
	"go-reverse-proxy/internal/config"
	"go-reverse-proxy/internal/handler"
	"go-reverse-proxy/internal/loadbalancer"
	"log"
	"net/http"
)

func main() {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	lb := loadbalancer.NewLoadBalancer(servers)

	cfg, err := config.LoadConfig("config/config.json")

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	proxyHandler := handler.NewProxyHandler(lb)

	log.Printf("Starting reverse proxy on port %s", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, proxyHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
