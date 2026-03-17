import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "axios";
import "./assets/styles.css";

axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

const app = createApp(App);
app.use(router);
app.mount("#app");