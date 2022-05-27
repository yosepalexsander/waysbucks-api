package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/yosepalexsander/waysbucks-api/handler"
	"github.com/yosepalexsander/waysbucks-api/interactor"
	customMiddleware "github.com/yosepalexsander/waysbucks-api/middleware"
)

func NewRouter(r *chi.Mux, h *interactor.AppHandler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Register)
			r.Post("/login", h.Login)

			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.Authentication)
				r.Get("/profile", h.GetUser)
				r.Put("/profile", h.UpdateUser)
				r.Delete("/profile", h.DeleteUser)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(customMiddleware.AdminOnly)
			r.Get("/", h.GetUsers)
		})

		r.Route("/address", func(r chi.Router) {
			r.Use(customMiddleware.Authentication)
			r.Get("/", h.GetUserAddresses)
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

		r.Route("/carts", func(r chi.Router) {
			r.Use(customMiddleware.Authentication)
			r.Get("/", h.GetCarts)
			r.Post("/", h.CreateCart)
			r.Put("/{cartID}", h.UpdateCart)
			r.Delete("/{cartID}", h.DeleteCart)
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Use(customMiddleware.Authentication)
			r.Post("/", h.CreateTransaction)
			r.Get("/{transactionID}", h.GetTransaction)
			r.With(customMiddleware.AdminOnly).Get("/", h.GetTransactions)
		})

		r.Route("/user-transactions", func(r chi.Router) {
			r.Use(customMiddleware.Authentication)
			r.Get("/", h.GetUserTransactions)
		})

		r.Post("/notification", h.PaymentNotification)

		r.Route("/upload", func(r chi.Router) {
			r.Use(customMiddleware.Authentication)
			r.Post("/", handler.UploadImage)
			r.Post("/avatar", handler.UploadAvatar)
		})
	})
}
