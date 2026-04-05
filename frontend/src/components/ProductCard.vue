<template>
  <div v-if="product" class="product-card" @click="goToProductPage">
    <div class="product-image-container">
      <div class="product-image">
        <img :src="product.image" :alt="product.name" loading="lazy" />
      </div>
    </div>
    <div class="product-info">
      <h3 class="product-title">{{ product.name }}</h3>
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

      <!-- Режим сборки ПК: показываем кнопку "В сборку" -->
      <div v-if="isBuilderMode">
        <button @click.stop="addToBuild" class="add-to-build-btn">
          <i class="fas fa-tools"></i> В сборку
        </button>
      </div>

      <!-- Обычный режим: показываем управление корзиной -->
      <div v-else>
        <div v-if="quantity > 0" class="quantity-controls" @click.stop>
          <button @click="decrement" class="qty-btn">
            <i class="fas fa-minus"></i>
          </button>
          <span class="quantity">{{ quantity }}</span>
          <button @click="increment" class="qty-btn">
            <i class="fas fa-plus"></i>
          </button>
        </div>
        <button v-else @click.stop="handleAddToCart" class="add-to-cart">
          <i class="fas fa-cart-plus"></i> Добавить
        </button>
      </div>
    </div>
  </div>
  <div v-else class="product-card error">Товар временно недоступен</div>
</template>

<script setup>
import { useRouter, useRoute } from "vue-router";
import { computed } from "vue";

const props = defineProps({
  product: { type: Object, required: true },
  quantity: { type: Number, default: 0 },
});

const emit = defineEmits(["add-to-cart", "update-cart", "add-to-build"]);
const router = useRouter();
const route = useRoute();

const isBuilderMode = computed(() => route.query.from === "builder");

const goToProductPage = () => {
  if (route.name === "Search") {
    router.push({
      path: `/product/${props.product.id}`,
      query: {
        from: "search",
        q: route.query.q || "",
      },
    });
  } else {
    router.push(`/product/${props.product.id}`);
  }
};

const handleAddToCart = () => {
  emit("add-to-cart", props.product);
};

const addToBuild = () => {
  emit("add-to-build", props.product);
};

const increment = () => {
  emit("update-cart", props.product.id, props.quantity + 1);
};

const decrement = () => {
  if (props.quantity === 1) {
    emit("update-cart", props.product.id, 0);
  } else {
    emit("update-cart", props.product.id, props.quantity - 1);
  }
};
</script>

<style scoped>
.product-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  transition:
    transform 0.3s,
    box-shadow 0.3s;
  border: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  cursor: pointer;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.product-image-container {
  position: relative;
  width: 100%;
  padding-top: 100%;
  overflow: hidden;
  background-color: white;
}

.product-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  object-position: center;
  transition: transform 0.5s;
  padding: 10px;
  background-color: white;
}

.product-info {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.product-title {
  margin: 0;
  font-size: 1rem;
  color: #333;
  font-weight: 600;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  height: 2.8em;
  transition: color 0.3s;
}

.product-title:hover {
  color: #4a90e2;
}

.product-rating {
  display: flex;
  align-items: center;
  gap: 4px;
}

.fa-star {
  color: #ddd;
  font-size: 0.85rem;
}

.fa-star.filled {
  color: #f1c40f;
}

.rating-value {
  margin-left: 4px;
  color: #666;
  font-size: 0.85rem;
  font-weight: 500;
}

.product-price {
  font-size: 1.2rem;
  font-weight: bold;
  color: #2c3e50;
  margin: 4px 0;
}

.quantity-controls {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  background: #f8f9fa;
  border-radius: 30px;
  padding: 8px 12px;
  border: 1px solid #e0e0e0;
  width: fit-content;
  margin: 0 auto;
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
  padding: 10px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-weight: 500;
  margin-top: 4px;
}

.add-to-cart:hover {
  background-color: #2980b9;
}

.add-to-build-btn {
  width: 100%;
  padding: 10px;
  background-color: #6c5ce7;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-weight: 500;
  margin-top: 4px;
}

.add-to-build-btn:hover {
  background-color: #5a4bcf;
}

@media (max-width: 768px) {
  .product-image-container {
    padding-top: 100%;
  }

  .product-info {
    padding: 12px;
  }

  .product-title {
    font-size: 0.95rem;
  }

  .product-price {
    font-size: 1.1rem;
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
