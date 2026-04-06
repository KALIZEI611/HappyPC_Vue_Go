package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type BuildHandler struct {
    buildRepo *repository.BuildRepository
    db        *gorm.DB
}

func NewBuildHandler(buildRepo *repository.BuildRepository, db *gorm.DB) *BuildHandler {
    return &BuildHandler{
        buildRepo: buildRepo,
        db:        db,
    }
}

type saveBuildRequest struct {
    Name       string                 `json:"name"`
    Components map[string]interface{} `json:"components"`
}

func (h *BuildHandler) SaveBuild(w http.ResponseWriter, r *http.Request) {
    user := GetUserFromContext(r.Context())
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req saveBuildRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if req.Name == "" {
        http.Error(w, "Build name is required", http.StatusBadRequest)
        return
    }

    componentsJSON, err := json.Marshal(req.Components)
    if err != nil {
        http.Error(w, "Failed to serialize components", http.StatusInternalServerError)
        return
    }

    build := &models.Build{
        UserID:     user.ID,
        Name:       req.Name,
        Components: componentsJSON,
    }

    if err := h.buildRepo.Create(build); err != nil {
        http.Error(w, "Failed to save build", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":      build.ID,
        "message": "Build saved successfully",
    })
}

func (h *BuildHandler) GetUserBuilds(w http.ResponseWriter, r *http.Request) {
    user := GetUserFromContext(r.Context())
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    builds, err := h.buildRepo.GetUserBuilds(user.ID)
    if err != nil {
        http.Error(w, "Failed to fetch builds", http.StatusInternalServerError)
        return
    }

    type buildResponse struct {
        ID         uint                   `json:"id"`
        CreatedAt  time.Time              `json:"created_at"`
        Name       string                 `json:"name"`
        Components map[string]interface{} `json:"components"`
    }
    
    response := make([]buildResponse, 0, len(builds))
    for _, build := range builds {
        var components map[string]interface{}
        if err := json.Unmarshal(build.Components, &components); err != nil {
            components = make(map[string]interface{})
        }
        response = append(response, buildResponse{
            ID:         build.ID,
            CreatedAt:  build.CreatedAt,
            Name:       build.Name,
            Components: components,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *BuildHandler) DeleteBuild(w http.ResponseWriter, r *http.Request) {
    user := GetUserFromContext(r.Context())
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    idParam := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "Invalid build ID", http.StatusBadRequest)
        return
    }

    var build models.Build
    if err := h.db.First(&build, id).Error; err != nil {
        http.Error(w, "Build not found", http.StatusNotFound)
        return
    }

    if build.UserID != user.ID {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    if err := h.buildRepo.Delete(uint(id)); err != nil {
        http.Error(w, "Failed to delete build", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}