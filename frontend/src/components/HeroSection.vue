<template>
  <section class="hero">
    <div class="container">
      <h1>Комплектующие для вашего ПК</h1>
      <p>Более 1000 товаров от ведущих производителей</p>
    </div>
  </section>
</template>

<script setup></script>

<style scoped>
.hero {
  background-color: #4a90e2;
  color: white;
  padding: 80px 0;
  text-align: center;
  margin-bottom: 40px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.hero h1 {
  font-size: 3rem;
  margin-bottom: 20px;
  animation: fadeInUp 0.8s ease;
}

.hero p {
  font-size: 1.2rem;
  opacity: 0.95;
  animation: fadeInUp 1s ease;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .hero h1 {
    font-size: 2rem;
  }

  .hero p {
    font-size: 1rem;
  }
}

@media (max-width: 480px) {
  .hero h1 {
    font-size: 1.5rem;
  }
}
</style>
