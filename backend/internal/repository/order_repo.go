package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
    return r.db.Create(order).Error
}

func (r *OrderRepository) GetUserOrders(userID uint) ([]models.Order, error) {
    var orders []models.Order
    err := r.db.Preload("Items.Product").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
    return orders, err
}