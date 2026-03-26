package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"unicode/utf8"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
)

type AuthHandler struct {
    userRepo    *repository.UserRepository
    sessionRepo *repository.SessionRepository
}

func NewAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthHandler {
    return &AuthHandler{
        userRepo:    userRepo,
        sessionRepo: sessionRepo,
    }
}

type registerRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if req.Username == "" || req.Email == "" || req.Password == "" {
        http.Error(w, "All fields are required", http.StatusBadRequest)
        return
    }

    if utf8.RuneCountInString(req.Username) > 10 {
        http.Error(w, "Username must be at most 10 characters", http.StatusBadRequest)
        return
    }

    hash, err := utils.HashPassword(req.Password)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hash,
    }

    if err := h.userRepo.Create(user); err != nil {
        http.Error(w, "User with this email or username already exists", http.StatusConflict)
        return
    }

    token, _ := utils.GenerateSessionToken()
    session := &models.Session{
        Token:     token,
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    h.sessionRepo.Create(session)

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
        MaxAge:   86400,
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
    })
}


type loginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req loginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    user, err := h.userRepo.FindByEmail(req.Email)
    if err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    token, _ := utils.GenerateSessionToken()
    session := &models.Session{
        Token:     token,
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    h.sessionRepo.Create(session)

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
        MaxAge:   86400,
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
    })
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err == nil {
        h.sessionRepo.DeleteByToken(cookie.Value)
    }
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        MaxAge:   -1,
    })
    w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    session, err := h.sessionRepo.FindByToken(cookie.Value)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    user, err := h.userRepo.FindByID(session.UserID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":         user.ID,
        "username":   user.Username,
        "email":      user.Email,
        "created_at": user.CreatedAt.Format(time.RFC3339),
    })
}