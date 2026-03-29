<template>
  <div class="app">
    <Navbar :cart-count="cartCount" />
    <Breadcrumbs v-if="showBreadcrumbs" :items="breadcrumbItems" />
    <main class="main-content" :class="{ 'no-breadcrumbs': !showBreadcrumbs }">
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
import { useRoute } from "vue-router";
import Navbar from "./components/Navbar.vue";
import Breadcrumbs from "./components/Breadcrumbs.vue";
import { useBreadcrumbs } from "./composables/useBreadcrumbs";

const { breadcrumbItems } = useBreadcrumbs();
const route = useRoute();

const cart = ref([]);

const cartCount = computed(() =>
  cart.value.reduce((acc, item) => acc + item.quantity, 0)
);

const showBreadcrumbs = computed(() => {
  return ["Search", "Category", "Product"].includes(route.name);
});

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
  padding-top: 20px; /* отступ после хлебных крошек */
  padding-bottom: 60px;
}
.main-content.no-breadcrumbs {
  padding-top: 20px; /* оставляем небольшой отступ сверху, если хлебных крошек нет */
}
</style>
