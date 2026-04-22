<template>
  <div class="auth-page">
    <div class="container">
      <div class="auth-card">
        <h2>Регистрация</h2>
        <form @submit.prevent="register">
          <div class="form-group">
            <label>Имя пользователя</label>
            <input
              type="text"
              v-model="username"
              @input="validateUsername"
              required
            />
            <small class="hint">Максимум 10 символов</small>
            <span v-if="usernameError" class="error-text">{{
              usernameError
            }}</span>
          </div>
          <div class="form-group">
            <label>Email</label>
            <input
              type="email"
              v-model="email"
              @input="validateEmail"
              required
            />
            <span v-if="emailError" class="error-text">{{ emailError }}</span>
          </div>
          <div class="form-group">
            <label>Пароль</label>
            <input
              type="password"
              v-model="password"
              @input="validatePassword"
              required
            />
            <div class="password-requirements">
              <small :class="{ valid: passwordValid.length }"
                >✓ Минимум 8 символов</small
              >
              <small :class="{ valid: passwordValid.upper }"
                >✓ Заглавная буква</small
              >
              <small :class="{ valid: passwordValid.digit }">✓ Цифра</small>
              <small :class="{ valid: passwordValid.special }"
                >✓ Спецсимвол (!@#$%^&*)</small
              >
            </div>
            <span v-if="passwordError" class="error-text">{{
              passwordError
            }}</span>
          </div>
          <button type="submit" :disabled="loading || !isFormValid">
            Зарегистрироваться
          </button>
          <p v-if="error" class="error">{{ error }}</p>
          <p>Уже есть аккаунт? <router-link to="/login">Войти</router-link></p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { fetchUser, userCache } from "../utils/cache";

const router = useRouter();
const username = ref("");
const email = ref("");
const password = ref("");
const loading = ref(false);
const error = ref("");

const usernameError = ref("");
const emailError = ref("");
const passwordError = ref("");
const passwordValid = ref({
  length: false,
  upper: false,
  digit: false,
  special: false,
});

const validateUsername = () => {
  if (username.value.length > 10) {
    usernameError.value = "Имя пользователя не должно превышать 10 символов";
  } else if (username.value.length === 0) {
    usernameError.value = "";
  } else {
    usernameError.value = "";
  }
};

const validateEmail = () => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (email.value && !emailRegex.test(email.value)) {
    emailError.value = "Введите корректный email";
  } else {
    emailError.value = "";
  }
};

const validatePassword = () => {
  const pwd = password.value;
  passwordValid.value.length = pwd.length >= 8;
  passwordValid.value.upper = /[A-Z]/.test(pwd);
  passwordValid.value.digit = /[0-9]/.test(pwd);
  passwordValid.value.special = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(
    pwd,
  );

  if (pwd && !passwordValid.value.length) {
    passwordError.value = "Пароль должен содержать минимум 8 символов";
  } else if (pwd && !passwordValid.value.upper) {
    passwordError.value = "Пароль должен содержать заглавную букву";
  } else if (pwd && !passwordValid.value.digit) {
    passwordError.value = "Пароль должен содержать цифру";
  } else if (pwd && !passwordValid.value.special) {
    passwordError.value = "Пароль должен содержать специальный символ";
  } else {
    passwordError.value = "";
  }
};

const isFormValid = computed(() => {
  return (
    username.value.trim() !== "" &&
    email.value.trim() !== "" &&
    password.value !== "" &&
    usernameError.value === "" &&
    emailError.value === "" &&
    passwordError.value === "" &&
    passwordValid.value.length &&
    passwordValid.value.upper &&
    passwordValid.value.digit &&
    passwordValid.value.special
  );
});

const register = async () => {
  loading.value = true;
  error.value = "";

  validateUsername();
  validateEmail();
  validatePassword();

  if (!isFormValid.value) {
    error.value = "Пожалуйста, исправьте ошибки в форме";
    loading.value = false;
    return;
  }

  try {
    const response = await axios.post("/api/register", {
      username: username.value,
      email: email.value,
      password: password.value,
    });

    if (response.data.token) {
      localStorage.setItem("session_token", response.data.token);
    }

    userCache.clear();
    await fetchUser();
    router.push("/");
  } catch (err) {
    error.value = err.response?.data || "Ошибка регистрации";
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}
.container {
  max-width: 400px;
  width: 100%;
}
.auth-card {
  background: white;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}
h2 {
  margin-bottom: 24px;
  text-align: center;
}
.form-group {
  margin-bottom: 16px;
}
label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
}
input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
}
button {
  width: 100%;
  padding: 12px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
}
button:disabled {
  opacity: 0.7;
}
.error {
  color: #e74c3c;
  margin-top: 12px;
  text-align: center;
}
.password-requirements {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 5px;
  font-size: 0.75rem;
}
.password-requirements small {
  color: #999;
}
.password-requirements small.valid {
  color: #27ae60;
}
.error-text {
  color: #e74c3c;
  font-size: 0.75rem;
  display: block;
  margin-top: 4px;
}
.auth-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}
.container {
  max-width: 400px;
  width: 100%;
}
.auth-card {
  background: white;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}
h2 {
  margin-bottom: 24px;
  text-align: center;
}
.form-group {
  margin-bottom: 16px;
}
label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
}
input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
}
button {
  width: 100%;
  padding: 12px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
}
button:disabled {
  opacity: 0.7;
}
</style>
