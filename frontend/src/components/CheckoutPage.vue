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
            <div v-if="deliveryMethod === 'pickup'" class="pickup-map">
              <h3>Пункт самовывоза</h3>
              <p>Старый Оскол, микрорайон Майский, 8</p>
              <iframe
                width="100%"
                height="300"
                frameborder="0"
                style="border: 0; border-radius: 8px"
                src="https://yandex.ru/map-widget/v1/?ll=37.912777%2C51.317990&z=17&pt=37.912777,51.317990,flag"
                allowfullscreen
              >
              </iframe>
            </div>
          </div>

          <div class="form-section">
            <h2>Способ оплаты</h2>
            <div class="payment-options">
              <label class="radio-label">
                <input type="radio" v-model="paymentMethod" value="online" />
                Онлайн (карта)
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

            <div v-if="paymentMethod === 'online'" class="card-form">
              <div class="form-group">
                <label>Номер карты</label>
                <input
                  type="text"
                  v-model="cardNumber"
                  placeholder="0000 0000 0000 0000"
                  maxlength="19"
                  @input="onCardNumberInput"
                />
                <div class="card-icons">
                  <i class="fab fa-cc-visa" :class="{ active: cardType === 'visa' }"></i>
                  <i
                    class="fab fa-cc-mastercard"
                    :class="{ active: cardType === 'mastercard' }"
                  ></i>
                  <i
                    class="fas fa-credit-card"
                    :class="{ active: cardType === 'mir' }"
                  ></i>
                </div>
                <span v-if="cardErrors.number" class="error-text">{{
                  cardErrors.number
                }}</span>
              </div>

              <div class="form-row">
                <div class="form-group half">
                  <label>Срок действия (MM/YY)</label>
                  <input
                    type="text"
                    v-model="cardDate"
                    placeholder="MM/YY"
                    maxlength="5"
                    @input="onCardDateInput"
                  />
                  <span v-if="cardErrors.date" class="error-text">{{
                    cardErrors.date
                  }}</span>
                </div>
                <div class="form-group half">
                  <label>CVV/CVC</label>
                  <input
                    type="password"
                    v-model="cardCvv"
                    placeholder="123"
                    maxlength="3"
                    @input="onCardCvvInput"
                  />
                  <span v-if="cardErrors.cvv" class="error-text">{{
                    cardErrors.cvv
                  }}</span>
                </div>
              </div>
            </div>

            <div v-if="paymentMethod === 'credit'" class="credit-form">
              <div class="form-group">
                <label>ФИО полностью</label>
                <input
                  type="text"
                  v-model="creditFullName"
                  placeholder="Иванов Иван Иванович"
                />
                <span v-if="creditErrors.fullName" class="error-text">{{
                  creditErrors.fullName
                }}</span>
              </div>
              <div class="form-row">
                <div class="form-group half">
                  <label>Серия паспорта</label>
                  <input
                    type="text"
                    v-model="creditPassportSeries"
                    placeholder="1234"
                    maxlength="4"
                    @input="onPassportSeriesInput"
                  />
                  <span v-if="creditErrors.passportSeries" class="error-text">{{
                    creditErrors.passportSeries
                  }}</span>
                </div>
                <div class="form-group half">
                  <label>Номер паспорта</label>
                  <input
                    type="text"
                    v-model="creditPassportNumber"
                    placeholder="123456"
                    maxlength="6"
                    @input="onPassportNumberInput"
                  />
                  <span v-if="creditErrors.passportNumber" class="error-text">{{
                    creditErrors.passportNumber
                  }}</span>
                </div>
              </div>
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

const cardNumber = ref("");
const cardDate = ref("");
const cardCvv = ref("");
const cardErrors = ref({ number: "", date: "", cvv: "" });
const cardType = ref(null);

const creditFullName = ref("");
const creditPassportSeries = ref("");
const creditPassportNumber = ref("");
const creditErrors = ref({ fullName: "", passportSeries: "", passportNumber: "" });

const subtotal = computed(() =>
  props.cart.reduce((sum, item) => sum + item.product.price * item.quantity, 0)
);
const totalPrice = computed(() =>
  deliveryMethod.value === "delivery" ? subtotal.value + 200 : subtotal.value
);

const updateCardType = (num) => {
  const digits = num.replace(/\s/g, "");
  if (!digits) {
    cardType.value = null;
    return;
  }
  if (digits.startsWith("4")) {
    cardType.value = "visa";
  } else if (/^5[1-5]/.test(digits) || /^222[1-9]|^22[3-9]|^2[3-6]/.test(digits)) {
    cardType.value = "mastercard";
  } else if (/^220[0-4]/.test(digits)) {
    cardType.value = "mir";
  } else {
    cardType.value = null;
  }
};

const validateCardNumber = (num) => {
  const digits = num.replace(/\s/g, "");
  if (!digits) return false;
  return /^\d{16}$/.test(digits);
};

const validateCardDate = (date) => {
  const match = date.match(/^(\d{2})\/(\d{2})$/);
  if (!match) return false;
  const month = parseInt(match[1], 10);
  const year = parseInt(match[2], 10);
  if (month < 1 || month > 12) return false;
  const currentYear = new Date().getFullYear() % 100;
  const currentMonth = new Date().getMonth() + 1;
  if (year < currentYear || (year === currentYear && month < currentMonth)) return false;
  return true;
};

const validateCvv = (cvv) => /^\d{3}$/.test(cvv);

const validateCreditFullName = (name) => {
  return name.trim().split(/\s+/).length >= 3;
};

const validatePassportSeries = (series) => {
  return /^\d{4}$/.test(series);
};

const validatePassportNumber = (number) => {
  return /^\d{6}$/.test(number);
};

const onCardNumberInput = (e) => {
  let value = e.target.value.replace(/\D/g, "");
  if (value.length > 16) value = value.slice(0, 16);
  const parts = value.match(/.{1,4}/g);
  cardNumber.value = parts ? parts.join(" ") : value;
  updateCardType(cardNumber.value);

  if (validateCardNumber(cardNumber.value)) {
    cardErrors.value.number = "";
  }
};

const onCardDateInput = (e) => {
  let value = e.target.value.replace(/\D/g, "");
  if (value.length > 4) value = value.slice(0, 4);
  if (value.length >= 3) {
    cardDate.value = `${value.slice(0, 2)}/${value.slice(2)}`;
  } else {
    cardDate.value = value;
  }

  if (validateCardDate(cardDate.value)) {
    cardErrors.value.date = "";
  }
};

const onCardCvvInput = (e) => {
  let value = e.target.value.replace(/\D/g, "");
  if (value.length > 3) value = value.slice(0, 3);
  cardCvv.value = value;

  if (validateCvv(cardCvv.value)) {
    cardErrors.value.cvv = "";
  }
};

const onPassportSeriesInput = (e) => {
  let value = e.target.value.replace(/\D/g, "");
  if (value.length > 4) value = value.slice(0, 4);
  creditPassportSeries.value = value;

  if (validatePassportSeries(creditPassportSeries.value)) {
    creditErrors.value.passportSeries = "";
  }
};

const onPassportNumberInput = (e) => {
  let value = e.target.value.replace(/\D/g, "");
  if (value.length > 6) value = value.slice(0, 6);
  creditPassportNumber.value = value;

  if (validatePassportNumber(creditPassportNumber.value)) {
    creditErrors.value.passportNumber = "";
  }
};

const submitOrder = async () => {
  if (!user.value) {
    router.push("/login");
    return;
  }
  if (deliveryMethod.value === "delivery" && !address.value.trim()) {
    alert("Укажите адрес доставки");
    return;
  }

  let hasError = false;

  if (paymentMethod.value === "online") {
    if (!validateCardNumber(cardNumber.value)) {
      cardErrors.value.number = "Номер карты должен содержать 16 цифр";
      hasError = true;
    }
    if (!validateCardDate(cardDate.value)) {
      cardErrors.value.date = "Неверный формат или срок действия";
      hasError = true;
    }
    if (!validateCvv(cardCvv.value)) {
      cardErrors.value.cvv = "CVV должен содержать 3 цифры";
      hasError = true;
    }
  } else if (paymentMethod.value === "credit") {
    if (!validateCreditFullName(creditFullName.value)) {
      creditErrors.value.fullName = "Введите полное ФИО (Фамилия Имя Отчество)";
      hasError = true;
    }
    if (!validatePassportSeries(creditPassportSeries.value)) {
      creditErrors.value.passportSeries = "Серия паспорта должна содержать 4 цифры";
      hasError = true;
    }
    if (!validatePassportNumber(creditPassportNumber.value)) {
      creditErrors.value.passportNumber = "Номер паспорта должен содержать 6 цифр";
      hasError = true;
    }
  }

  if (hasError) return;

  loading.value = true;
  try {
    const items = props.cart.map((item) => ({
      product_id: item.product.id,
      quantity: item.quantity,
      price: item.product.price,
    }));

    const response = await axios.post("/api/orders", {
      items,
      delivery_method: deliveryMethod.value,
      delivery_address: deliveryMethod.value === "delivery" ? address.value : "",
      payment_method: paymentMethod.value,
      total_price: totalPrice.value,
    });

    const orderId = response.data.id;
    emit("cart-cleared");
    router.push({ path: "/order-success", query: { order_id: orderId } });
  } catch (err) {
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
.credit-form {
  margin-top: 20px;
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
.card-form {
  margin-top: 15px;
}
.form-group {
  margin-bottom: 15px;
}
.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
}
.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
}
.form-row {
  display: flex;
  gap: 15px;
}
.form-group.half {
  flex: 1;
}
.card-icons {
  display: flex;
  gap: 10px;
  margin-top: 5px;
}
.card-icons i {
  font-size: 1.5rem;
  color: #ccc;
}
.card-icons i.active {
  color: #4a90e2;
}
.error-text {
  color: #e74c3c;
  font-size: 0.8rem;
  margin-top: 4px;
  display: block;
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
