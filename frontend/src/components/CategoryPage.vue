<template>
  <div class="category-page">
    <div class="container">
      <div v-if="category" class="category-header">
        <div class="category-icon-large">
          <img :src="category.icon" :alt="category.name" />
        </div>
        <div class="category-info">
          <h1>{{ category.name }}</h1>
          <p class="category-count">{{ filteredProducts.length }} товаров</p>
        </div>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Загрузка товаров...</p>
      </div>
      <div v-else-if="error" class="error-state">
        <i class="fas fa-exclamation-circle"></i>
        <p>{{ error }}</p>
        <button @click="fetchCategoryData(route.params.id)" class="retry-btn">
          <i class="fas fa-redo"></i> Повторить
        </button>
      </div>

      <div v-else-if="category" class="content-layout">
        <button class="mobile-filter-toggle" @click="showMobileFilters = true">
          <i class="fas fa-sliders-h"></i> Фильтры
        </button>

        <aside class="filter-sidebar">
          <div class="filter-header">
            <h3>Фильтры</h3>
            <button
              v-if="hasActiveFilters"
              @click="resetFilters"
              class="reset-filters"
            >
              <i class="fas fa-times"></i> Сбросить
            </button>
          </div>

          <div class="filter-section">
            <h4>Цена, ₽</h4>
            <div class="price-inputs">
              <input
                type="number"
                v-model.number="filters.priceMin"
                placeholder="от"
                min="0"
                :max="priceMax"
              />
              <span>—</span>
              <input
                type="number"
                v-model.number="filters.priceMax"
                placeholder="до"
                :min="priceMin"
              />
            </div>
          </div>

          <div class="filter-section">
            <h4>Бренд</h4>
            <div class="brand-list">
              <label
                v-for="brand in availableBrands"
                :key="brand"
                class="brand-checkbox"
              >
                <input
                  type="checkbox"
                  :value="brand"
                  v-model="filters.brands"
                />
                {{ brand }}
              </label>
            </div>
          </div>

          <div class="filter-section">
            <h4>Рейтинг</h4>
            <select v-model="filters.minRating" class="rating-select">
              <option :value="0">Любой</option>
              <option :value="4">4 и выше</option>
              <option :value="4.5">4.5 и выше</option>
            </select>
          </div>
        </aside>

        <transition name="slide">
          <div v-if="showMobileFilters" class="mobile-filter-panel">
            <div class="mobile-filter-header">
              <h3>Фильтры</h3>
              <button @click="showMobileFilters = false">
                <i class="fas fa-times"></i>
              </button>
            </div>
            <div class="filter-section">
              <h4>Цена, ₽</h4>
              <div class="price-inputs">
                <input
                  type="number"
                  v-model.number="filters.priceMin"
                  placeholder="от"
                />
                <span>—</span>
                <input
                  type="number"
                  v-model.number="filters.priceMax"
                  placeholder="до"
                />
              </div>
            </div>
            <div class="filter-section">
              <h4>Бренд</h4>
              <div class="brand-list">
                <label
                  v-for="brand in availableBrands"
                  :key="brand"
                  class="brand-checkbox"
                >
                  <input
                    type="checkbox"
                    :value="brand"
                    v-model="filters.brands"
                  />
                  {{ brand }}
                </label>
              </div>
            </div>
            <div class="filter-section">
              <h4>Рейтинг</h4>
              <select v-model="filters.minRating">
                <option :value="0">Любой</option>
                <option :value="4">4 и выше</option>
                <option :value="4.5">4.5 и выше</option>
              </select>
            </div>
            <button
              v-if="hasActiveFilters"
              @click="resetFilters"
              class="reset-filters"
            >
              Сбросить фильтры
            </button>
            <button @click="showMobileFilters = false" class="apply-filters">
              Применить
            </button>
          </div>
        </transition>

        <section class="products-section">
          <div class="products-header">
            <div class="results-info">
              Показано {{ filteredProducts.length }} из
              {{ products.length }} товаров
            </div>
            <div class="sort-selector">
              <label for="sort">Сортировка:</label>
              <select id="sort" v-model="sortBy">
                <option value="default">По умолчанию</option>
                <option value="price-asc">Цена: по возрастанию</option>
                <option value="price-desc">Цена: по убыванию</option>
                <option value="rating">По рейтингу</option>
                <option value="name">По названию</option>
              </select>
            </div>
          </div>

          <div v-if="filteredProducts.length === 0" class="no-products">
            <i class="fas fa-box-open"></i>
            <p>Товары не найдены. Попробуйте изменить фильтры.</p>
          </div>

          <div v-else class="products-grid">
            <ProductCard
              v-for="product in filteredProducts"
              :key="product.id"
              :product="product"
              :quantity="getQuantity(product.id)"
              @add-to-cart="$emit('add-to-cart', $event)"
              @update-cart="
                (productId, newQty) => $emit('update-cart', productId, newQty)
              "
            />
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useRoute } from "vue-router";
import axios from "axios";
import ProductCard from "./ProductCard.vue";
import { categoryProductsCache } from "../utils/cache";

const route = useRoute();
const props = defineProps({
  cart: { type: Array, required: true },
});
const emit = defineEmits(["add-to-cart", "update-cart"]);

const loading = ref(false);
const error = ref(null);
const category = ref(null);
const products = ref([]);
const sortBy = ref("default");
const filters = ref({
  priceMin: null,
  priceMax: null,
  brands: [],
  minRating: 0,
});
const showMobileFilters = ref(false);

const getQuantity = (productId) => {
  const item = props.cart.find((i) => i.product.id === productId);
  return item ? item.quantity : 0;
};

const fetchCategoryData = async (categoryId) => {
  if (!categoryId) {
    error.value = "ID категории не указан";
    loading.value = false;
    return;
  }

  const cached = categoryProductsCache.get(categoryId);
  if (cached) {
    category.value = cached.category;
    products.value = cached.products;
    return;
  }

  loading.value = true;
  error.value = null;
  try {
    const { data } = await axios.get(`/category/${categoryId}/products`);
    if (!data.category) throw new Error("Категория не найдена");
    const categoryData = { ...data.category, icon: data.category.icon_url };
    category.value = categoryData;
    products.value = data.products;
    categoryProductsCache.set(categoryId, {
      category: categoryData,
      products: data.products,
    });
  } catch (err) {
    error.value = err.response?.data || err.message;
    category.value = null;
  } finally {
    loading.value = false;
  }
};

watch(
  () => route.params.id,
  (newId) => {
    fetchCategoryData(newId);
  },
  { immediate: true },
);

const priceMin = computed(() =>
  products.value.length ? Math.min(...products.value.map((p) => p.price)) : 0,
);
const priceMax = computed(() =>
  products.value.length
    ? Math.max(...products.value.map((p) => p.price))
    : 100000,
);
const availableBrands = computed(() => {
  const brands = products.value.map((p) => p.name.split(" ")[0]);
  return [...new Set(brands)].sort();
});
const hasActiveFilters = computed(() => {
  return (
    filters.value.priceMin !== null ||
    filters.value.priceMax !== null ||
    filters.value.brands.length > 0 ||
    filters.value.minRating > 0
  );
});

const filteredProducts = computed(() => {
  let filtered = [...products.value];
  if (filters.value.priceMin !== null && !isNaN(filters.value.priceMin)) {
    filtered = filtered.filter((p) => p.price >= filters.value.priceMin);
  }
  if (filters.value.priceMax !== null && !isNaN(filters.value.priceMax)) {
    filtered = filtered.filter((p) => p.price <= filters.value.priceMax);
  }
  if (filters.value.brands.length > 0) {
    filtered = filtered.filter((p) =>
      filters.value.brands.includes(p.name.split(" ")[0]),
    );
  }
  if (filters.value.minRating > 0) {
    filtered = filtered.filter((p) => p.rating >= filters.value.minRating);
  }

  switch (sortBy.value) {
    case "price-asc":
      filtered.sort((a, b) => a.price - b.price);
      break;
    case "price-desc":
      filtered.sort((a, b) => b.price - a.price);
      break;
    case "rating":
      filtered.sort((a, b) => b.rating - a.rating);
      break;
    case "name":
      filtered.sort((a, b) => a.name.localeCompare(b.name));
      break;
  }
  return filtered;
});

const resetFilters = () => {
  filters.value = {
    priceMin: null,
    priceMax: null,
    brands: [],
    minRating: 0,
  };
};
</script>

<style scoped>
.category-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px 20px 40px;
}
.category-header {
  display: flex;
  align-items: center;
  gap: 30px;
  background: white;
  border-radius: 20px;
  padding: 40px;
  margin-bottom: 30px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  border: 1px solid #e0e0e0;
}
.category-icon-large {
  width: 120px;
  height: 120px;
  background-color: #4a90e2;
  border-radius: 30px;
  padding: 25px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.category-icon-large img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  filter: brightness(0) invert(1);
}
.category-info h1 {
  font-size: 2.5rem;
  color: #333;
  margin-bottom: 10px;
}
.category-count {
  color: #666;
  font-size: 1.1rem;
}
.loading-state,
.error-state {
  text-align: center;
  padding: 80px 20px;
  background: white;
  border-radius: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
}
.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #4a90e2;
  border-radius: 50%;
  margin: 0 auto 20px;
  animation: spin 1s linear infinite;
}
@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
.error-state i {
  font-size: 3rem;
  color: #e74c3c;
  margin-bottom: 20px;
}
.retry-btn {
  padding: 12px 30px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s;
}
.retry-btn:hover {
  background-color: #357abd;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(74, 144, 226, 0.3);
}
.content-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 30px;
  margin-top: 20px;
}
.mobile-filter-toggle {
  display: none;
  width: 100%;
  padding: 12px;
  margin-bottom: 20px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  align-items: center;
  justify-content: center;
  gap: 10px;
}
.mobile-filter-toggle:hover {
  background-color: #357abd;
}
.filter-sidebar {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  border: 1px solid #e0e0e0;
  height: fit-content;
  position: sticky;
  top: 100px;
}
.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid #4a90e2;
}
.filter-header h3 {
  font-size: 1.2rem;
  color: #333;
  margin: 0;
}
.reset-filters {
  background: none;
  border: none;
  color: #e74c3c;
  cursor: pointer;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: 5px;
}
.reset-filters:hover {
  text-decoration: underline;
}
.filter-section {
  margin-bottom: 25px;
}
.filter-section h4 {
  color: #555;
  font-size: 1rem;
  margin-bottom: 12px;
  font-weight: 600;
}
.price-inputs {
  display: flex;
  align-items: center;
  gap: 8px;
}
.price-inputs input {
  width: 100%;
  padding: 8px 12px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 0.9rem;
}
.price-inputs input:focus {
  outline: none;
  border-color: #4a90e2;
}
.brand-list {
  max-height: 200px;
  overflow-y: auto;
  padding-right: 5px;
}
.brand-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
  cursor: pointer;
  font-size: 0.95rem;
  color: #666;
}
.brand-checkbox input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}
.rating-select {
  width: 100%;
  padding: 8px 12px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 0.95rem;
  color: #333;
}
.mobile-filter-panel {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: white;
  z-index: 2000;
  padding: 20px;
  overflow-y: auto;
  transform: translateX(-100%);
  transition: transform 0.3s ease;
}
.slide-enter-active .mobile-filter-panel,
.slide-leave-active .mobile-filter-panel {
  transform: translateX(0);
}
.mobile-filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 10px;
  border-bottom: 2px solid #4a90e2;
}
.mobile-filter-header h3 {
  font-size: 1.3rem;
  color: #333;
}
.mobile-filter-header button {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
}
.apply-filters {
  width: 100%;
  padding: 14px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  margin-top: 20px;
  cursor: pointer;
}
.apply-filters:hover {
  background-color: #357abd;
}
.products-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  background: white;
  padding: 15px 20px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}
.results-info {
  color: #666;
  font-size: 0.95rem;
}
.sort-selector {
  display: flex;
  align-items: center;
  gap: 10px;
}
.sort-selector label {
  color: #666;
  font-size: 0.95rem;
}
.sort-selector select {
  padding: 8px 15px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 0.95rem;
  color: #333;
  cursor: pointer;
  transition: all 0.3s;
}
.sort-selector select:hover,
.sort-selector select:focus {
  border-color: #4a90e2;
  outline: none;
}
.no-products {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
}
.no-products i {
  font-size: 4rem;
  color: #999;
  margin-bottom: 20px;
}
.no-products p {
  color: #666;
  font-size: 1.1rem;
}
.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 25px;
}
@media (max-width: 1024px) {
  .content-layout {
    grid-template-columns: 240px 1fr;
    gap: 20px;
  }
}
@media (max-width: 768px) {
  .content-layout {
    grid-template-columns: 1fr;
  }
  .filter-sidebar {
    display: none;
  }
  .mobile-filter-toggle {
    display: flex;
  }
  .category-header {
    flex-direction: column;
    text-align: center;
    padding: 30px;
  }
  .category-info h1 {
    font-size: 2rem;
  }
  .products-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  .sort-selector {
    width: 100%;
  }
  .sort-selector select {
    flex: 1;
  }
  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: 15px;
  }
}
@media (max-width: 480px) {
  .category-icon-large {
    width: 100px;
    height: 100px;
    padding: 20px;
  }
  .category-info h1 {
    font-size: 1.5rem;
  }
  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 12px;
  }
}
</style>
