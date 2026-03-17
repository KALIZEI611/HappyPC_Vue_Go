package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

func (h *CartHandler) getOrCreateCartID(w http.ResponseWriter, r *http.Request) string {
    cookie, err := r.Cookie("cart_id")
    if err == nil && cookie.Value != "" {
        return cookie.Value
    }
    newID := fmt.Sprintf("cart_%d", os.Getpid())
    http.SetCookie(w, &http.Cookie{
        Name:  "cart_id",
        Value: newID,
        Path:  "/",
    })
    return newID
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
    cartID := h.getOrCreateCartID(w, r)
    var items []models.CartItem
    if err := h.db.Preload("Product").Where("cart_id = ?", cartID).Find(&items).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}

type addToCartRequest struct {
    ProductID uint `json:"product_id"`
    Quantity  int  `json:"quantity"`
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
    cartID := h.getOrCreateCartID(w, r)
    var req addToCartRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    if req.Quantity <= 0 {
        req.Quantity = 1
    }
    var existing models.CartItem
    err := h.db.Where("cart_id = ? AND product_id = ?", cartID, req.ProductID).First(&existing).Error
    if err == nil {
        existing.Quantity += req.Quantity
        h.db.Save(&existing)
    } else {
        item := models.CartItem{
            CartID:    cartID,
            ProductID: req.ProductID,
            Quantity:  req.Quantity,
        }
        h.db.Create(&item)
    }
    w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
    cartID := h.getOrCreateCartID(w, r)
    productIDParam := chi.URLParam(r, "productId")
    var productID uint
    fmt.Sscanf(productIDParam, "%d", &productID)

    var req struct {
        Quantity int `json:"quantity"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if req.Quantity <= 0 {
        h.db.Where("cart_id = ? AND product_id = ?", cartID, productID).Delete(&models.CartItem{})
    } else {
        h.db.Model(&models.CartItem{}).Where("cart_id = ? AND product_id = ?", cartID, productID).Update("quantity", req.Quantity)
    }
    w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
    cartID := h.getOrCreateCartID(w, r)
    productIDParam := chi.URLParam(r, "productId")
    var productID uint
    fmt.Sscanf(productIDParam, "%d", &productID)

    h.db.Where("cart_id = ? AND product_id = ?", cartID, productID).Delete(&models.CartItem{})
    w.WriteHeader(http.StatusOK)
}