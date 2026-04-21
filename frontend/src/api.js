import axios from 'axios';

// Берем URL из переменной окружения или используем URL бэкенда
const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://happypcvuego-production.up.railway.app';

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

export default api;