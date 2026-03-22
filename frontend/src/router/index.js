import { createRouter, createWebHistory } from "vue-router";
import HomePage from "../components/HomePage.vue";
import CategoryPage from "../components/CategoryPage.vue";
import CartPage from "../components/CartPage.vue";
import SearchPage from "../components/SearchPage.vue";
import ProductPage from "../components/ProductPage.vue";
import LoginPage from "../components/LoginPage.vue";
import RegisterPage from "../components/RegisterPage.vue";

const routes = [
  { path: "/", name: "Home", component: HomePage },
  { path: "/category/:id", name: "Category", component: CategoryPage },
  { path: "/cart", name: "Cart", component: CartPage },
  {
    path: "/search",
    name: "Search",
    component: SearchPage,
    props: (route) => ({ searchQuery: route.query.q || "" }),
  },
  { path: "/product/:id", name: "Product", component: ProductPage },
  { path: "/login", name: "Login", component: LoginPage },
  { path: "/register", name: "Register", component: RegisterPage },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { top: 0 };
  },
});

export default router;