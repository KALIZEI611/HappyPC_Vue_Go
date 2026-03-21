<template>
  <div class="cart-page">
    <div class="container">
      <h1>Корзина</h1>

      <div v-if="cart.length === 0" class="empty-cart">
        <i class="fas fa-shopping-cart"></i>
        <p>Корзина пуста</p>
        <router-link to="/" class="continue-shopping"
          >Продолжить покупки</router-link
        >
      </div>

      <div v-else class="cart-content">
        <div class="cart-items">
          <div v-for="item in cart" :key="item.product.id" class="cart-item">
            <router-link
              :to="'/product/' + item.product.id"
              class="item-image-link"
            >
              <div class="item-image">
                <img :src="item.product.image" :alt="item.product.name" />
              </div>
            </router-link>

            <div class="item-info">
              <router-link
                :to="'/product/' + item.product.id"
                class="item-title-link"
              >
                <h3>{{ item.product.name }}</h3>
              </router-link>
              <div class="item-price">
                {{ item.product.price.toLocaleString() }} ₽
              </div>
            </div>

            <div class="quantity-controls">
              <button @click="decrement(item.product.id)" class="qty-btn">
                <i class="fas fa-minus"></i>
              </button>
              <span class="quantity">{{ item.quantity }}</span>
              <button @click="increment(item.product.id)" class="qty-btn">
                <i class="fas fa-plus"></i>
              </button>
            </div>

            <div class="item-total">
              {{ (item.product.price * item.quantity).toLocaleString() }} ₽
            </div>

            <button @click="removeItem(item.product.id)" class="remove-btn">
              <i class="fas fa-trash"></i>
            </button>
          </div>
        </div>

        <div class="cart-summary">
          <h3>Итого</h3>
          <div class="summary-row">
            <span>Товаров:</span>
            <span>{{ totalItems }}</span>
          </div>
          <div class="summary-row total">
            <span>Общая сумма:</span>
            <span>{{ totalPrice.toLocaleString() }} ₽</span>
          </div>
          <button class="checkout-btn">Оформить заказ</button>
          <router-link to="/" class="continue-shopping"
            >Продолжить покупки</router-link
          >
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  cart: { type: Array, required: true },
});

const emit = defineEmits(["update-cart", "remove-from-cart"]);

const totalItems = computed(() =>
  props.cart.reduce((sum, item) => sum + item.quantity, 0),
);

const totalPrice = computed(() =>
  props.cart.reduce((sum, item) => sum + item.product.price * item.quantity, 0),
);

const increment = (productId) => {
  const item = props.cart.find((i) => i.product.id === productId);
  if (item) emit("update-cart", productId, item.quantity + 1);
};

const decrement = (productId) => {
  const item = props.cart.find((i) => i.product.id === productId);
  if (item) {
    const newQty = item.quantity - 1;
    emit("update-cart", productId, newQty);
  }
};

const removeItem = (productId) => {
  emit("remove-from-cart", productId);
};
</script>

<style scoped>
.item-image-link,
.item-title-link {
  text-decoration: none;
  color: inherit;
  cursor: pointer;
}

.item-image-link:hover .item-image {
  opacity: 0.8;
}

.item-title-link:hover h3 {
  color: #4a90e2;
}

.cart-page {
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
  margin-bottom: 30px;
  font-size: 2rem;
}

.empty-cart {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.empty-cart i {
  font-size: 5rem;
  color: #999;
  margin-bottom: 20px;
}

.empty-cart p {
  color: #666;
  font-size: 1.2rem;
  margin-bottom: 20px;
}

.continue-shopping {
  display: inline-block;
  padding: 12px 30px;
  background-color: #4a90e2;
  color: white;
  text-decoration: none;
  border-radius: 8px;
  transition: background-color 0.3s;
  border: none;
  cursor: pointer;
  font-size: 1rem;
}

.continue-shopping:hover {
  background-color: #357abd;
}

.cart-content {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 30px;
}

.cart-items {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.cart-item {
  display: flex;
  align-items: center;
  padding: 20px 0;
  border-bottom: 1px solid #eee;
  gap: 15px;
  flex-wrap: wrap;
}

.cart-item:last-child {
  border-bottom: none;
}

.item-image-link {
  flex-shrink: 0;
}

.item-image {
  width: 100px;
  height: 100px;
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
  transition: opacity 0.2s;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 10px;
}

.item-info {
  flex: 2;
  min-width: 200px;
}

.item-title-link {
  text-decoration: none;
  color: #333;
}

.item-title-link h3 {
  font-size: 1.1rem;
  margin-bottom: 5px;
  transition: color 0.2s;
}

.item-price {
  font-size: 1.1rem;
  font-weight: bold;
  color: #2c3e50;
}

.quantity-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f8f9fa;
  border-radius: 30px;
  padding: 4px;
  border: 1px solid #e0e0e0;
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
  font-size: 1.1rem;
}

.item-total {
  font-size: 1.2rem;
  font-weight: bold;
  color: #2c3e50;
  min-width: 100px;
  text-align: right;
}

.remove-btn {
  background: none;
  border: none;
  color: #e74c3c;
  cursor: pointer;
  padding: 10px;
  border-radius: 50%;
  transition: all 0.3s;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
}

.remove-btn:hover {
  background-color: #fee;
  color: #c0392b;
  transform: scale(1.1);
}

.cart-summary {
  background: white;
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  height: fit-content;
  position: sticky;
  top: 100px;
}

.cart-summary h3 {
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid #4a90e2;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  color: #666;
}

.summary-row.total {
  font-size: 1.2rem;
  font-weight: bold;
  color: #333;
  border-top: 2px solid #eee;
  padding-top: 15px;
  margin-top: 15px;
}

.checkout-btn {
  width: 100%;
  padding: 14px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 15px;
}

.checkout-btn:hover {
  background-color: #357abd;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(74, 144, 226, 0.3);
}

@media (max-width: 768px) {
  .cart-content {
    grid-template-columns: 1fr;
  }

  .cart-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }

  .item-info {
    width: 100%;
  }

  .quantity-controls {
    align-self: center;
  }

  .item-total {
    align-self: flex-end;
  }

  .cart-summary {
    position: static;
  }
}

@media (max-width: 480px) {
  .qty-btn {
    width: 32px;
    height: 32px;
    font-size: 0.9rem;
  }

  .quantity {
    min-width: 25px;
  }
}
</style>
