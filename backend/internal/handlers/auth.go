package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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

func validatePassword(password string) error {
    if len(password) < 8 {
        return fmt.Errorf("пароль должен содержать минимум 8 символов")
    }
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    if !hasUpper {
        return fmt.Errorf("пароль должен содержать хотя бы одну заглавную букву")
    }
    hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
    if !hasDigit {
        return fmt.Errorf("пароль должен содержать хотя бы одну цифру")
    }
    hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
    if !hasSpecial {
        return fmt.Errorf("пароль должен содержать хотя бы один специальный символ (!@#$%^&* и т.д.)")
    }
    return nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Неверный запрос", http.StatusBadRequest)
        return
    }

    if req.Username == "" || req.Email == "" || req.Password == "" {
        http.Error(w, "Все поля обязательны для заполнения", http.StatusBadRequest)
        return
    }

    if utf8.RuneCountInString(req.Username) > 10 {
        http.Error(w, "Имя пользователя не должно превышать 10 символов", http.StatusBadRequest)
        return
    }

    if err := validatePassword(req.Password); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    hash, err := utils.HashPassword(req.Password)
    if err != nil {
        http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
        return
    }

    user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hash,
    }

    if err := h.userRepo.Create(user); err != nil {
        if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
            http.Error(w, "Пользователь с таким email уже существует", http.StatusConflict)
        } else {
            http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
        }
        return
    }

    token, _ := utils.GenerateSessionToken()
    session := &models.Session{
        Token:     token,
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    if err := h.sessionRepo.Create(session); err != nil {
        http.Error(w, "Failed to create session", http.StatusInternalServerError)
        return
    }
    h.sessionRepo.Create(session)
    

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteNoneMode,
        Secure:   true, 
        MaxAge:   86400,
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
        "token":    token,
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
    if err := h.sessionRepo.Create(session); err != nil {
        http.Error(w, "Failed to create session", http.StatusInternalServerError)
        return
    }
    h.sessionRepo.Create(session)

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteNoneMode,
        Secure:   true, 
        MaxAge:   86400,
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
        "token":    token,
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
    var token string

    authHeader := r.Header.Get("Authorization")
    if strings.HasPrefix(authHeader, "Bearer ") {
        token = strings.TrimPrefix(authHeader, "Bearer ")
    }

    if token == "" {
        cookie, err := r.Cookie("session_token")
        if err == nil {
            token = cookie.Value
        }
    }

    if token == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    session, err := h.sessionRepo.FindByToken(token)
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