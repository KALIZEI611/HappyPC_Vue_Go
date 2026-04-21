import axios from "axios";
import api from "../api";

const favoritesService = {
  async getFavorites() {
    const { data } = await api.get("/api/favorites");
    return data;
  },
  
  async addToFavorites(productId) {
    await axios.post("/api/favorites", { product_id: productId });
  },
  
  async removeFromFavorites(productId) {
    await axios.delete(`/api/favorites/${productId}`);
  },
  
  async toggleFavorite(productId) {
    const { data } = await axios.post("/api/favorites/toggle", { product_id: productId });
    return data.isFavorite;
  }
};

export default favoritesService;