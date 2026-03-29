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
            <p>Здесь будут отображаться ваши заказы.</p>
          </div>

          <div v-else-if="activeTab === 'builds'" class="tab-content">
            <h2>Мои сборки</h2>
            <p>Здесь будут отображаться ваши сборки ПК.</p>
          </div>

          <div v-else-if="activeTab === 'favorites'" class="tab-content">
            <h2>Избранное</h2>
            <p>Здесь будут отображаться избранные товары.</p>
          </div>

          <div v-else-if="activeTab === 'feedback'" class="tab-content">
            <h2>Обратная связь</h2>
            <p>Форма обратной связи будет здесь.</p>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { userCache } from "../utils/cache";

const activeTab = ref("profile");
const user = ref(null);
const loading = ref(true);

const menuItems = [
  { id: "profile", name: "Мой профиль", icon: "fas fa-user" },
  { id: "orders", name: "Мои заказы", icon: "fas fa-shopping-bag" },
  { id: "builds", name: "Мои сборки", icon: "fas fa-tv" },
  { id: "favorites", name: "Избранное", icon: "fas fa-heart" },
  { id: "feedback", name: "Обратная связь", icon: "fas fa-envelope" },
];

const formatDate = (dateString) => {
  if (!dateString) return "";
  const date = new Date(dateString);
  return date.toLocaleDateString("ru-RU", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
};

const fetchProfile = async () => {
  loading.value = true;
  try {
    const data = await userCache.fetch();
    user.value = data;
  } catch (err) {
    console.error("Ошибка загрузки профиля:", err);
    user.value = null;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchProfile();
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
}
</style>
