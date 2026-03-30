import axios from "axios";
import { ref } from "vue";

// Реактивная переменная пользователя
export const user = ref(null);
export let userPromise = null;

export const fetchUser = async () => {
  if (user.value !== null) return user.value;
  if (userPromise) return userPromise;
  userPromise = axios
    .get("/api/me")
    .then((res) => {
      user.value = res.data;
      userPromise = null;
      return user.value;
    })
    .catch((err) => {
      user.value = null;
      userPromise = null;
      throw err;
    });
  return userPromise;
};

export const clearUser = () => {
  user.value = null;
  userPromise = null;
};

const storage = {
  get(key) {
    const item = localStorage.getItem(key);
    if (!item) return null;
    try {
      return JSON.parse(item);
    } catch {
      return null;
    }
  },
  set(key, value) {
    localStorage.setItem(key, JSON.stringify(value));
  },
  remove(key) {
    localStorage.removeItem(key);
  },
  clear() {
    localStorage.clear();
  },
};

export const homeCategoriesCache = {
  get() {
    return storage.get("homeCategories");
  },
  set(value) {
    storage.set("homeCategories", value);
  },
  clear() {
    storage.remove("homeCategories");
  },
};

export const categoryProductsCache = {
  get(id) {
    const data = storage.get("categoryProducts");
    return data ? data[id] : null;
  },
  set(id, value) {
    const data = storage.get("categoryProducts") || {};
    data[id] = value;
    storage.set("categoryProducts", data);
  },
  clear() {
    storage.remove("categoryProducts");
  },
};

export const productCache = {
  get(id) {
    const data = storage.get("product");
    return data ? data[id] : null;
  },
  set(id, value) {
    const data = storage.get("product") || {};
    data[id] = value;
    storage.set("product", data);
  },
  clear() {
    storage.remove("product");
  },
};

export const allProductsCache = {
  get() {
    return storage.get("allProducts");
  },
  set(value) {
    storage.set("allProducts", value);
  },
  clear() {
    storage.remove("allProducts");
  },
};

export const userCache = {
  data: null,
  promise: null,
  get() {
    if (this.data) return this.data;
    const stored = storage.get("user");
    if (stored) {
      this.data = stored;
      return stored;
    }
    return null;
  },
  set(value) {
    this.data = value;
    storage.set("user", value);
  },
  clear() {
    this.data = null;
    this.promise = null;
    storage.remove("user");
  },
  async fetch() {
    if (this.data) return this.data;
    if (this.promise) return this.promise;
    this.promise = axios
      .get("/api/me")
      .then((res) => {
        this.data = res.data;
        storage.set("user", this.data);
        this.promise = null;
        return this.data;
      })
      .catch((err) => {
        this.promise = null;
        throw err;
      });
    return this.promise;
  },
};

export const categoriesCache = {
  data: null,
  set(value) {
    this.data = value;
  },
  get() {
    return this.data;
  },
  clear() {
    this.data = null;
  },
};

export function clearAllCaches() {
  homeCategoriesCache.clear();
  categoryProductsCache.clear();
  productCache.clear();
  allProductsCache.clear();
  userCache.clear();
  storage.clear();
}
