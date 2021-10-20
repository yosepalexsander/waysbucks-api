package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/yosepalexsander/waysbucks-api/db"
	"github.com/yosepalexsander/waysbucks-api/interactor"
	"github.com/yosepalexsander/waysbucks-api/router"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	var dbStore db.DBStore
	db.Connect(&dbStore)
	interactor := interactor.Interactor{DB: dbStore.DB}
	appHandler := interactor.NewAppHandler()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{"*"},
		MaxAge:         300,
	}).Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	router.NewRouter(r, appHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 30,
	}

	log.Println("Server Started on port 8080")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	gracefullShutdown(server)
}

func gracefullShutdown(server *http.Server) {

	// Listen for syscall signals for process to interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig

	// Shutdown signal with grace period of 30 seconds
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
		signal.Stop(sig)
	}()
	go func() {
		// Trigger graceful shutdown
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}

		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Fatalf("graceful shutdown timed out")
		}
	}()

	// Wait for server context to be stopped
	log.Println("shutting down")
	os.Exit(0)
}
