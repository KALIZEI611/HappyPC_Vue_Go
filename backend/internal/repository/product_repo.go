package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
    return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByCategoryID(categoryID uint) ([]models.Product, error) {
    var products []models.Product
    err := r.db.Where("category_id = ?", categoryID).Find(&products).Error
    return products, err
}

func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
    var product models.Product
    err := r.db.First(&product, id).Error
    return &product, err
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
    var products []models.Product
    err := r.db.Find(&products).Error
    return products, err
}