package router

import (
	"github.com/go-chi/chi/v5"
	customMiddleware "github.com/yosepalexsander/waysbucks-api/handler/middleware"
	"github.com/yosepalexsander/waysbucks-api/interactor"
)

func NewRouter(r *chi.Mux, h *interactor.AppHandler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		
		r.Route("/users", func(r chi.Router) {
			r.Use(customMiddleware.Authentication)
			r.Get("/{userID}", h.GetUser)
			r.Put("/{userID}", h.UpdateUser)
			r.Delete("/{userID}", h.DeleteUser)

			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.AdminOnly)
				r.Get("/", h.GetUsers)
			})
		})
		r.Route("/address", func(r chi.Router) {
			r.Use(customMiddleware.Authentication)
			r.Get("/", h.GetUserAddress)
			r.Post("/", h.CreateAddress)
			r.Put("/{addressID}", h.UpdateAddress)
			r.Delete("/{addressID}", h.DeleteAddress)
		})
		r.Route("/products", func(r chi.Router) {
			r.Get("/", h.GetProducts)
			r.Get("/{productID}", h.GetProduct)
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.Authentication)
				r.Use(customMiddleware.AdminOnly)
				r.Post("/", h.CreateProduct)
				r.Put("/{productID}", h.UpdateProduct)
				r.Delete("/{productID}", h.DeleteProduct)
			})
		})
		r.Route("/toppings", func(r chi.Router) {
			r.Get("/", h.GetToppings)

			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.Authentication)
				r.Use(customMiddleware.AdminOnly)
				r.Post("/", h.CreateTopping)
				r.Put("/{toppingID}", h.UpdateTopping)
				r.Delete("/{toppingID}", h.DeleteTopping)
			})
			
		})
	})
}