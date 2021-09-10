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
	"github.com/yosepalexsander/waysbucks-api/db"
)

func main()  {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	
	var dbEnv db.DB
	db.Connect(&dbEnv)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome"))
		})
	})
	
	server := http.Server{
		Addr: "0.0.0.0:8080", 
		Handler: r,
	}
	log.Printf("Server Started on: %s", server.Addr)
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	
	// Listen for syscall signals for process to interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<- sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()
		go func() {
			<- shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatalf("graceful shutdown timed out")
			}
		}()
		
		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)

		}

		log.Print("Server Stopped")
		serverStopCtx()
	}()

	serverErr := server.ListenAndServe()
	if serverErr != nil  && serverErr != http.ErrServerClosed{
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<- serverCtx.Done()
}