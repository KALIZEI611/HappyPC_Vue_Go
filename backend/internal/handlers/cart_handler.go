package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"backend/internal/models"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type CartHandler struct {
    db *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
    return &CartHandler{db: db}
}

func (h *CartHandler) getUserIDFromContext(r *http.Request) (uint, error) {
    user := GetUserFromContext(r.Context())
    if user == nil {
        return 0, fmt.Errorf("user not authenticated")
    }
    return user.ID, nil
}

type addToCartRequest struct {
    ProductID uint `json:"product_id"`
    Quantity  int  `json:"quantity"`
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
    userID, err := h.getUserIDFromContext(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    var items []models.CartItem
    if err := h.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
    userID, err := h.getUserIDFromContext(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    var req addToCartRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    if req.Quantity <= 0 {
        req.Quantity = 1
    }
    var existing models.CartItem
    err = h.db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existing).Error
    if err == nil {
        existing.Quantity += req.Quantity
        h.db.Save(&existing)
    } else {
        item := models.CartItem{
            UserID:    userID,
            ProductID: req.ProductID,
            Quantity:  req.Quantity,
        }
        h.db.Create(&item)
    }
    w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
    userID, err := h.getUserIDFromContext(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    productIDParam := chi.URLParam(r, "productId")
    productID, err := strconv.Atoi(productIDParam)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    var req struct {
        Quantity int `json:"quantity"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if req.Quantity <= 0 {
        h.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.CartItem{})
    } else {
        h.db.Model(&models.CartItem{}).Where("user_id = ? AND product_id = ?", userID, productID).Update("quantity", req.Quantity)
    }
    w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
    userID, err := h.getUserIDFromContext(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    productIDParam := chi.URLParam(r, "productId")
    productID, err := strconv.Atoi(productIDParam)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    h.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.CartItem{})
    w.WriteHeader(http.StatusOK)
}