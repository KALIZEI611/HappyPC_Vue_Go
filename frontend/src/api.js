import axios from 'axios';

// Используем переменную окружения или URL бэкенда по умолчанию
const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://happypcvuego-production.up.railway.app';

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

export default api;