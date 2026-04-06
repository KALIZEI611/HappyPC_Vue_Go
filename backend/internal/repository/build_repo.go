package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type BuildRepository struct {
    db *gorm.DB
}

func NewBuildRepository(db *gorm.DB) *BuildRepository {
    return &BuildRepository{db: db}
}

func (r *BuildRepository) Create(build *models.Build) error {
    return r.db.Create(build).Error
}

func (r *BuildRepository) GetUserBuilds(userID uint) ([]models.Build, error) {
    var builds []models.Build
    err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&builds).Error
    return builds, err
}

func (r *BuildRepository) Delete(id uint) error {
    return r.db.Delete(&models.Build{}, id).Error
}