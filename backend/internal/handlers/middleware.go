package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            cookie, err := r.Cookie("session_token")
            if err != nil {
                next.ServeHTTP(w, r)
                return
            }

            session, err := sessionRepo.FindByToken(cookie.Value)
            if err != nil {
                next.ServeHTTP(w, r)
                return
            }

            // Обновляем время истечения
            session.ExpiresAt = time.Now().Add(15 * time.Minute)
            sessionRepo.Update(session)

            user, err := userRepo.FindByID(session.UserID)
            if err != nil {
                next.ServeHTTP(w, r)
                return
            }

            ctx := context.WithValue(r.Context(), UserContextKey, user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserFromContext(ctx context.Context) *models.User {
    user, ok := ctx.Value(UserContextKey).(*models.User)
    if !ok {
        return nil
    }
    return user
}