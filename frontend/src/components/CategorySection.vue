<template>
  <div class="category-section" :id="id">
    <div class="category-header">
      <div class="category-icon">
        <img :src="category.icon" :alt="category.name" />
      </div>
      <h3>{{ category.name }}</h3>
    </div>

    <div class="products-grid">
      <ProductCard
        v-for="product in products"
        :key="product.id"
        :product="product"
        @add-to-cart="$emit('add-to-cart', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import ProductCard from "./ProductCard.vue";

const props = defineProps({
  id: {
    type: String,
    default: "",
  },
  category: {
    type: Object,
    required: true,
  },
  products: {
    type: Array,
    required: true,
  },
});

defineEmits(["add-to-cart"]);
</script>

<style scoped>
.category-section {
  margin-bottom: 60px;
  background-color: white;
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid #e0e0e0;
  transition: all 0.3s ease;
  scroll-margin-top: 100px;
}

.category-section:hover {
  box-shadow: 0 8px 30px rgba(74, 144, 226, 0.15);
}

.category-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 30px;
  border-bottom: 2px solid #4a90e2;
  padding-bottom: 15px;
}

.category-icon {
  width: 50px;
  height: 50px;
  background-color: #4a90e2;
  border-radius: 12px;
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.category-icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  filter: brightness(0) invert(1);
}

.category-header h3 {
  font-size: 1.8rem;
  margin: 0;
  color: #333;
  font-weight: 600;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 25px;
}

@media (max-width: 768px) {
  .category-section {
    padding: 20px;
  }

  .category-header h3 {
    font-size: 1.5rem;
  }

  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
  }
}

@media (max-width: 480px) {
  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 12px;
  }

  .category-icon {
    width: 40px;
    height: 40px;
    padding: 8px;
  }
}
</style>
