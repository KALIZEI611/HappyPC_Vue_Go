<template>
  <div class="checkout-page">
    <div class="container">
      <h1>Оформление заказа</h1>

      <div v-if="cart.length === 0" class="empty-cart">
        <i class="fas fa-shopping-cart"></i>
        <p>Корзина пуста</p>
        <router-link to="/" class="continue-shopping">Продолжить покупки</router-link>
      </div>

      <div v-else class="checkout-content">
        <form @submit.prevent="submitOrder">
          <div class="form-section">
            <h2>Данные пользователя</h2>
            <div class="user-info">
              <p><strong>Email:</strong> {{ user.email }}</p>
            </div>
          </div>

          <div class="form-section">
            <h2>Доставка</h2>
            <div class="delivery-options">
              <label class="radio-label">
                <input type="radio" v-model="deliveryMethod" value="pickup" />
                Самовывоз (бесплатно)
              </label>
              <label class="radio-label">
                <input type="radio" v-model="deliveryMethod" value="delivery" />
                Доставка (200 ₽)
              </label>
            </div>
            <div v-if="deliveryMethod === 'delivery'" class="address-input">
              <label for="address">Адрес доставки</label>
              <input type="text" id="address" v-model="address" required />
            </div>
          </div>

          <div class="form-section">
            <h2>Способ оплаты</h2>
            <div class="payment-options">
              <label class="radio-label">
                <input type="radio" v-model="paymentMethod" value="online" />
                Онлайн (карта, СБП)
              </label>
              <label class="radio-label">
                <input type="radio" v-model="paymentMethod" value="cash" />
                При получении
              </label>
              <label class="radio-label">
                <input type="radio" v-model="paymentMethod" value="credit" />
                В кредит
              </label>
            </div>
          </div>

          <div class="order-summary">
            <div class="total-row">
              <span>Товаров на сумму:</span>
              <span>{{ subtotal.toLocaleString() }} ₽</span>
            </div>
            <div class="total-row" v-if="deliveryMethod === 'delivery'">
              <span>Доставка:</span>
              <span>200 ₽</span>
            </div>
            <div class="total-row total">
              <span>Итого:</span>
              <span>{{ totalPrice.toLocaleString() }} ₽</span>
            </div>
            <button type="submit" :disabled="loading" class="submit-btn">
              {{ loading ? "Оформление..." : "Оформить заказ" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { user } from "../utils/cache";

const props = defineProps({
  cart: { type: Array, required: true },
});
const emit = defineEmits(["cart-cleared"]);

const router = useRouter();
const loading = ref(false);

const deliveryMethod = ref("pickup");
const address = ref("");
const paymentMethod = ref("online");

const subtotal = computed(() =>
  props.cart.reduce((sum, item) => sum + item.product.price * item.quantity, 0)
);
const totalPrice = computed(() =>
  deliveryMethod.value === "delivery" ? subtotal.value + 200 : subtotal.value
);

const submitOrder = async () => {
  if (!user.value) {
    router.push("/login");
    return;
  }
  if (deliveryMethod.value === "delivery" && !address.value.trim()) {
    alert("Укажите адрес доставки");
    return;
  }

  loading.value = true;
  try {
    const items = props.cart.map((item) => ({
      product_id: item.product.id,
      quantity: item.quantity,
      price: item.product.price,
    }));

    console.log("Sending order data:", {
      items,
      delivery_method: deliveryMethod.value,
      delivery_address: deliveryMethod.value === "delivery" ? address.value : "",
      payment_method: paymentMethod.value,
      total_price: totalPrice.value,
    });

    const response = await axios.post("/api/orders", {
      items,
      delivery_method: deliveryMethod.value,
      delivery_address: deliveryMethod.value === "delivery" ? address.value : "",
      payment_method: paymentMethod.value,
      total_price: totalPrice.value,
    });

    const orderId = response.data.id;
    emit("cart-cleared"); // сигнал App.vue обновить корзину
    router.push({ path: "/order-success", query: { order_id: orderId } });
  } catch (err) {
    console.error("Order error:", err);
    const message = err.response?.data || "Ошибка оформления заказа";
    alert(message);
  } finally {
    loading.value = false;
  }
};
</script>


<style scoped>
.checkout-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 80px 0;
}
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 20px;
}
h1 {
  margin-bottom: 30px;
  color: #333;
}
.empty-cart {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
}
.form-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}
.form-section h2 {
  font-size: 1.3rem;
  margin-bottom: 16px;
  color: #333;
}
.user-info p {
  margin: 8px 0;
  font-size: 1rem;
}
.delivery-options,
.payment-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.radio-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}
.address-input {
  margin-top: 16px;
}
.address-input label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
}
.address-input input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 8px;
}
.order-summary {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}
.total-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 1rem;
}
.total-row.total {
  font-size: 1.2rem;
  font-weight: bold;
  border-top: 1px solid #e0e0e0;
  padding-top: 12px;
  margin-top: 12px;
}
.submit-btn {
  width: 100%;
  padding: 14px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  margin-top: 20px;
}
.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}
.continue-shopping {
  display: inline-block;
  padding: 12px 30px;
  background-color: #4a90e2;
  color: white;
  text-decoration: none;
  border-radius: 8px;
}
</style>
