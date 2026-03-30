<template>
  <div class="order-success">
    <div class="container">
      <i class="fas fa-check-circle"></i>
      <h1>Заказ оформлен!</h1>
      <p>Спасибо за покупку. Номер вашего заказа: {{ orderId }}</p>
      <router-link to="/" class="btn">Вернуться на главную</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
const orderId = ref(new URLSearchParams(window.location.search).get("order_id") || "");
</script>

<style scoped>
.order-success {
  min-height: 100vh;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
}
.container {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
.fa-check-circle {
  font-size: 4rem;
  color: #27ae60;
  margin-bottom: 20px;
}
.btn {
  display: inline-block;
  margin-top: 20px;
  padding: 10px 20px;
  background: #4a90e2;
  color: white;
  text-decoration: none;
  border-radius: 6px;
}
</style>
