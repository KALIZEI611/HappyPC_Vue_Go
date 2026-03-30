import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "axios";
import "./assets/styles.css";
import { user } from "./utils/cache";

axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || "";
axios.defaults.withCredentials = true;

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const currentPath = window.location.pathname;
      const requestUrl = error.config.url;
      if (
        requestUrl === "/api/me" ||
        currentPath === "/login" ||
        currentPath === "/register"
      ) {
        return Promise.reject(error);
      }
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);

const app = createApp(App);
app.use(router);
app.mount("#app");

setInterval(
  () => {
    if (user.value) {
      axios.get("/api/me").catch(() => {});
    }
  },
  14 * 60 * 1000,
);
