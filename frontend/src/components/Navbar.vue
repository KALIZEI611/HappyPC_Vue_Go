<template>
  <nav class="navbar">
    <div class="nav-container">
      <div class="nav-logo" @click="goToHome">
        <i class="fas fa-microchip"></i>
        <span>HappyPC</span>
      </div>

      <div class="nav-search">
        <input
          type="text"
          v-model="localSearchQuery"
          placeholder="Поиск запчастей..."
          @keyup.enter="performSearch"
        />
        <i class="fas fa-search" @click="performSearch"></i>
      </div>

      <ul class="nav-menu">
        <li v-for="item in displayMenuItems" :key="item.id">
          <a href="#" @click.prevent="handleMenuItemClick(item)">
            <i :class="item.icon"></i>
            <span>{{ item.name }}</span>
          </a>
        </li>
      </ul>

      <div class="nav-auth" v-if="user">
        <button @click="logout" class="logout-btn">Выйти</button>
      </div>

      <router-link to="/cart" class="nav-cart">
        <i class="fas fa-shopping-cart"></i>
        <span class="cart-badge" v-if="cartCount > 0">{{ cartCount }}</span>
      </router-link>

      <button class="mobile-menu-btn" @click="toggleMobileMenu">
        <i class="fas fa-bars"></i>
      </button>
    </div>

    <div v-if="mobileMenuOpen" class="mobile-menu">
      <ul>
        <li v-for="item in displayMenuItems" :key="item.id">
          <a href="#" @click.prevent="handleMenuItemClick(item)">
            <i :class="item.icon"></i>
            <span>{{ item.name }}</span>
          </a>
        </li>
        <li v-if="user">
          <button @click="logout" class="mobile-logout-btn">Выйти</button>
        </li>
      </ul>
    </div>
  </nav>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import { clearAllCaches, user, fetchUser, clearUser } from "../utils/cache";

const props = defineProps({
  cartCount: { type: Number, default: 0 },
});

const router = useRouter();
const route = useRoute();
const localSearchQuery = ref("");
const mobileMenuOpen = ref(false);

const baseMenuItems = [
  { id: 1, name: "Категории", icon: "fas fa-th-large", action: "home" },
  { id: 2, name: "Сборка ПК", icon: "fas fa-tv", action: "pc-assembly" },
  { id: 3, name: "Профиль", icon: "fas fa-user", action: "profile" },
];

const displayMenuItems = computed(() => {
  if (user.value) {
    return baseMenuItems.map((item) => {
      if (item.action === "profile") {
        return { ...item, name: user.value.username };
      }
      return item;
    });
  }
  return baseMenuItems;
});

const logout = async () => {
  try {
    await axios.post("/api/logout");
    clearUser(); // сбрасываем реактивного пользователя
    clearAllCaches(); // очищаем все кэши
    router.push("/");
  } catch (err) {
    console.error("Ошибка выхода:", err);
  }
};

const performSearch = () => {
  const query = localSearchQuery.value.trim();
  if (query) {
    router.push({ path: "/search", query: { q: query } });
    localSearchQuery.value = "";
  }
};

const handleMenuItemClick = (item) => {
  mobileMenuOpen.value = false;
  switch (item.action) {
    case "home":
      router.push("/");
      break;
    case "pc-assembly":
      console.log("Переход на страницу сборки ПК");
      break;
    case "profile":
      if (user.value) {
        router.push("/profile");
      } else {
        router.push("/login");
      }
      break;
  }
};

const goToHome = () => router.push("/");
const toggleMobileMenu = () => (mobileMenuOpen.value = !mobileMenuOpen.value);

onMounted(async () => {
  await fetchUser();
});
</script>

<style scoped>
.navbar {
  background: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  position: fixed;
  width: 100%;
  top: 0;
  z-index: 1000;
  border-bottom: 2px solid #4a90e2;
}

.nav-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1.5rem;
  font-weight: bold;
  color: #2c3e50;
  cursor: pointer;
  transition: opacity 0.3s;
}

.nav-logo:hover {
  opacity: 0.8;
}

.nav-logo i {
  font-size: 2rem;
  color: #4a90e2;
}

.nav-search {
  flex: 1;
  max-width: 400px;
  margin: 0 2rem;
  position: relative;
}

.nav-search input {
  width: 100%;
  padding: 0.8rem 1rem;
  padding-right: 2.5rem;
  border: 2px solid #e0e0e0;
  border-radius: 25px;
  font-size: 1rem;
  transition: all 0.3s;
}

.nav-search input:focus {
  outline: none;
  border-color: #4a90e2;
  box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.1);
}

.nav-search i {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #999;
}

.nav-menu {
  display: flex;
  list-style: none;
  gap: 2rem;
}

.nav-menu a {
  text-decoration: none;
  color: #555;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 5px;
  transition: color 0.3s;
  cursor: pointer;
}

.nav-menu a:hover {
  color: #4a90e2;
}

.nav-menu i {
  font-size: 1.1rem;
}

.nav-auth {
  margin-left: 1rem;
  display: flex;
  align-items: center;
}

.logout-btn {
  background: none;
  border: none;
  color: #e74c3c;
  cursor: pointer;
  font-size: 0.9rem;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  transition: background 0.2s;
}

.logout-btn:hover {
  background: #fee;
}

.nav-cart {
  position: relative;
  cursor: pointer;
  margin-left: 1rem;
}

.nav-cart i {
  font-size: 1.5rem;
  color: #555;
  transition: color 0.3s;
}

.nav-cart:hover i {
  color: #4a90e2;
}

.cart-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background-color: #e74c3c;
  color: white;
  font-size: 0.75rem;
  font-weight: bold;
  min-width: 18px;
  height: 18px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
}

.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #555;
}

.mobile-menu {
  display: none;
  background: white;
  padding: 1rem;
  border-top: 1px solid #eee;
}

.mobile-menu ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.mobile-menu li {
  margin: 0.5rem 0;
}

.mobile-menu a {
  text-decoration: none;
  color: #555;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0.5rem;
  cursor: pointer;
}

.mobile-menu a:hover {
  background: #f8f9fa;
  color: #4a90e2;
}

.mobile-logout-btn {
  width: 100%;
  background: none;
  border: none;
  color: #e74c3c;
  cursor: pointer;
  font-size: 1rem;
  padding: 0.5rem;
  text-align: left;
}

.mobile-logout-btn:hover {
  background: #fee;
}

@media (max-width: 1024px) {
  .nav-menu {
    display: none;
  }
  .mobile-menu-btn {
    display: block;
  }
  .mobile-menu {
    display: block;
  }
  .nav-search {
    max-width: 250px;
  }
  .nav-auth {
    display: none;
  }
}

@media (max-width: 768px) {
  .nav-container {
    padding: 1rem;
  }
  .nav-logo span {
    display: none;
  }
  .nav-search {
    margin: 0 1rem;
    max-width: 150px;
  }
}
</style>
