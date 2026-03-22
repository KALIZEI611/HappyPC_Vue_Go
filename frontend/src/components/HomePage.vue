<template>
  <div class="home-page">
    <HeroSection />
    <CategoriesGrid :categories="categories" />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import axios from "axios";
import HeroSection from "./HeroSection.vue";
import CategoriesGrid from "./CategoriesGrid.vue";

const categories = ref([]);

const fetchCategories = async () => {
  try {
    const { data: cats } = await axios.get("/categories");
    if (!Array.isArray(cats)) {
      console.error("Ответ от сервера не является массивом:", cats);
      categories.value = [];
      return;
    }
    const promises = cats.map(async (cat) => {
      const productsRes = await axios.get(`/${cat.url_key}`);
      return {
        id: cat.id,
        name: cat.name,
        icon: cat.icon_url,
        products: productsRes.data,
      };
    });
    categories.value = await Promise.all(promises);
  } catch (err) {
    console.error("Ошибка загрузки категорий:", err);
  }
};

onMounted(fetchCategories);
</script>

<style scoped>
.home-page {
  min-height: 100vh;
}
</style>
