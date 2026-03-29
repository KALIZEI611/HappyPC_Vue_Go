<template>
  <nav class="breadcrumbs">
    <div class="container">
      <ul>
        <li v-for="(item, index) in items" :key="index">
          <router-link :to="item.path" v-if="item.path">
            {{ item.name }}
          </router-link>
          <span v-else>{{ item.name }}</span>
          <span v-if="index < items.length - 1" class="separator">/</span>
        </li>
      </ul>
    </div>
  </nav>
</template>

<script setup>
defineProps({
  items: {
    type: Array,
    required: true,
  },
});
</script>

<style scoped>
.breadcrumbs {
  background: #f8f9fa;
  padding: 12px 0;
  border-bottom: 1px solid #e0e0e0;
  margin-top: 80px; 
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

ul {
  display: flex;
  flex-wrap: wrap;
  list-style: none;
  margin: 0;
  padding: 0;
  font-size: 0.9rem;
}

li {
  display: flex;
  align-items: center;
  color: #666;
}

a {
  color: #4a90e2;
  text-decoration: none;
  transition: opacity 0.2s;
}

a:hover {
  opacity: 0.8;
}

.separator {
  margin: 0 8px;
  color: #999;
}
</style>
