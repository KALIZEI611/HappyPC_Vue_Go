<template>
  <div class="app">
    <Navbar :cart-count="cartCount" />
    <main class="main-content">
      <router-view
        :cart="cart"
        @add-to-cart="addToCart"
        @update-cart="updateCart"
        @remove-from-cart="removeFromCart"
      />
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import Navbar from "./components/Navbar.vue";

const cart = ref([]);

const cartCount = computed(() =>
  cart.value.reduce((acc, item) => acc + item.quantity, 0)
);

const addToCart = (product) => {
  const existing = cart.value.find((item) => item.product.id === product.id);
  if (existing) {
    existing.quantity += 1;
  } else {
    cart.value.push({ product, quantity: 1 });
  }
};

const updateCart = (productId, newQuantity) => {
  if (newQuantity <= 0) {
    cart.value = cart.value.filter((item) => item.product.id !== productId);
  } else {
    const item = cart.value.find((item) => item.product.id === productId);
    if (item) item.quantity = newQuantity;
  }
};

const removeFromCart = (productId) => {
  cart.value = cart.value.filter((item) => item.product.id !== productId);
};
</script>

<style scoped>
.app {
  min-height: 100vh;
  background-color: #f5f5f5;
}
.main-content {
  padding-top: 80px;
  padding-bottom: 60px;
}
</style>
