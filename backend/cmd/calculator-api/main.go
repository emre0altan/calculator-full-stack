package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"calculator/backend/internal/httpapi"
)

var listenAndServe = (*http.Server).ListenAndServe
var fatalf = log.Fatalf

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           newHandler(os.Getenv("WEB_DIR")),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("calculator listening on port %s", port)
	if err := listenAndServe(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fatalf("serve calculator API: %v", err)
	}
}

func newHandler(webDir string) http.Handler {
	if webDir == "" {
		return httpapi.NewHandler(nil)
	}
	return httpapi.NewHandler(http.FileServer(http.Dir(webDir)))
}
