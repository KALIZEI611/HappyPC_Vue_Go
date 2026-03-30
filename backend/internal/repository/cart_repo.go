package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type CartRepository struct {
    db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
    return &CartRepository{db: db}
}

func (r *CartRepository) ClearCart(userID uint) error {
    return r.db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}