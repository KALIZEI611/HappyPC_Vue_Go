import { createRouter, createWebHistory } from "vue-router";
import axios from "axios";
import HomePage from "../components/HomePage.vue";
import CategoryPage from "../components/CategoryPage.vue";
import CartPage from "../components/CartPage.vue";
import SearchPage from "../components/SearchPage.vue";
import ProductPage from "../components/ProductPage.vue";
import LoginPage from "../components/LoginPage.vue";
import RegisterPage from "../components/RegisterPage.vue";
import ProfilePage from "../components/ProfilePage.vue";
import CheckoutPage from "../components/CheckoutPage.vue";
import OrderSuccessPage from "../components/OrderSuccessPage.vue";
import OrdersPage from "../components/OrdersPage.vue";

axios.defaults.withCredentials = true;

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
  {
    path: "/profile",
    name: "Profile",
    component: ProfilePage,
    meta: { requiresAuth: true },
  },
  
  { path: "/checkout", name: "Checkout", component: CheckoutPage, meta: { requiresAuth: true } },
  { path: "/order-success", name: "OrderSuccess", component: OrderSuccessPage },
  {path: "/profile/orders", name: "Orders", component: OrdersPage, meta: { requiresAuth: true }},
  
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { top: 0 };
  },
});

router.beforeEach(async (to, from, next) => {
  if (to.meta.requiresAuth) {
    try {
      const { data } = await axios.get("/api/me");
      if (data) {
        next();
      } else {
        next("/login");
      }
    } catch (err) {
      console.error("Auth guard error:", err);
      next("/login");
    }
  } else {
    next();
  }
});

export default router;
