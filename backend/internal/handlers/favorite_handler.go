package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

type FavoriteHandler struct {
	favoriteRepo *repository.FavoriteRepository
}

func NewFavoriteHandler(favoriteRepo *repository.FavoriteRepository) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteRepo: favoriteRepo,
	}
}

func (h *FavoriteHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	favorites, err := h.favoriteRepo.GetUserFavorites(user.ID)
	if err != nil {
		http.Error(w, "Failed to get favorites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorites)
}

func (h *FavoriteHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ProductID uint `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ProductID == 0 {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	if err := h.favoriteRepo.Add(user.ID, req.ProductID); err != nil {
		http.Error(w, "Failed to add favorite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product added to favorites",
	})
}

func (h *FavoriteHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.favoriteRepo.Remove(user.ID, uint(productID)); err != nil {
		http.Error(w, "Failed to remove favorite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product removed from favorites",
	})
}

func (h *FavoriteHandler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ProductID uint `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ProductID == 0 {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	isFav, err := h.favoriteRepo.IsFavorite(user.ID, req.ProductID)
	if err != nil {
		http.Error(w, "Failed to check favorite", http.StatusInternalServerError)
		return
	}

	if isFav {
		err = h.favoriteRepo.Remove(user.ID, req.ProductID)
	} else {
		err = h.favoriteRepo.Add(user.ID, req.ProductID)
	}

	if err != nil {
		http.Error(w, "Failed to toggle favorite", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"isFavorite": !isFav})
}