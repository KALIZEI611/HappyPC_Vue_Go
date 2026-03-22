import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
      '/categories': 'http://localhost:8080',
      '/Processors': 'http://localhost:8080',
      '/Video_cards': 'http://localhost:8080',
      '/RAM': 'http://localhost:8080',
      '/SSD': 'http://localhost:8080',
      '/Motherboards': 'http://localhost:8080',
      '/Power_supplies': 'http://localhost:8080',
      '/Buildings': 'http://localhost:8080',
      '/Monitors': 'http://localhost:8080',
      '/Coolers': 'http://localhost:8080',
      '/Keyboards': 'http://localhost:8080',
      '/Mice': 'http://localhost:8080',
      '/Headphones': 'http://localhost:8080',
      '/cart': 'http://localhost:8080',
    }
  }
});