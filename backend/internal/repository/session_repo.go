package repository

import (
	"backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type SessionRepository struct {
    db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
    return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *models.Session) error {
    return r.db.Create(session).Error
}

func (r *SessionRepository) FindByToken(token string) (*models.Session, error) {
    var session models.Session
    err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error
    return &session, err
}

func (r *SessionRepository) DeleteByToken(token string) error {
    return r.db.Where("token = ?", token).Delete(&models.Session{}).Error
}