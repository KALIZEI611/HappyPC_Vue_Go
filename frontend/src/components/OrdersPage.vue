<template>
  <div class="orders-page">
    <div class="container">
      <h1>Мои заказы</h1>
      <div v-if="loading" class="loading">Загрузка...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else-if="orders.length === 0" class="empty">У вас пока нет заказов.</div>
      <div v-else class="orders-list">
        <div v-for="order in orders" :key="order.id" class="order-card">
          <div class="order-header">
            <span class="order-id">Заказ №{{ order.id }}</span>
            <span class="order-date">{{ formatDate(order.created_at) }}</span>
            <span class="order-status">{{ getStatusText(order.status) }}</span>
          </div>
          <div class="order-items">
            <div v-for="item in order.items" :key="item.id" class="order-item">
              <img
                :src="item.product.image"
                :alt="item.product.name"
                class="item-image"
              />
              <div class="item-info">
                <div class="item-name">{{ item.product.name }}</div>
                <div class="item-quantity">× {{ item.quantity }}</div>
                <div class="item-price">
                  {{ (item.price * item.quantity).toLocaleString() }} ₽
                </div>
              </div>
            </div>
          </div>
          <div class="order-total">
            Итого: <strong>{{ order.total_price.toLocaleString() }} ₽</strong>
          </div>
          <div class="order-details">
            <div>
              <strong>Способ доставки:</strong>
              {{ getDeliveryText(order.delivery_method) }}
            </div>
            <div v-if="order.delivery_address">
              <strong>Адрес:</strong> {{ order.delivery_address }}
            </div>
            <div>
              <strong>Способ оплаты:</strong> {{ getPaymentText(order.payment_method) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import axios from "axios";
import api from "../api";

const orders = ref([]);
const loading = ref(true);
const error = ref(null);

const formatDate = (dateString) => {
  if (!dateString) return "";
  const date = new Date(dateString);
  return date.toLocaleDateString("ru-RU", {
    year: "numeric",
    month: "long",
    day: "numeric",
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

const fetchOrders = async () => {
  loading.value = true;
  error.value = null;
  try {
    const { data } = await api.get("/api/orders");
    orders.value = data;
  } catch (err) {
    error.value = err.response?.data || "Не удалось загрузить заказы";
  } finally {
    loading.value = false;
  }
};

onMounted(fetchOrders);
</script>

<style scoped>
.orders-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 80px 0;
}
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}
h1 {
  margin-bottom: 30px;
  color: #333;
}
.order-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}
.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 10px;
  margin-bottom: 15px;
  border-bottom: 1px solid #eee;
}
.order-id {
  font-weight: bold;
  font-size: 1.1rem;
}
.order-date {
  color: #666;
  font-size: 0.9rem;
}
.order-status {
  padding: 4px 8px;
  background: #f0f0f0;
  border-radius: 4px;
  font-size: 0.8rem;
}
.order-items {
  margin-bottom: 15px;
}
.order-item {
  display: flex;
  gap: 15px;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f5f5f5;
}
.item-image {
  width: 60px;
  height: 60px;
  object-fit: contain;
}
.item-info {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.item-name {
  flex: 2;
}
.item-quantity,
.item-price {
  width: 80px;
  text-align: right;
}
.order-total {
  text-align: right;
  font-size: 1.1rem;
  margin-bottom: 15px;
}
.order-details {
  background: #f9f9f9;
  padding: 10px;
  border-radius: 8px;
  font-size: 0.9rem;
}
.loading,
.error,
.empty {
  text-align: center;
  padding: 40px;
  background: white;
  border-radius: 12px;
}
</style>
