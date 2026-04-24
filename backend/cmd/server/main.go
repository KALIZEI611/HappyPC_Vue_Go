package main

import (
	"log"
	"net/http"
	"os"

	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
    db, err := database.ConnectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    database.MigrateAndSeed(db)

    categoryRepo := repository.NewCategoryRepository(db)
    productRepo := repository.NewProductRepository(db)
    userRepo := repository.NewUserRepository(db)
    sessionRepo := repository.NewSessionRepository(db)
    cartRepo := repository.NewCartRepository(db)

    categoryHandler := handlers.NewCategoryHandler(categoryRepo, productRepo)
    cartHandler := handlers.NewCartHandler(db)
    authHandler := handlers.NewAuthHandler(userRepo, sessionRepo)
    orderHandler := handlers.NewOrderHandler(db, cartRepo)

    buildRepo := repository.NewBuildRepository(db)
    buildHandler := handlers.NewBuildHandler(buildRepo, db)

    favoriteRepo := repository.NewFavoriteRepository(db)
    favoriteHandler := handlers.NewFavoriteHandler(favoriteRepo)
    
    feedbackHandler := handlers.NewFeedbackHandler(db)

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{
            "http://localhost",
            "http://localhost:5173",
            "http://127.0.0.1",
            "http://127.0.0.1:5173",
            "https://happy-pc-vue-6enowbewt-kalizei611s-projects.vercel.app", // ← Добавьте ЭТОТ URL
            "https://happy-pc-vue-*.vercel.app", // ← Можно добавить wildcard для всех preview
        },
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        AllowCredentials: true,
    }))

    r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    r.Get("/categories", categoryHandler.GetAll)
    r.Get("/{urlKey}", categoryHandler.GetProductsByURLKey)
    r.Get("/category/{id}/products", categoryHandler.GetCategoryByID)

    r.Route("/api", func(r chi.Router) {
        r.Post("/register", authHandler.Register)
        r.Post("/login", authHandler.Login)
        r.Post("/logout", authHandler.Logout)
        r.Get("/me", authHandler.Me)

        r.Group(func(r chi.Router) {
            r.Use(handlers.AuthMiddleware(sessionRepo, userRepo))
            r.Get("/cart", cartHandler.GetCart)
            r.Post("/cart", cartHandler.AddToCart)
            r.Put("/cart/{productId}", cartHandler.UpdateCartItem)
            r.Delete("/cart/{productId}", cartHandler.RemoveFromCart)
            r.Delete("/cart", orderHandler.ClearCart)
            r.Get("/orders", orderHandler.GetUserOrders)
            r.Post("/orders", orderHandler.CreateOrder)
            r.Post("/builds", buildHandler.SaveBuild)
            r.Get("/builds", buildHandler.GetUserBuilds)
            r.Delete("/builds/{id}", buildHandler.DeleteBuild)
            r.Get("/favorites", favoriteHandler.GetFavorites)
            r.Post("/favorites", favoriteHandler.AddFavorite)
            r.Delete("/favorites/{id}", favoriteHandler.RemoveFavorite)
            r.Post("/favorites/toggle", favoriteHandler.ToggleFavorite)
            r.Post("/feedback", feedbackHandler.SendFeedback)
            r.Put("/user/profile", authHandler.UpdateProfile)
        })
    })

    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on port %s", port)
    http.ListenAndServe(":"+port, r)
}