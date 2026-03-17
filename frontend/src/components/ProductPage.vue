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

          <button
            @click="handleAddToCart"
            class="add-to-cart"
            :class="{ 'in-cart': isInCart }"
          >
            <i :class="isInCart ? 'fas fa-check' : 'fas fa-cart-plus'"></i>
            {{ isInCart ? "В корзине" : "Добавить в корзину" }}
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

const props = defineProps({
  cart: { type: Array, required: true },
});

const emit = defineEmits(["add-to-cart"]);

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const error = ref(null);
const product = ref(null);

const fetchProduct = async () => {
  loading.value = true;
  error.value = null;

  const productId = parseInt(route.params.id);
  if (isNaN(productId)) {
    error.value = "Неверный идентификатор товара";
    loading.value = false;
    return;
  }

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
    const found = allProducts.find((p) => p.id === productId);
    if (found) {
      product.value = found;
    } else {
      error.value = "Товар не найден";
    }
  } catch (err) {
    console.error("Ошибка загрузки товара:", err);
    error.value = "Не удалось загрузить товар. Попробуйте позже.";
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchProduct();
});

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

const isInCart = computed(() => {
  return props.cart.some((item) => item.product.id === product.value?.id);
});

const handleAddToCart = () => {
  if (isInCart.value) {
    router.push("/cart");
  } else {
    emit("add-to-cart", product.value);
  }
};

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
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 20px;
}

.add-to-cart.in-cart {
  background-color: #27ae60;
}

.add-to-cart.in-cart:hover {
  background-color: #2ecc71;
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
}
</style>
