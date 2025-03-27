package main

import (
	"fmt"
	"log"
	"net/http"
	"go-compress-showdown/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running")
	})
	mux.HandleFunc("/gzip", handler.GzipHandler) // Add GzipHandler route
	mux.HandleFunc("/none", handler.NoneHandler) // Add NoneHandler route
	mux.HandleFunc("/brotli", handler.BrotliHandler) // Add BrotliHandler route
	mux.HandleFunc("/zstd", handler.ZstdHandler) // Add ZstdHandler route

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}