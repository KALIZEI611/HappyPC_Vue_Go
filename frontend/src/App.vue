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
        @add-to-build="handleAddToBuild"
        @cart-cleared="handleCartCleared"
      />
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import Navbar from "./components/Navbar.vue";
import Breadcrumbs from "./components/Breadcrumbs.vue";
import { useBreadcrumbs } from "./composables/useBreadcrumbs";
import { cartService } from "./services/cartService";
import { user, fetchUser } from "./utils/cache";

const { breadcrumbItems } = useBreadcrumbs();
const route = useRoute();
const router = useRouter();

const cart = ref([]);

const cartCount = computed(() =>
  cart.value.reduce((acc, item) => acc + item.quantity, 0)
);

const showBreadcrumbs = computed(() => {
  return ["Search", "Category", "Product"].includes(route.name);
});

const fetchCart = async () => {
  if (!user.value) return;
  try {
    const items = await cartService.getCart();
    cart.value = items;
  } catch (err) {
    console.error("Ошибка загрузки корзины:", err);
  }
};

const handleAddToBuild = (product) => {
  const returnKey = sessionStorage.getItem("pcBuilderReturnKey");
  if (returnKey) {
    router.push(`/pc-builder?product_id=${product.id}`);
  } else {
    console.log("Режим выбора компонента не активен");
    alert("Сначала выберите компонент для замены на странице сборки ПК");
  }
};

const addToCart = async (product) => {
  if (!user.value) {
    router.push("/login");
    return;
  }
  try {
    await cartService.addToCart(product.id, 1);
    await fetchCart();
  } catch (err) {
    console.error("Ошибка добавления в корзину:", err);
  }
};

const handleCartCleared = async () => {
  cart.value = [];
  await fetchCart();
};

const updateCart = async (productId, newQuantity) => {
  if (!user.value) return;
  try {
    if (newQuantity <= 0) {
      await cartService.removeCartItem(productId);
    } else {
      await cartService.updateCartItem(productId, newQuantity);
    }
    await fetchCart();
  } catch (err) {
    console.error("Ошибка обновления корзины:", err);
  }
};

const removeFromCart = async (productId) => {
  if (!user.value) return;
  try {
    await cartService.removeCartItem(productId);
    await fetchCart();
  } catch (err) {
    console.error("Ошибка удаления из корзины:", err);
  }
};
watch(
  () => route.path,
  (newPath) => {
    if (newPath === "/cart" && user.value) {
      fetchCart();
    }
  }
);
watch(user, async (newUser) => {
  if (newUser) {
    await fetchCart();
  } else {
    cart.value = [];
  }
});

onMounted(async () => {
  try {
    await fetchUser();
    if (user.value) {
      await fetchCart();
    }
  } catch {}
});
</script>

<style scoped>
.app {
  min-height: 100vh;
  background-color: #f5f5f5;
}
.main-content {
  padding-top: 20px;
  padding-bottom: 60px;
}
.main-content.no-breadcrumbs {
  padding-top: 20px;
}
</style>
