package repository

import (
	"backend/internal/models"
	"backend/internal/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    err := r.db.First(&user, id).Error
    return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *UserRepository) UpdatePassword(userID uint, newPasswordHash string) error {
    return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password_hash", newPasswordHash).Error
}

func (r *UserRepository) CheckPassword(userID uint, password string) (bool, error) {
    var user models.User
    if err := r.db.First(&user, userID).Error; err != nil {
        return false, err
    }
    return utils.CheckPasswordHash(password, user.PasswordHash), nil
}