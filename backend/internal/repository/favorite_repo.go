package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Add(userID, productID uint) error {
	favorite := &models.Favorite{
		UserID:    userID,
		ProductID: productID,
	}
	return r.db.Create(favorite).Error
}

func (r *FavoriteRepository) Remove(userID, productID uint) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.Favorite{}).Error
}

func (r *FavoriteRepository) GetUserFavorites(userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := r.db.Preload("Product").Where("user_id = ?", userID).Order("created_at desc").Find(&favorites).Error
	return favorites, err
}

func (r *FavoriteRepository) IsFavorite(userID, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error
	return count > 0, err
}