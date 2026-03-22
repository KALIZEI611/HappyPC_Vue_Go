package main

import (
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
    cfg := config.Load()
    db := database.InitDB(cfg)

    categoryRepo := repository.NewCategoryRepository(db)
    productRepo := repository.NewProductRepository(db)

    categoryHandler := handlers.NewCategoryHandler(categoryRepo, productRepo)
    cartHandler := handlers.NewCartHandler(db)

    userRepo := repository.NewUserRepository(db)
    sessionRepo := repository.NewSessionRepository(db)
    authHandler := handlers.NewAuthHandler(userRepo, sessionRepo)

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{
            "http://localhost",
            "http://localhost:5173",
            "http://127.0.0.1",
            "http://127.0.0.1:5173",
        },
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        AllowCredentials: true,
    }))

    r.Get("/categories", categoryHandler.GetAll)
    r.Get("/{urlKey}", categoryHandler.GetProductsByURLKey)
    r.Get("/category/{id}/products", categoryHandler.GetCategoryByID)

    r.Get("/cart", cartHandler.GetCart)
    r.Post("/cart", cartHandler.AddToCart)
    r.Put("/cart/{productId}", cartHandler.UpdateCartItem)
    r.Delete("/cart/{productId}", cartHandler.RemoveFromCart)

    r.Post("/api/register", authHandler.Register)
    r.Post("/api/login", authHandler.Login)
    r.Post("/api/logout", authHandler.Logout)
    r.Get("/api/me", authHandler.Me)

    r.Group(func(r chi.Router) {
        r.Use(handlers.AuthMiddleware(sessionRepo, userRepo))
        r.Get("/api/profile", func(w http.ResponseWriter, r *http.Request) {
            user := handlers.GetUserFromContext(r.Context())
            // просто заглушка, позже можно реализовать
            _ = user
        })
    })

    port := cfg.ServerPort
    log.Printf("Server starting on port %s", port)
    http.ListenAndServe(":"+port, r)
}