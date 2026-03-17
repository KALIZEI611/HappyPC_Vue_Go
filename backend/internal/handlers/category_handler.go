package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
    categoryRepo *repository.CategoryRepository
    productRepo  *repository.ProductRepository
}

func NewCategoryHandler(categoryRepo *repository.CategoryRepository, productRepo *repository.ProductRepository) *CategoryHandler {
    return &CategoryHandler{
        categoryRepo: categoryRepo,
        productRepo:  productRepo,
    }
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    categories, err := h.categoryRepo.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) GetProductsByURLKey(w http.ResponseWriter, r *http.Request) {
    urlKey := chi.URLParam(r, "urlKey")
    cat, err := h.categoryRepo.GetByURLKey(urlKey)
    if err != nil {
        http.Error(w, "Category not found", http.StatusNotFound)
        return
    }
    products, err := h.productRepo.GetByCategoryID(cat.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }

    cat, err := h.categoryRepo.GetByID(uint(id))
    if err != nil {
        http.Error(w, "Category not found", http.StatusNotFound)
        return
    }

    products, err := h.productRepo.GetByCategoryID(cat.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := struct {
        Category models.Category `json:"category"`
        Products []models.Product `json:"products"`
    }{
        Category: *cat,
        Products: products,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}