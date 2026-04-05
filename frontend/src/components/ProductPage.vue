<template>
  <div class="product-page">
    <div class="container">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Загрузка товара...</p>
      </div>
      <div v-else-if="error" class="error-state">
        <i class="fas fa-exclamation-circle"></i>
        <p>{{ error }}</p>
        <button @click="fetchProduct" class="retry-btn">
          <i class="fas fa-redo"></i> Повторить
        </button>
      </div>
      <div v-else-if="product" class="product-detail">
        <div class="product-gallery">
          <img :src="product.image" :alt="product.name" />
        </div>
        <div class="product-info">
          <h1>{{ product.name }}</h1>
          <div class="product-rating">
            <i
              v-for="star in 5"
              :key="star"
              class="fas fa-star"
              :class="{ filled: star <= Math.floor(product.rating) }"
            ></i>
            <span class="rating-value">{{ product.rating }}</span>
          </div>
          <div class="product-price">{{ product.price.toLocaleString() }} ₽</div>

          <div class="product-description" v-if="product.description">
            <h3>Описание</h3>
            <p>{{ product.description }}</p>
          </div>

          <div class="product-specs" v-if="Object.keys(parsedSpecs).length">
            <h3>Характеристики</h3>
            <ul>
              <li v-for="(value, key) in parsedSpecs" :key="key">
                <strong>{{ formatSpecKey(key) }}:</strong> {{ value }}
              </li>
            </ul>
          </div>

          <div v-if="quantity > 0" class="quantity-controls" @click.stop>
            <button @click="decrement" class="qty-btn">
              <i class="fas fa-minus"></i>
            </button>
            <span class="quantity">{{ quantity }}</span>
            <button @click="increment" class="qty-btn">
              <i class="fas fa-plus"></i>
            </button>
          </div>

          <button v-else @click="handleAddToCart" class="add-to-cart">
            <i class="fas fa-cart-plus"></i> Добавить в корзину
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import axios from "axios";
import { specMapping } from "../constants/specMapping";
import { allProductsCache, homeCategoriesCache } from "../utils/cache";
import { useBreadcrumbs } from "../composables/useBreadcrumbs";

const props = defineProps({
  cart: { type: Array, required: true },
});
const emit = defineEmits(["add-to-cart", "update-cart"]);

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const error = ref(null);
const product = ref(null);
let isFetching = false;

const { setBreadcrumbs } = useBreadcrumbs();

const quantity = computed(() => {
  const cartItem = props.cart.find((item) => item.product.id === product.value?.id);
  return cartItem ? cartItem.quantity : 0;
});

const fetchProduct = async () => {
  loading.value = true;
  error.value = null;

  const productId = parseInt(route.params.id);
  if (isNaN(productId)) {
    error.value = "Неверный идентификатор товара";
    loading.value = false;
    return;
  }

  const cachedAllProducts = allProductsCache.get();
  if (cachedAllProducts) {
    const found = cachedAllProducts.find((p) => p.id === productId);
    if (found) {
      product.value = found;
      await setProductBreadcrumbs(found);
      loading.value = false;
      return;
    }
  }

  if (isFetching) return;
  isFetching = true;

  try {
    const { data: cats } = await axios.get("/categories");
    const promises = cats.map(async (cat) => {
      try {
        const response = await axios.get(`/${cat.url_key}`);
        return response.data;
      } catch {
        return [];
      }
    });
    const results = await Promise.all(promises);
    const allProducts = results.flat();
    allProductsCache.set(allProducts);
    const found = allProducts.find((p) => p.id === productId);
    if (found) {
      product.value = found;
      await setProductBreadcrumbs(found);
    } else {
      error.value = "Товар не найден";
    }
  } catch (err) {
    console.error("Ошибка загрузки товара:", err);
    error.value = "Не удалось загрузить товар. Попробуйте позже.";
  } finally {
    loading.value = false;
    isFetching = false;
  }
};

const setProductBreadcrumbs = async (productData) => {
  if (route.query.from === "search" && route.query.q) {
    const searchQuery = route.query.q;
    setBreadcrumbs([
      { name: "Главная", path: "/" },
      {
        name: `Поиск: ${searchQuery}`,
        path: `/search?q=${encodeURIComponent(searchQuery)}`,
      },
      { name: productData.name, path: `/product/${productData.id}` },
    ]);
    return;
  }

  let categoryName = null;
  const cachedCategories = homeCategoriesCache.get();
  if (cachedCategories) {
    const category = cachedCategories.find((c) => c.id === productData.category_id);
    if (category) categoryName = category.name;
  }
  if (!categoryName && productData.category_id) {
    try {
      const { data: category } = await axios.get(
        `/category/${productData.category_id}/products`
      );
      categoryName = category.category.name;
    } catch {
      categoryName = "Категория";
    }
  }
  setBreadcrumbs([
    { name: "Главная", path: "/" },
    {
      name: categoryName || "Категория",
      path: `/category/${productData.category_id}`,
    },
    { name: productData.name, path: `/product/${productData.id}` },
  ]);
};

onMounted(() => {
  fetchProduct();
});

const increment = () => {
  emit("update-cart", product.value.id, quantity.value + 1);
};

const decrement = () => {
  if (quantity.value === 1) {
    emit("update-cart", product.value.id, 0);
  } else {
    emit("update-cart", product.value.id, quantity.value - 1);
  }
};

const handleAddToCart = () => {
  emit("add-to-cart", product.value);
};

const parsedSpecs = computed(() => {
  if (!product.value?.specs) return {};
  if (typeof product.value.specs === "string") {
    try {
      return JSON.parse(product.value.specs);
    } catch {
      return {};
    }
  }
  return product.value.specs;
});

const formatSpecKey = (key) => specMapping[key] || key;
</script>

<style scoped>
.product-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 40px 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
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

.product-detail {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  background: white;
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  border: 1px solid #e0e0e0;
}

.product-gallery {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: white;
  border-radius: 16px;
  padding: 20px;
  min-height: 600px;
}

.product-gallery img {
  width: 100%;
  height: auto;
  max-height: 800px;
  object-fit: contain;
}

.product-info h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 15px;
}

.product-rating {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-bottom: 15px;
}

.fa-star {
  color: #ddd;
  font-size: 1.2rem;
}

.fa-star.filled {
  color: #f1c40f;
}

.rating-value {
  margin-left: 8px;
  color: #666;
  font-size: 1rem;
  font-weight: 500;
}

.product-price {
  font-size: 2rem;
  font-weight: bold;
  color: #2c3e50;
  margin-bottom: 25px;
}

.product-description,
.product-specs {
  margin-bottom: 25px;
}

.product-description h3,
.product-specs h3 {
  color: #333;
  font-size: 1.2rem;
  margin-bottom: 10px;
  border-bottom: 2px solid #4a90e2;
  padding-bottom: 5px;
  display: inline-block;
}

.product-description p {
  color: #666;
  line-height: 1.6;
}

.product-specs ul {
  list-style: none;
  padding: 0;
}

.product-specs li {
  padding: 8px 0;
  border-bottom: 1px solid #eee;
  color: #555;
}

.product-specs li strong {
  color: #333;
  min-width: 150px;
  display: inline-block;
}

.quantity-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 12px;
  background: #f8f9fa;
  border-radius: 30px;
  padding: 8px 12px;
  border: 1px solid #e0e0e0;
  width: fit-content;
  margin: 20px auto 0;
  box-sizing: border-box;
}

.qty-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: white;
  color: #4a90e2;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
  font-size: 1rem;
}

.qty-btn:hover:not(:disabled) {
  background: #4a90e2;
  color: white;
  transform: scale(1.1);
  box-shadow: 0 4px 10px rgba(74, 144, 226, 0.3);
}

.qty-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: #f1f1f1;
  color: #999;
}

.quantity {
  font-weight: 600;
  color: #333;
  min-width: 30px;
  text-align: center;
  font-size: 1.2rem;
}

.add-to-cart {
  width: 100%;
  padding: 15px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 20px;
}

.add-to-cart:hover {
  background-color: #2980b9;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(52, 152, 219, 0.3);
}

@media (max-width: 768px) {
  .product-detail {
    grid-template-columns: 1fr;
    padding: 20px;
  }

  .product-gallery {
    min-height: 400px;
  }

  .product-gallery img {
    max-height: 500px;
  }

  .product-info h1 {
    font-size: 1.5rem;
  }

  .product-price {
    font-size: 1.5rem;
  }

  .quantity-controls {
    padding: 6px 10px;
    gap: 8px;
  }

  .qty-btn {
    width: 30px;
    height: 30px;
    font-size: 0.8rem;
  }

  .quantity {
    min-width: 24px;
    font-size: 1rem;
  }
}
</style>
