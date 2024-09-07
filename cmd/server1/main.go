package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from the target server0!")
	log.Printf("Request received")
}

func main() {
	http.HandleFunc("/", handler)
	port := ":8082"
	log.Printf("Target server0 listening on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
