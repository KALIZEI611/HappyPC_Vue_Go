<template>
  <div class="auth-page">
    <div class="container">
      <div class="auth-card">
        <h2>Вход</h2>
        <form @submit.prevent="login">
          <div class="form-group">
            <label>Email</label>
            <input type="email" v-model="email" required />
          </div>
          <div class="form-group">
            <label>Пароль</label>
            <input type="password" v-model="password" required />
          </div>
          <button type="submit" :disabled="loading">Войти</button>
          <p v-if="error" class="error">{{ error }}</p>
          <p>
            Нет аккаунта?
            <router-link to="/register">Зарегистрироваться</router-link>
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { fetchUser, userCache } from "../utils/cache";

const router = useRouter();
const email = ref("");
const password = ref("");
const loading = ref(false);
const error = ref("");

const login = async () => {
  loading.value = true;
  error.value = "";
  try {
    const response = await axios.post("/api/login", {
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
    error.value = err.response?.data || "Ошибка входа";
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
</style>
