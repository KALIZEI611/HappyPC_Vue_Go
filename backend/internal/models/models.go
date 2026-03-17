package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Category struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    CreatedAt time.Time      `json:"-"`
    UpdatedAt time.Time      `json:"-"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    Name      string         `json:"name" gorm:"not null"`
    URLKey    string         `json:"url_key" gorm:"uniqueIndex;not null"`
    IconURL   string         `json:"icon_url"`
    Products  []Product      `json:"products,omitempty"`
}

type Product struct {
    ID          uint           `json:"id" gorm:"primarykey"`
    CreatedAt   time.Time      `json:"-"`
    UpdatedAt   time.Time      `json:"-"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    Name        string         `json:"name" gorm:"not null"`
    Price       float64        `json:"price" gorm:"not null"`
    Image       string         `json:"image"`
    Rating      float32        `json:"rating" gorm:"default:4.5"`
    Description string         `json:"description"`
    Specs       json.RawMessage `json:"specs" gorm:"type:jsonb"`
    CategoryID  uint           `json:"category_id"`
    Category    Category       `json:"category,omitempty"`
}

type CartItem struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    CreatedAt time.Time      `json:"-"`
    UpdatedAt time.Time      `json:"-"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    CartID    string         `json:"cart_id" gorm:"index"`
    ProductID uint           `json:"product_id"`
    Product   Product        `json:"product"`
    Quantity  int            `json:"quantity" gorm:"not null;default:1"`
}