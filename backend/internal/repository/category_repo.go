package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
    db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
    return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
    var categories []models.Category
    err := r.db.Find(&categories).Error
    return categories, err
}

func (r *CategoryRepository) GetByID(id uint) (*models.Category, error) {
    var category models.Category
    err := r.db.First(&category, id).Error
    return &category, err
}

func (r *CategoryRepository) GetByURLKey(urlKey string) (*models.Category, error) {
    var category models.Category
    err := r.db.Where("url_key = ?", urlKey).First(&category).Error
    return &category, err
}