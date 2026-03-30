import axios from "axios";

export const cartService = {
  async getCart() {
    const { data } = await axios.get("/api/cart");
    return data;
  },
  async addToCart(productId, quantity = 1) {
    await axios.post("/api/cart", { product_id: productId, quantity });
  },
  async updateCartItem(productId, quantity) {
    await axios.put(`/api/cart/${productId}`, { quantity });
  },
  async removeCartItem(productId) {
    await axios.delete(`/api/cart/${productId}`);
  },
};
