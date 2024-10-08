package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	custom_log "github.com/Arinji2/meme-backend/logger"
	"github.com/Arinji2/meme-backend/routes/oauth"
	"github.com/Arinji2/meme-backend/routes/user"
	sql_db "github.com/Arinji2/meme-backend/sql"
	session_dal "github.com/Arinji2/meme-backend/sql/dal/session"
	"github.com/Arinji2/meme-backend/types"
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

	r.Group((func(r chi.Router) {
		r.Use(DBMiddleware)
		r.Use(AuthMiddleware)
		r.Route("/users", (func(r chi.Router) {
			r.Get("/by-email", user.GetUserByEmailHandler)
		}))
	}))

	r.Route("/oauth2-redirect", (func(r chi.Router) {
		r.Use(DBMiddleware)
		r.Get("/google", oauth.RegisterWithGoogleOauth)
	}))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
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

func DBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := sql_db.GetDB()
		ctx := context.WithValue(r.Context(), types.DBKey, db)
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("X-Session-ID")

		session, err := session_dal.GetUserBySession(r.Context(), sessionID)
		if err != nil {

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), types.UserIDKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
