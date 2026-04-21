<template>
  <div class="search-page">
    <div class="container">
      <h1>Поиск: «{{ searchQuery }}»</h1>
      <p v-if="!loading && filteredProducts.length === 0" class="no-results">
        Ничего не найдено. Попробуйте изменить запрос.
      </p>

      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Загрузка товаров...</p>
      </div>

      <div v-else-if="filteredProducts.length > 0" class="products-section">
        <div class="results-info">Найдено {{ filteredProducts.length }} товаров</div>
        <div class="products-grid">
          <ProductCard
            v-for="product in filteredProducts"
            :key="product.id"
            :product="product"
            :quantity="getQuantity(product.id)"
            @add-to-cart="$emit('add-to-cart', $event)"
            @update-cart="(productId, newQty) => $emit('update-cart', productId, newQty)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import axios from "axios";
import ProductCard from "./ProductCard.vue";
import { allProductsCache } from "../utils/cache";
import { useBreadcrumbs } from "../composables/useBreadcrumbs";
import api from "../api";

const props = defineProps({
  searchQuery: { type: String, default: "" },
  cart: { type: Array, required: true },
});

const route = useRoute();
const { setBreadcrumbs } = useBreadcrumbs();

const loading = ref(false);
const allProducts = ref([]);
let isFetching = false;

const fetchAllProducts = async () => {
  const cached = allProductsCache.get();
  if (cached) {
    allProducts.value = cached;
    return;
  }

  if (isFetching) return;

  isFetching = true;
  loading.value = true;
  try {
    const { data: cats } = await api.get("/categories");
    const promises = cats.map(async (cat) => {
      try {
        const response = await api.get(`/${cat.url_key}`);
        return response.data;
      } catch {
        return [];
      }
    });
    const results = await Promise.all(promises);
    const products = results.flat();
    allProducts.value = products;
    allProductsCache.set(products);
  } catch (err) {
    console.error("Ошибка загрузки товаров:", err);
  } finally {
    loading.value = false;
    isFetching = false;
  }
};

const filteredProducts = computed(() => {
  if (!props.searchQuery) return [];
  const query = props.searchQuery.toLowerCase().trim();
  return allProducts.value.filter((product) =>
    product.name.toLowerCase().includes(query)
  );
});

const getQuantity = (productId) => {
  const cartItem = props.cart.find((item) => item.product.id === productId);
  return cartItem ? cartItem.quantity : 0;
};

const updateBreadcrumbs = () => {
  setBreadcrumbs([
    { name: "Главная", path: "/" },
    {
      name: `Поиск: ${props.searchQuery}`,
      path: `/search?q=${encodeURIComponent(props.searchQuery)}`,
    },
  ]);
};

onMounted(() => {
  fetchAllProducts();
  updateBreadcrumbs();
});

watch(() => props.searchQuery, updateBreadcrumbs);
</script>

<style scoped>
.search-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 40px 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

h1 {
  color: #333;
  margin-bottom: 20px;
  font-size: 2rem;
}

.no-results {
  text-align: center;
  color: #666;
  font-size: 1.2rem;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.loading-state {
  text-align: center;
  padding: 80px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
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

.products-section {
  background: white;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.results-info {
  color: #666;
  margin-bottom: 20px;
  font-size: 1rem;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 25px;
}

@media (max-width: 768px) {
  h1 {
    font-size: 1.5rem;
  }
  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: 15px;
  }
}

@media (max-width: 480px) {
  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 12px;
  }
}
</style>
