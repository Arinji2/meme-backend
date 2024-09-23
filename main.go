package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	custom_log "github.com/Arinji2/meme-backend/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
)

func main() {
	r := chi.NewRouter()
	r.Use(SkipLoggingMiddleware)

	err := godotenv.Load()
	if err != nil {
		isProduction := os.Getenv("ENVIRONMENT") == "PRODUCTION"
		if !isProduction {
			log.Fatal("Error loading .env file")
		} else {
			custom_log.Logger.Warn("Using Production Environment")
		}
	} else {
		custom_log.Logger.Warn("Using Development Environment")
	}

	r.Get("/", healthHandler)

	r.Get("/health", healthCheckHandler)

	http.ListenAndServe(":8080", r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Meme Backend: Request Received")
	w.Write([]byte("Meme Backend: Request Received"))
	key := r.URL.Query().Get("key")

	if key != "" && key != os.Getenv("ACCESS_KEY") {
		render.Status(r, http.StatusUnauthorized)
		return
	}

	render.Status(r, http.StatusOK)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Meme Backend: Health Check"))
	render.Status(r, http.StatusOK)
}

func SkipLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		middleware.Logger(next).ServeHTTP(w, r)
	})
}
