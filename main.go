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
	"github.com/rs/cors"
	"github.com/yosepalexsander/waysbucks-api/db"
	"github.com/yosepalexsander/waysbucks-api/handler"
	customMiddleware "github.com/yosepalexsander/waysbucks-api/handler/middleware"
	"github.com/yosepalexsander/waysbucks-api/persistance"
)

type Env struct {
	user handler.UserServer
	address handler.AddressServer
} 

func main()  {
	var dbStore db.DBStore
	db.Connect(&dbStore)
	env := Env{
		user: handler.UserServer{
			Repo: persistance.UserRepo{DB: dbStore.DB},
		},
		address: handler.AddressServer{
			Repo: persistance.AddressRepo{DB: dbStore.DB},
		},
	}

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

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", env.user.Register)
		r.Post("/login", env.user.Login)

		r.Route("/users", func(r chi.Router) {
			r.Use(customMiddleware.Authentication)
			r.Get("/", env.user.GetUsers)
			r.Get("/{userID}", env.user.GetUser)
			r.Put("/{userID}", env.user.UpdateUser)
			r.Delete("/{userID}", env.user.DeleteUser)
		})
		r.Route("/address", func(r chi.Router) {
			r.Use(customMiddleware.Authentication)
			r.Get("/", env.address.GetUserAddress)
			r.Post("/", env.address.CreateUserAddress)
			r.Put("/{addressID}", env.address.UpdateUserAddress)
			r.Delete("/{addressID}", env.address.DeleteAddress)
		})

	})
	
	server := &http.Server{
		Addr: "0.0.0.0:8080", 
		Handler: r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 30,
	}
	log.Printf("Server Started")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
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