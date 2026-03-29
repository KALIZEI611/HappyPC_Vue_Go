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
import { homeCategoriesCache } from "../utils/cache";
import { useBreadcrumbs } from "../composables/useBreadcrumbs";

const { setBreadcrumbs } = useBreadcrumbs();

const categories = ref([]);
let loading = false;

const fetchCategories = async () => {
  if (loading || categories.value.length > 0) return;

  const cached = homeCategoriesCache.get();
  if (cached) {
    categories.value = cached;
    return;
  }

  loading = true;
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

    const result = await Promise.all(promises);
    categories.value = result;
    homeCategoriesCache.set(result);
  } catch (err) {
    console.error("Ошибка загрузки категорий:", err);
  } finally {
    loading = false;
  }
};

onMounted(() => {
  fetchCategories();
  setBreadcrumbs([{ name: "Главная", path: "/" }]);
});
</script>

<style scoped>
.home-page {
  min-height: 100vh;
}
</style>
 