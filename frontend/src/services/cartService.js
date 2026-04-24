
import api from "../api";

export const cartService = {
  async getCart() {
    const { data } = await api.get("/api/cart");
    return data;
  },
  async addToCart(productId, quantity = 1) {
    await api.post("/api/cart", { product_id: productId, quantity });
  },
  async updateCartItem(productId, quantity) {
    await api.put(`/api/cart/${productId}`, { quantity });
  },
  async removeCartItem(productId) {
    await api.delete(`/api/cart/${productId}`);
  },
  async clearCart() {
    await api.delete("/api/cart");
  },
};

export default cartService;