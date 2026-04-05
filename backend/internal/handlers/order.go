package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/internal/models"
	"backend/internal/repository"

	"gorm.io/gorm"
)

type OrderHandler struct {
    db        *gorm.DB
    cartRepo  *repository.CartRepository
    orderRepo *repository.OrderRepository
}

func NewOrderHandler(db *gorm.DB, cartRepo *repository.CartRepository) *OrderHandler {
    return &OrderHandler{
        db:        db,
        cartRepo:  cartRepo,
        orderRepo: repository.NewOrderRepository(db),
    }
}

type createOrderRequest struct {
    Items           []struct {
        ProductID uint    `json:"product_id"`
        Quantity  int     `json:"quantity"`
        Price     float64 `json:"price"`
    } `json:"items"`
    DeliveryMethod  string  `json:"delivery_method"`
    DeliveryAddress string  `json:"delivery_address"`
    PaymentMethod   string  `json:"payment_method"`
    TotalPrice      float64 `json:"total_price"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    user := GetUserFromContext(r.Context())
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req createOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Decode error:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("Received order request: %+v", req)

    if len(req.Items) == 0 {
        http.Error(w, "Order must contain at least one item", http.StatusBadRequest)
        return
    }

    order := &models.Order{
        UserID:          user.ID,
        DeliveryMethod:  req.DeliveryMethod,
        DeliveryAddress: req.DeliveryAddress,
        PaymentMethod:   req.PaymentMethod,
        TotalPrice:      req.TotalPrice,
        Status:          "pending",
    }

    if err := h.orderRepo.Create(order); err != nil {
        http.Error(w, "Failed to create order", http.StatusInternalServerError)
        return
    }

    for _, item := range req.Items {
        orderItem := &models.OrderItem{
            OrderID:   order.ID,
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            Price:     item.Price,
        }
        if err := h.db.Create(orderItem).Error; err != nil {
            http.Error(w, "Failed to create order item", http.StatusInternalServerError)
            return
        }
    }

    if err := h.cartRepo.ClearCart(user.ID); err != nil {
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":      order.ID,
        "message": "Order created successfully",
    })
}

func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
    user := GetUserFromContext(r.Context())
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var orders []models.Order
    err := h.db.Preload("Items.Product").Where("user_id = ?", user.ID).Order("created_at desc").Find(&orders).Error
    if err != nil {
        http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
    user := GetUserFromContext(r.Context())
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    if err := h.cartRepo.ClearCart(user.ID); err != nil {
        http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}