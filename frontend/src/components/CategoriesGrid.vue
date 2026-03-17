<template>
  <section class="categories-grid-section">
    <div class="container">
      <h2 class="section-title">Категории запчастей</h2>
      <div class="categories-grid">
        <router-link
          v-for="category in categories"
          :key="category.id"
          :to="`/category/${category.id}`"
          class="category-card"
        >
          <div class="category-icon">
            <img :src="category.icon" :alt="category.name" />
          </div>
          <h3>{{ category.name }}</h3>
          <span class="category-count">{{ category.products.length }} товаров</span>
        </router-link>
      </div>
    </div>
  </section>
</template>

<script setup>
defineProps({
  categories: {
    type: Array,
    required: true,
  },
});
</script>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}
.section-title {
  color: #333;
  font-size: 2rem;
  margin-bottom: 40px;
  text-align: center;
  position: relative;
  padding-bottom: 15px;
}
.section-title::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 60px;
  height: 3px;
  background-color: #4a90e2;
  border-radius: 2px;
}
.categories-grid-section {
  margin-bottom: 40px;
}
.categories-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 25px;
  margin-top: 30px;
}
.category-card {
  background: white;
  border-radius: 16px;
  padding: 30px 20px;
  text-align: center;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  border: 1px solid #e0e0e0;
  animation: fadeIn 0.5s ease-out;
  display: block;
}
.category-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 25px rgba(74, 144, 226, 0.2);
  border-color: #4a90e2;
}
.category-icon {
  width: 90px;
  height: 90px;
  margin: 0 auto 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #4a90e2;
  border-radius: 50%;
  padding: 20px;
  transition: all 0.3s ease;
}
.category-card:hover .category-icon {
  transform: scale(1.1);
  background-color: #357abd;
}
.category-icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  filter: brightness(0) invert(1);
}
.category-card h3 {
  color: #333;
  font-size: 1.2rem;
  margin-bottom: 10px;
  font-weight: 600;
}
.category-count {
  color: #666;
  font-size: 0.95rem;
  background: #f0f0f0;
  padding: 5px 12px;
  border-radius: 25px;
  display: inline-block;
}
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
@media (max-width: 1024px) {
  .categories-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;
  }
}
@media (max-width: 768px) {
  .section-title {
    font-size: 1.5rem;
  }
  .categories-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 15px;
  }
  .category-card {
    padding: 20px 15px;
  }
  .category-icon {
    width: 70px;
    height: 70px;
    padding: 15px;
  }
  .category-card h3 {
    font-size: 1rem;
  }
}
@media (max-width: 480px) {
  .categories-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  .category-card {
    padding: 20px 15px;
  }
}
</style>
