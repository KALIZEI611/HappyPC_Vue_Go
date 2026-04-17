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
            <button @click="handleLogout" class="nav-btn logout-btn">
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

          <div v-else-if="activeTab === 'favorites'" class="tab-content">
            <h2>Избранное</h2>

            <div v-if="favoritesLoading" class="loading">Загрузка...</div>
            <div v-else-if="favorites.length === 0" class="no-favorites">
              <i class="far fa-heart"></i>
              <p>У вас пока нет избранных товаров</p>
              <router-link to="/" class="browse-link">Перейти к покупкам</router-link>
            </div>
            <div v-else class="favorites-grid">
              <div v-for="item in favorites" :key="item.id" class="favorite-item">
                <button
                  @click="removeFromFavorites(item.product_id || item.Product?.id)"
                  class="remove-favorite"
                  title="Удалить из избранного"
                >
                  <i class="fas fa-times"></i>
                </button>

                <router-link
                  :to="'/product/' + (item.product_id || item.Product?.id)"
                  class="favorite-link"
                >
                  <img
                    :src="item.product?.image || item.Product?.image"
                    :alt="item.product?.name || item.Product?.name"
                  />
                  <h3>{{ item.product?.name || item.Product?.name }}</h3>
                  <div class="price">
                    {{ (item.product?.price || item.Product?.price).toLocaleString() }} ₽
                  </div>
                </router-link>

                <!-- Контролы количества -->
                <div class="favorite-actions">
                  <div
                    v-if="getQuantity(item.product_id || item.Product?.id) > 0"
                    class="quantity-controls"
                  >
                    <button
                      @click="
                        updateCartQuantity(
                          item.product_id || item.Product?.id,
                          getQuantity(item.product_id || item.Product?.id) - 1
                        )
                      "
                      class="qty-btn"
                    >
                      <i class="fas fa-minus"></i>
                    </button>
                    <span class="quantity">{{
                      getQuantity(item.product_id || item.Product?.id)
                    }}</span>
                    <button
                      @click="
                        updateCartQuantity(
                          item.product_id || item.Product?.id,
                          getQuantity(item.product_id || item.Product?.id) + 1
                        )
                      "
                      class="qty-btn"
                    >
                      <i class="fas fa-plus"></i>
                    </button>
                  </div>
                  <button
                    v-else
                    @click="addToCartFromFavorites(item.product || item.Product)"
                    class="add-to-cart-fav"
                  >
                    <i class="fas fa-cart-plus"></i> В корзину
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div v-else-if="activeTab === 'feedback'" class="tab-content">
            <h2>Обратная связь</h2>

            <div class="feedback-container">
              <div class="feedback-form-wrapper">
                <form @submit.prevent="sendFeedback" class="feedback-form">
                  <div class="form-group-modern">
                    <label class="form-label">
                      <i class="fas fa-tag"></i>
                      <span>Тема обращения</span>
                    </label>
                    <div class="select-wrapper">
                      <select v-model="feedbackForm.subject" required class="form-select">
                        <option value="" disabled>Выберите тему</option>
                        <option value="order">
                          <i class="fas fa-shopping-bag"></i> Вопрос о заказе
                        </option>
                        <option value="product">
                          <i class="fas fa-microchip"></i> Проблема с товаром
                        </option>
                        <option value="delivery">
                          <i class="fas fa-truck"></i> Доставка
                        </option>
                        <option value="payment">
                          <i class="fas fa-credit-card"></i> Оплата
                        </option>
                        <option value="other">
                          <i class="fas fa-comment"></i> Другое
                        </option>
                      </select>
                      <i class="fas fa-chevron-down select-arrow"></i>
                    </div>
                  </div>

                  <div class="form-group-modern">
                    <label class="form-label">
                      <i class="fas fa-envelope"></i>
                      <span>Сообщение</span>
                    </label>
                    <div class="textarea-wrapper">
                      <textarea
                        v-model="feedbackForm.message"
                        rows="6"
                        placeholder="Опишите вашу проблему или вопрос..."
                        required
                        class="form-textarea"
                      ></textarea>
                      <div class="textarea-footer">
                        <span class="char-count"
                          >{{ feedbackForm.message.length }}/1000</span
                        >
                        <span class="markdown-hint">
                          <i class="fas fa-markdown"></i> Поддерживается Markdown
                        </span>
                      </div>
                    </div>
                  </div>

                  <div class="form-group-modern checkbox-group">
                    <label class="checkbox-label">
                      <input type="checkbox" v-model="feedbackForm.copyToEmail" />
                      <span class="checkbox-custom"></span>
                      <span class="checkbox-text">
                        <i class="fas fa-paper-plane"></i>
                        Отправить копию на email
                      </span>
                    </label>
                  </div>

                  <div class="form-actions">
                    <button
                      type="submit"
                      :disabled="feedbackLoading"
                      class="submit-btn-modern"
                    >
                      <i
                        :class="
                          feedbackLoading
                            ? 'fas fa-spinner fa-spin'
                            : 'fas fa-paper-plane'
                        "
                      ></i>
                      {{ feedbackLoading ? "Отправка..." : "Отправить сообщение" }}
                    </button>
                  </div>

                  <div v-if="feedbackSuccess" class="alert-success">
                    <i class="fas fa-check-circle"></i>
                    <span>{{ feedbackSuccess }}</span>
                  </div>

                  <div v-if="feedbackError" class="alert-error">
                    <i class="fas fa-exclamation-circle"></i>
                    <span>{{ feedbackError }}</span>
                  </div>
                </form>
              </div>

              <div class="feedback-info-modern">
                <h3>
                  <i class="fas fa-question-circle"></i>
                  Часто задаваемые вопросы
                </h3>

                <div class="faq-list">
                  <div class="faq-item" v-for="(faq, index) in faqList" :key="index">
                    <div class="faq-question" @click="toggleFaq(index)">
                      <i class="fas fa-question"></i>
                      <span>{{ faq.question }}</span>
                      <i
                        :class="faq.open ? 'fas fa-chevron-up' : 'fas fa-chevron-down'"
                      ></i>
                    </div>
                    <div class="faq-answer" v-show="faq.open">
                      <i class="fas fa-reply"></i>
                      <p>{{ faq.answer }}</p>
                    </div>
                  </div>
                </div>

                <div class="contact-info">
                  <h4>
                    <i class="fas fa-headset"></i>
                    Служба поддержки
                  </h4>
                  <div class="contact-details">
                    <div class="contact-item">
                      <i class="fas fa-phone-alt"></i>
                      <span>+7 (800) 123-45-67</span>
                    </div>
                    <div class="contact-item">
                      <i class="fas fa-clock"></i>
                      <span>Пн-Пт: 9:00 - 20:00</span>
                    </div>
                    <div class="contact-item">
                      <i class="fas fa-envelope"></i>
                      <span>support@happypc.ru</span>
                    </div>
                  </div>
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
import { ref, reactive, onMounted, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import { user, fetchUser, logout } from "../utils/cache";
import { allProductsCache } from "../utils/cache";
import favoritesService from "../services/favoritesService";
import cartService from "../services/cartService";

// Определяем emit
const emit = defineEmits(["add-to-cart", "update-cart", "cart-cleared"]);

const router = useRouter();
const route = useRoute();
const activeTab = ref("profile");
const userData = ref(null);
const loading = ref(true);
const orders = ref([]);
const ordersLoading = ref(false);
const builds = ref([]);
const buildsLoading = ref(false);
const favorites = ref([]);
const favoritesLoading = ref(false);
const cart = ref([]); // Добавляем корзину

// Функция для получения количества товара в корзине
const getQuantity = (productId) => {
  const cartItem = cart.value.find(
    (item) => item.product?.id === productId || item.product_id === productId
  );
  return cartItem ? cartItem.quantity : 0;
};

// Функция для обновления количества в корзине
const updateCartQuantity = async (productId, newQuantity) => {
  if (!user.value) {
    router.push("/login");
    return;
  }

  try {
    if (newQuantity <= 0) {
      await cartService.removeCartItem(productId);
    } else {
      await cartService.updateCartItem(productId, newQuantity);
    }
    await fetchCart();
    emit("update-cart", productId, newQuantity);
  } catch (err) {
    console.error("Ошибка обновления корзины:", err);
  }
};

// Функция для загрузки корзины
const fetchCart = async () => {
  if (!user.value) return;
  try {
    const items = await cartService.getCart();
    cart.value = items;
  } catch (err) {
    console.error("Ошибка загрузки корзины:", err);
  }
};

const feedbackForm = reactive({
  subject: "",
  message: "",
  copyToEmail: false,
});
const feedbackLoading = ref(false);
const feedbackSuccess = ref("");
const feedbackError = ref("");

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

const toggleFaq = (index) => {
  faqList.value[index].open = !faqList.value[index].open;
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

const fetchFavorites = async () => {
  if (!user.value) return;
  favoritesLoading.value = true;
  try {
    const data = await favoritesService.getFavorites();
    favorites.value = data;
  } catch (err) {
    console.error("Ошибка загрузки избранного:", err);
  } finally {
    favoritesLoading.value = false;
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
      const fullProduct = allProducts.find((p) => p.id === comp.id);
      if (fullProduct) {
        buildToLoad[type] = fullProduct;
      } else {
        buildToLoad[type] = comp;
      }
    }
  }

  localStorage.setItem("pcBuilderBuild", JSON.stringify(buildToLoad));
  router.push("/pc-builder");
};

const removeFromFavorites = async (productId) => {
  if (!productId) return;
  try {
    await favoritesService.removeFromFavorites(productId);
    await fetchFavorites();
  } catch (err) {
    console.error("Ошибка удаления из избранного:", err);
  }
};

const addToCartFromFavorites = async (product) => {
  if (!product) return;
  emit("add-to-cart", product);
  // Обновляем корзину после добавления
  setTimeout(() => {
    fetchCart();
  }, 500);
};

const sendFeedback = async () => {
  feedbackLoading.value = true;
  feedbackSuccess.value = "";
  feedbackError.value = "";

  try {
    await axios.post("/api/feedback", feedbackForm);
    feedbackSuccess.value = "Сообщение отправлено! Мы ответим вам в ближайшее время.";
    feedbackForm.subject = "";
    feedbackForm.message = "";
    feedbackForm.copyToEmail = false;
  } catch (err) {
    feedbackError.value = err.response?.data || "Ошибка отправки сообщения";
  } finally {
    feedbackLoading.value = false;
  }
};

const handleLogout = async () => {
  await logout();
  emit("cart-cleared");
};

watch(activeTab, (newTab) => {
  if (newTab === "favorites") {
    fetchFavorites();
  }
  if (newTab === "orders") {
    fetchOrders();
  }
  if (newTab === "builds") {
    fetchBuilds();
  }
});

onMounted(async () => {
  await fetchUser();
  userData.value = user.value;
  loading.value = false;
  await fetchCart(); // Загружаем корзину при монтировании
  if (activeTab.value === "orders") {
    await fetchOrders();
  }
});
</script>

<style scoped>
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

/* Сайдбар */
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

.logout-btn {
  margin-top: auto;
  color: #e74c3c;
}

.logout-btn:hover {
  background: #ffebee;
  color: #c0392b;
}

/* Основной контент */
.profile-content {
  flex: 1;
  padding: 32px;
  overflow-y: auto;
  max-height: calc(100vh - 100px);
}

.tab-content h2 {
  margin-bottom: 24px;
  color: #333;
  font-size: 1.5rem;
  border-bottom: 2px solid #4a90e2;
  padding-bottom: 10px;
  display: inline-block;
}

/* Профиль */
.profile-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
  background: #f8f9fa;
  padding: 24px;
  border-radius: 12px;
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

/* Загрузка и ошибки */
.loading,
.error {
  text-align: center;
  padding: 60px 20px;
  color: #666;
  background: #f8f9fa;
  border-radius: 12px;
}

.loading {
  position: relative;
}

.loading::after {
  content: "";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #4a90e2;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: translate(-50%, -50%) rotate(0deg);
  }
  100% {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}

/* Заказы */
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

.order-item:last-child {
  border-bottom: none;
}

.item-name {
  color: #4a90e2;
  text-decoration: none;
  font-weight: 500;
}

.item-name:hover {
  text-decoration: underline;
}

.item-quantity,
.item-price,
.item-total {
  color: #555;
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

.delivery-info,
.payment-info {
  display: flex;
  gap: 5px;
}

.order-total {
  font-weight: bold;
  font-size: 1.1rem;
  color: #2c3e50;
}

/* Сборки */
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
  transition: box-shadow 0.2s;
}

.build-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
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
  color: #333;
}

.delete-build-btn {
  background: none;
  border: none;
  color: #e74c3c;
  cursor: pointer;
  font-size: 1.1rem;
  padding: 5px 10px;
  border-radius: 6px;
  transition: all 0.2s;
}

.delete-build-btn:hover {
  background: #fee;
  transform: scale(1.05);
}

.build-components {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 15px;
}

.build-component {
  font-size: 0.9rem;
  padding: 4px 0;
  color: #555;
}

.build-component strong {
  color: #333;
  min-width: 140px;
  display: inline-block;
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
  transition: background-color 0.2s;
}

.load-build-btn:hover {
  background-color: #357abd;
}

/* Избранное */
.favorites-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 20px;
}

.favorite-item {
  position: relative;
  background: white;
  border-radius: 12px;
  padding: 16px;
  border: 1px solid #e0e0e0;
  transition: transform 0.2s, box-shadow 0.2s;
}

.favorite-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.remove-favorite {
  position: absolute;
  top: 8px;
  right: 8px;
  background: #ffebee;
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  cursor: pointer;
  color: #e74c3c;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  z-index: 1;
}

.remove-favorite:hover {
  background: #e74c3c;
  color: white;
  transform: scale(1.05);
}

.favorite-link {
  text-decoration: none;
  text-align: center;
  display: block;
}

.favorite-link img {
  width: 100%;
  height: 150px;
  object-fit: contain;
  margin-bottom: 12px;
}

.favorite-link h3 {
  font-size: 0.9rem;
  color: #333;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price {
  font-size: 1.1rem;
  font-weight: bold;
  color: #2c3e50;
}

.favorite-actions {
  margin-top: 12px;
}

.favorite-actions .quantity-controls {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: #f8f9fa;
  border-radius: 30px;
  padding: 6px 10px;
  border: 1px solid #e0e0e0;
  width: fit-content;
  margin: 0 auto;
}

.favorite-actions .qty-btn {
  width: 32px;
  height: 32px;
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
  font-size: 0.9rem;
}

.favorite-actions .qty-btn:hover {
  background: #4a90e2;
  color: white;
  transform: scale(1.05);
}

.favorite-actions .quantity {
  font-weight: 600;
  color: #333;
  min-width: 28px;
  text-align: center;
  font-size: 1rem;
}

.add-to-cart-fav {
  width: 100%;
  padding: 10px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 0.9rem;
  transition: background-color 0.3s;
}

.add-to-cart-fav:hover {
  background-color: #357abd;
}

.no-favorites {
  text-align: center;
  padding: 60px 20px;
  background: #f8f9fa;
  border-radius: 12px;
}

.no-favorites i {
  font-size: 4rem;
  color: #ccc;
  margin-bottom: 20px;
}

.no-favorites p {
  color: #666;
  font-size: 1.1rem;
  margin-bottom: 20px;
}

.browse-link {
  display: inline-block;
  margin-top: 20px;
  padding: 12px 24px;
  background-color: #4a90e2;
  color: white;
  text-decoration: none;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.browse-link:hover {
  background-color: #357abd;
}

/* Обратная связь */
.feedback-form {
  max-width: 600px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
}

.form-group select,
.form-group textarea,
.form-group input[type="text"] {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 0.95rem;
  transition: border-color 0.2s;
}

.form-group select:focus,
.form-group textarea:focus,
.form-group input[type="text"]:focus {
  outline: none;
  border-color: #4a90e2;
}

.form-group textarea {
  resize: vertical;
}

.feedback-info {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #e0e0e0;
}

.feedback-info h4 {
  margin-bottom: 15px;
  color: #333;
}

.feedback-info ul {
  list-style: none;
  padding: 0;
}

.feedback-info li {
  margin-bottom: 15px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

.feedback-info strong {
  display: block;
  margin-bottom: 5px;
  color: #4a90e2;
}

.feedback-info p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}

.submit-feedback-btn {
  padding: 12px 24px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.submit-feedback-btn:hover:not(:disabled) {
  background-color: #357abd;
}

.submit-feedback-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.success-message {
  margin-top: 15px;
  padding: 10px;
  background: #e8f5e9;
  color: #27ae60;
  border-radius: 8px;
  text-align: center;
}

.error-message {
  margin-top: 15px;
  padding: 10px;
  background: #ffebee;
  color: #e74c3c;
  border-radius: 8px;
  text-align: center;
}

/* Пустые состояния */
.no-orders,
.no-builds {
  text-align: center;
  padding: 60px 20px;
  background: #f8f9fa;
  border-radius: 12px;
}

.no-orders p,
.no-builds p {
  color: #666;
  font-size: 1.1rem;
}

/* Адаптивность */
@media (max-width: 1024px) {
  .profile-sidebar {
    width: 220px;
  }

  .profile-content {
    padding: 24px;
  }
}

@media (max-width: 768px) {
  .profile-layout {
    flex-direction: column;
  }

  .profile-sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e0e0e0;
    padding: 10px 0;
  }

  .profile-nav {
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: center;
  }

  .nav-btn {
    width: auto;
    padding: 10px 16px;
  }

  .nav-btn.active {
    border-left: none;
    border-bottom: 2px solid #4a90e2;
  }

  .logout-btn {
    margin-top: 0;
  }

  .profile-content {
    padding: 20px;
    max-height: none;
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

  .favorites-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
  }

  .build-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .build-footer {
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
  }

  .build-component strong {
    min-width: 120px;
  }
}

@media (max-width: 480px) {
  .container {
    padding: 0 15px;
  }

  .profile-content {
    padding: 15px;
  }

  .tab-content h2 {
    font-size: 1.2rem;
  }

  .profile-info {
    padding: 15px;
  }

  .favorites-grid {
    grid-template-columns: 1fr;
  }

  .order-card {
    padding: 15px;
  }

  .build-card {
    padding: 15px;
  }

  .feedback-form {
    max-width: 100%;
  }

  .submit-feedback-btn {
    width: 100%;
  }
}
</style>
