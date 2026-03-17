<template>
  <div class="product-card" @click="goToProductPage">
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
      <button
        @click.stop="handleAddToCart"
        class="add-to-cart"
        :class="{ 'in-cart': isInCart }"
      >
        <i :class="isInCart ? 'fas fa-check' : 'fas fa-cart-plus'"></i>
        {{ isInCart ? "В корзине" : "Добавить" }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from "vue-router";

const props = defineProps({
  product: { type: Object, required: true },
  isInCart: { type: Boolean, default: false },
});

const emit = defineEmits(["add-to-cart"]);
const router = useRouter();

const goToProductPage = () => {
  router.push(`/product/${props.product.id}`);
};

const handleAddToCart = (event) => {
  event.stopPropagation();
  if (props.isInCart) {
    router.push("/cart");
  } else {
    emit("add-to-cart", props.product);
  }
};
</script>

<style scoped>
.product-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
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

.add-to-cart.in-cart {
  background-color: #27ae60;
}

.add-to-cart.in-cart:hover {
  background-color: #2ecc71;
}

.add-to-cart:hover {
  background-color: #2980b9;
}

.add-to-cart i {
  font-size: 1rem;
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
}
</style>
