package main

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/yosepalexsander/waysbucks-api/db"
	"github.com/yosepalexsander/waysbucks-api/interactor"
	"github.com/yosepalexsander/waysbucks-api/router"
)

func main()  {
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
		MaxAge: 300,
	}).Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.NewRouter(r, appHandler)
	
	server := &http.Server{
		Addr: "127.0.0.1:8080", 
		Handler: r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 40,
		IdleTimeout:  time.Second * 40,
	}
	log.Printf("Server Started")
	
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	gracefullShutdown(server)
}

func gracefullShutdown(server *http.Server) {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	
	// Listen for syscall signals for process to interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<- sig
	
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer func() {
			signal.Stop(sig)
			cancel()
		}()

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
		serverStopCtx()
	}()

	// Wait for server context to be stopped
	<- serverCtx.Done()
}
