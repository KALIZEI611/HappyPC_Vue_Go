package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/models"

	"gorm.io/gorm"
)

type FeedbackHandler struct {
	db *gorm.DB
}

func NewFeedbackHandler(db *gorm.DB) *FeedbackHandler {
	return &FeedbackHandler{db: db}
}

type feedbackRequest struct {
	Subject     string `json:"subject"`
	Message     string `json:"message"`
	CopyToEmail bool   `json:"copyToEmail"`
}

func (h *FeedbackHandler) SendFeedback(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req feedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Subject == "" || req.Message == "" {
		http.Error(w, "Subject and message are required", http.StatusBadRequest)
		return
	}

	feedback := &models.Feedback{
		UserID:      user.ID,
		Subject:     req.Subject,
		Message:     req.Message,
		Status:      "pending",
		CopyToEmail: req.CopyToEmail,
		CreatedAt:   time.Now(),
	}

	if err := h.db.Create(feedback).Error; err != nil {
		http.Error(w, "Failed to save feedback", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Feedback sent successfully",
	})
}