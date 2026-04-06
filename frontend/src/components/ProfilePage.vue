<template>
  <div class="profile-page">
    <div class="container">
      <div class="profile-layout">
        <aside class="profile-sidebar">
          <nav class="profile-nav">
            <button
              v-for="item in menuItems"
              :key="item.id"
              :class="['nav-btn', { active: activeTab === item.id }]"
              @click="activeTab = item.id"
            >
              <i :class="item.icon"></i>
              <span>{{ item.name }}</span>
            </button>
            <button @click="logout" class="nav-btn logout-btn">
              <i class="fas fa-sign-out-alt"></i>
              <span>Выйти</span>
            </button>
          </nav>
        </aside>

        <main class="profile-content">
          <div v-if="activeTab === 'profile'" class="tab-content">
            <h2>Мой профиль</h2>
            <div v-if="loading" class="loading">Загрузка...</div>
            <div v-else-if="user" class="profile-info">
              <div class="info-row">
                <strong>Имя пользователя:</strong>
                <span>{{ user.username }}</span>
              </div>
              <div class="info-row">
                <strong>Email:</strong>
                <span>{{ user.email }}</span>
              </div>
              <div class="info-row">
                <strong>Дата регистрации:</strong>
                <span>{{ formatDate(user.created_at) }}</span>
              </div>
            </div>
            <div v-else class="error">Не удалось загрузить данные профиля</div>
          </div>

          <div v-else-if="activeTab === 'orders'" class="tab-content">
            <h2>Мои заказы</h2>
            <div v-if="ordersLoading" class="loading">Загрузка заказов...</div>
            <div v-else-if="orders.length === 0" class="no-orders">
              <p>У вас пока нет заказов.</p>
            </div>
            <div v-else class="orders-list">
              <div v-for="order in orders" :key="order.id" class="order-card">
                <div class="order-header">
                  <span class="order-id">Заказ №{{ order.id }}</span>
                  <span class="order-date">{{ formatDate(order.created_at) }}</span>
                  <span class="order-status" :class="'status-' + order.status">
                    {{ getStatusText(order.status) }}
                  </span>
                </div>
                <div class="order-items">
                  <div v-for="item in order.items" :key="item.id" class="order-item">
                    <router-link :to="'/product/' + item.product_id" class="item-name">
                      {{ getProductName(item.product_id) }}
                    </router-link>
                    <span class="item-quantity">{{ item.quantity }} шт.</span>
                    <span class="item-price">{{ item.price.toLocaleString() }} ₽</span>
                    <span class="item-total"
                      >{{ (item.price * item.quantity).toLocaleString() }} ₽</span
                    >
                  </div>
                </div>
                <div class="order-footer">
                  <div class="delivery-info">
                    <span>Доставка: {{ getDeliveryText(order.delivery_method) }}</span>
                    <span v-if="order.delivery_address"
                      >({{ order.delivery_address }})</span
                    >
                  </div>
                  <div class="payment-info">
                    Оплата: {{ getPaymentText(order.payment_method) }}
                  </div>
                  <div class="order-total">
                    Итого: {{ order.total_price.toLocaleString() }} ₽
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else-if="activeTab === 'builds'" class="tab-content">
            <h2>Мои сборки</h2>
            <div v-if="buildsLoading" class="loading">Загрузка сборок...</div>
            <div v-else-if="builds.length === 0" class="no-builds">
              <p>У вас пока нет сохранённых сборок.</p>
            </div>
            <div v-else class="builds-list">
              <div v-for="build in builds" :key="build.id" class="build-card">
                <div class="build-header">
                  <h3>{{ build.name }}</h3>
                  <button @click="deleteBuild(build.id)" class="delete-build-btn">
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
                <div class="build-components">
                  <div
                    v-for="(comp, type) in build.components"
                    :key="type"
                    class="build-component"
                  >
                    <strong>{{ getComponentTypeName(type) }}:</strong> {{ comp.name }}
                  </div>
                </div>
                <div class="build-footer">
                  <span class="build-date">{{ formatDate(build.created_at) }}</span>
                  <button @click="loadBuild(build.components)" class="load-build-btn">
                    Загрузить сборку
                  </button>
                </div>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import axios from "axios";
import { user, fetchUser, logout } from "../utils/cache";
import { allProductsCache } from "../utils/cache";
import { useRouter, useRoute } from "vue-router";

const activeTab = ref("profile");
const userData = ref(null);
const loading = ref(true);
const orders = ref([]);
const ordersLoading = ref(false);
const builds = ref([]);
const buildsLoading = ref(false);
const router = useRouter();
const route = useRoute();

const menuItems = [
  { id: "profile", name: "Мой профиль", icon: "fas fa-user" },
  { id: "orders", name: "Мои заказы", icon: "fas fa-shopping-bag" },
  { id: "builds", name: "Мои сборки", icon: "fas fa-tv" },
  { id: "favorites", name: "Избранное", icon: "fas fa-heart" },
  { id: "feedback", name: "Обратная связь", icon: "fas fa-envelope" },
];

const getComponentTypeName = (type) => {
  const names = {
    cpu: "Процессор",
    gpu: "Видеокарта",
    motherboard: "Материнская плата",
    ram: "Оперативная память",
    psu: "Блок питания",
    cooler: "Кулер",
    storage: "SSD",
    case: "Корпус",
  };
  return names[type] || type;
};

const formatDate = (dateString) => {
  if (!dateString) return "";
  const date = new Date(dateString);
  return date.toLocaleDateString("ru-RU", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const getStatusText = (status) => {
  const statuses = {
    pending: "В обработке",
    paid: "Оплачен",
    shipped: "Отправлен",
    delivered: "Доставлен",
    cancelled: "Отменён",
  };
  return statuses[status] || status;
};

const getDeliveryText = (method) => {
  return method === "pickup" ? "Самовывоз" : "Доставка";
};

const getPaymentText = (method) => {
  const methods = {
    online: "Онлайн",
    cash: "При получении",
    credit: "В кредит",
  };
  return methods[method] || method;
};

const getProductName = (productId) => {
  const allProducts = allProductsCache.get();
  if (allProducts) {
    const product = allProducts.find((p) => p.id === productId);
    if (product) return product.name;
  }
  return `Товар #${productId}`;
};

const fetchOrders = async () => {
  if (!user.value) return;
  ordersLoading.value = true;
  try {
    const { data } = await axios.get("/api/orders");
    orders.value = data;
  } catch (err) {
    console.error("Ошибка загрузки заказов:", err);
  } finally {
    ordersLoading.value = false;
  }
};

const fetchBuilds = async () => {
  if (!user.value) return;
  buildsLoading.value = true;
  try {
    const { data } = await axios.get("/api/builds");
    builds.value = data;
  } catch (err) {
    console.error("Ошибка загрузки сборок:", err);
  } finally {
    buildsLoading.value = false;
  }
};

const deleteBuild = async (buildId) => {
  if (!confirm("Удалить эту сборку?")) return;
  try {
    await axios.delete(`/api/builds/${buildId}`);
    await fetchBuilds();
  } catch (err) {
    console.error("Ошибка удаления сборки:", err);
    alert("Не удалось удалить сборку");
  }
};

const loadBuild = async (components) => {
  // Загружаем все товары из кэша
  let allProducts = allProductsCache.get();
  if (!allProducts) {
    try {
      const { data: cats } = await axios.get("/categories");
      const promises = cats.map(async (cat) => {
        const res = await axios.get(`/${cat.url_key}`);
        return res.data;
      });
      const results = await Promise.all(promises);
      allProducts = results.flat();
      allProductsCache.set(allProducts);
    } catch (err) {
      console.error("Ошибка загрузки товаров:", err);
      alert("Не удалось загрузить товары");
      return;
    }
  }

  const buildToLoad = {};
  for (const [type, comp] of Object.entries(components)) {
    if (comp && comp.id) {
      // Ищем полный объект товара по ID
      const fullProduct = allProducts.find((p) => p.id === comp.id);
      if (fullProduct) {
        buildToLoad[type] = fullProduct;
      } else {
        // Если не нашли, используем сохранённые данные
        buildToLoad[type] = comp;
      }
    }
  }

  localStorage.setItem("pcBuilderBuild", JSON.stringify(buildToLoad));
  router.push("/pc-builder");
};
watch(activeTab, (newTab) => {
  if (newTab === "builds") {
    fetchBuilds();
  }
});

watch(activeTab, (newTab) => {
  if (newTab === "orders") {
    fetchOrders();
  }
});

onMounted(async () => {
  await fetchUser();
  userData.value = user.value;
  loading.value = false;
  if (activeTab.value === "orders") {
    await fetchOrders();
  }
});
</script>

<style scoped>
.builds-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.build-card {
  background: #f9f9f9;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #e0e0e0;
}

.build-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e0e0e0;
}

.build-header h3 {
  margin: 0;
  font-size: 1.2rem;
}

.delete-build-btn {
  background: none;
  border: none;
  color: #e74c3c;
  cursor: pointer;
  font-size: 1.1rem;
  padding: 5px 10px;
  border-radius: 6px;
}

.delete-build-btn:hover {
  background: #fee;
}

.build-components {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 15px;
}

.build-component {
  font-size: 0.9rem;
}

.build-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 10px;
  border-top: 1px solid #e0e0e0;
}

.build-date {
  color: #666;
  font-size: 0.8rem;
}

.load-build-btn {
  padding: 6px 12px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.load-build-btn:hover {
  background-color: #357abd;
}
.profile-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-top: 100px;
  padding-bottom: 40px;
  display: flex;
  flex-direction: column;
}

.container {
  flex: 1;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  width: 100%;
}

.profile-layout {
  display: flex;
  gap: 30px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  min-height: 500px;
  height: 100%;
}

.profile-sidebar {
  width: 260px;
  background: #f8f9fa;
  border-right: 1px solid #e0e0e0;
  padding: 20px 0;
  display: flex;
  flex-direction: column;
}

.profile-nav {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.nav-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 24px;
  background: none;
  border: none;
  width: 100%;
  text-align: left;
  font-size: 1rem;
  color: #555;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-btn i {
  width: 20px;
  font-size: 1.1rem;
}

.nav-btn:hover {
  background: #e9ecef;
  color: #4a90e2;
}

.nav-btn.active {
  background: #e9ecef;
  color: #4a90e2;
  border-left: 3px solid #4a90e2;
  font-weight: 500;
}

.profile-content {
  flex: 1;
  padding: 32px;
}

.tab-content h2 {
  margin-bottom: 24px;
  color: #333;
  font-size: 1.5rem;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-row {
  display: flex;
  border-bottom: 1px solid #e0e0e0;
  padding-bottom: 12px;
}

.info-row strong {
  width: 180px;
  color: #666;
  font-weight: 500;
}

.info-row span {
  color: #333;
  font-weight: 400;
}

.loading,
.error {
  text-align: center;
  padding: 40px;
  color: #666;
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.order-card {
  background: #f9f9f9;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #e0e0e0;
  transition: box-shadow 0.2s;
}

.order-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e0e0e0;
}

.order-id {
  font-weight: 600;
  font-size: 1.1rem;
  color: #333;
}

.order-date {
  color: #666;
  font-size: 0.9rem;
}

.order-status {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 500;
}

.status-pending {
  background: #fff3e0;
  color: #e67e22;
}
.status-paid {
  background: #e8f5e9;
  color: #27ae60;
}
.status-shipped {
  background: #e3f2fd;
  color: #2980b9;
}
.status-delivered {
  background: #e8f5e9;
  color: #27ae60;
}
.status-cancelled {
  background: #ffebee;
  color: #e74c3c;
}

.order-items {
  margin: 15px 0;
}

.order-item {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 15px;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
  font-size: 0.9rem;
}

.item-name {
  color: #4a90e2;
  text-decoration: none;
  font-weight: 500;
}
.item-name:hover {
  text-decoration: underline;
}

.order-footer {
  margin-top: 15px;
  padding-top: 10px;
  border-top: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
  font-size: 0.9rem;
  color: #555;
}

.order-total {
  font-weight: bold;
  font-size: 1.1rem;
  color: #2c3e50;
}

.no-orders {
  text-align: center;
  padding: 40px;
  color: #666;
}

@media (max-width: 768px) {
  .profile-layout {
    flex-direction: column;
  }
  .profile-sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e0e0e0;
  }
  .profile-nav {
    flex-direction: row;
    flex-wrap: wrap;
  }
  .nav-btn {
    width: auto;
    padding: 10px 16px;
  }
  .info-row {
    flex-direction: column;
    gap: 5px;
  }
  .info-row strong {
    width: auto;
  }
  .order-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  .order-item {
    grid-template-columns: 1fr;
    gap: 5px;
    border-bottom: 1px solid #eee;
    padding: 12px 0;
  }
  .order-footer {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
