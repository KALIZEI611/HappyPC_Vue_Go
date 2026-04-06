<template>
  <div class="pc-builder-page">
    <div class="container">
      <h1>Сборка ПК</h1>
      <p class="subtitle">Выберите комплектующие – система проверит совместимость.</p>

      <div v-if="loading" class="loading">Загрузка компонентов...</div>
      <div v-else class="builder-content">
        <div class="components-grid">
          <div v-for="comp in componentTypes" :key="comp.key" class="component-card">
            <h3>{{ comp.label }}</h3>
            <div v-if="selected[comp.key]" class="selected-item">
              <img :src="selected[comp.key].image" :alt="selected[comp.key].name" />
              <div class="selected-info">
                <span class="item-name">{{ selected[comp.key].name }}</span>
                <span class="item-price"
                  >{{ selected[comp.key].price.toLocaleString() }} ₽</span
                >
              </div>
              <button
                @click="removeComponent(comp.key)"
                class="remove-btn"
                title="Удалить"
              >
                <i class="fas fa-trash"></i>
              </button>
            </div>
            <div v-else class="select-placeholder">
              <p>Не выбран</p>
              <button @click="goToCategory(comp.categoryId, comp.key)" class="select-btn">
                Выбрать {{ comp.label.toLowerCase() }}
              </button>
            </div>
          </div>
        </div>

        <div class="compatibility-panel">
          <h3>Проверка совместимости</h3>
          <ul>
            <li v-for="(msg, idx) in compatibilityMessages" :key="idx" :class="msg.type">
              {{ msg.text }}
            </li>
          </ul>
          <div class="total-price">
            <strong>Итоговая цена:</strong> {{ totalPrice.toLocaleString() }} ₽
          </div>

          <div class="button-group">
            <button
              @click="addToCart"
              :disabled="!isCompatible || totalPrice === 0"
              class="add-to-cart-btn"
            >
              <i class="fas fa-shopping-cart"></i> Добавить сборку в корзину
            </button>

            <button @click="clearAll" class="clear-all-btn">
              <i class="fas fa-trash-alt"></i> Очистить всё
            </button>
          </div>

          <div class="save-build-section">
            <button
              v-if="!showSaveForm"
              @click="showSaveForm = true"
              class="save-build-btn"
            >
              <i class="fas fa-save"></i> Сохранить сборку
            </button>
            <div v-else class="save-form">
              <input
                type="text"
                v-model="buildName"
                placeholder="Название сборки"
                maxlength="50"
                @keyup.enter="saveBuild"
              />
              <button @click="saveBuild" class="confirm-btn" title="Сохранить">
                <i class="fas fa-check"></i>
              </button>
              <button
                @click="
                  showSaveForm = false;
                  buildName = '';
                "
                class="cancel-btn"
                title="Отмена"
              >
                <i class="fas fa-times"></i>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import { allProductsCache } from "../utils/cache";
import { user } from "../utils/cache";

const router = useRouter();
const route = useRoute();
const loading = ref(true);
const allProducts = ref([]);
const showSaveForm = ref(false);
const buildName = ref("");

const loadSavedBuild = () => {
  const saved = localStorage.getItem("pcBuilderBuild");
  if (saved) {
    try {
      return JSON.parse(saved);
    } catch {
      return {};
    }
  }
  return {};
};

const selected = ref(loadSavedBuild());

const componentTypes = ref([
  { key: "cpu", label: "Процессор", categoryId: 1, items: [] },
  { key: "gpu", label: "Видеокарта", categoryId: 2, items: [] },
  { key: "motherboard", label: "Материнская плата", categoryId: 5, items: [] },
  { key: "ram", label: "Оперативная память", categoryId: 3, items: [] },
  { key: "psu", label: "Блок питания", categoryId: 6, items: [] },
  { key: "cooler", label: "Кулер", categoryId: 9, items: [] },
  { key: "storage", label: "SSD", categoryId: 4, items: [] },
  { key: "case", label: "Корпус", categoryId: 7, items: [] },
]);

const compatibilityMessages = ref([]);
const isCompatible = ref(false);

const totalPrice = computed(() => {
  let sum = 0;
  for (const key in selected.value) {
    if (selected.value[key]) sum += selected.value[key].price;
  }
  return sum;
});

const saveBuild = async () => {
  if (!user.value) {
    router.push("/login");
    return;
  }
  if (!buildName.value.trim()) {
    alert("Введите название сборки");
    return;
  }
  if (totalPrice.value === 0) {
    alert("Сборка пуста. Добавьте хотя бы один компонент.");
    return;
  }

  try {
    const componentsToSave = {};
    for (const key in selected.value) {
      if (selected.value[key]) {
        componentsToSave[key] = {
          id: selected.value[key].id,
          name: selected.value[key].name,
          price: selected.value[key].price,
        };
      }
    }

    await axios.post("/api/builds", {
      name: buildName.value,
      components: componentsToSave,
    });

    alert("Сборка сохранена!");

    showSaveForm.value = false;
    buildName.value = "";

    selected.value = {
      cpu: null,
      gpu: null,
      motherboard: null,
      ram: null,
      psu: null,
      cooler: null,
      storage: null,
      case: null,
    };

    localStorage.removeItem("pcBuilderBuild");

    updateCompatibility();
  } catch (err) {
    console.error("Ошибка сохранения сборки:", err);
    alert("Не удалось сохранить сборку");
  }
};

const clearAll = () => {
  if (totalPrice.value === 0) {
    alert("Сборка уже пуста");
    return;
  }

  if (confirm("Вы уверены, что хотите очистить всю сборку?")) {
    selected.value = {
      cpu: null,
      gpu: null,
      motherboard: null,
      ram: null,
      psu: null,
      cooler: null,
      storage: null,
      case: null,
    };

    localStorage.removeItem("pcBuilderBuild");
    updateCompatibility();

    alert("Сборка очищена");
  }
};

const saveBuildToLocalStorage = () => {
  localStorage.setItem("pcBuilderBuild", JSON.stringify(selected.value));
};

const removeComponent = (key) => {
  selected.value[key] = null;
  saveBuildToLocalStorage();
  updateCompatibility();
};

const goToCategory = (categoryId, componentKey) => {
  sessionStorage.setItem("pcBuilderReturnKey", componentKey);
  router.push(`/category/${categoryId}?from=builder`);
};

const updateCompatibility = () => {
  const msgs = [];
  let compatible = true;

  const cpu = selected.value.cpu;
  const mb = selected.value.motherboard;
  const ram = selected.value.ram;
  const psu = selected.value.psu;
  const gpu = selected.value.gpu;
  const pcCase = selected.value.case;

  if (cpu && mb) {
    if (cpu.socket !== mb.socket) {
      msgs.push({
        type: "error",
        text: `Сокет процессора (${cpu.socket}) не совместим с материнской платой (${mb.socket})`,
      });
      compatible = false;
    } else {
      msgs.push({ type: "success", text: `Сокет совместим: ${cpu.socket}` });
    }
  }

  if (ram && mb) {
    if (ram.ram_type !== mb.ram_type) {
      msgs.push({
        type: "error",
        text: `Тип памяти RAM (${ram.ram_type}) не поддерживается материнской платой (${mb.ram_type})`,
      });
      compatible = false;
    } else {
      msgs.push({
        type: "success",
        text: `Тип памяти совместим: ${ram.ram_type}`,
      });
    }
  }

  let totalPower = (cpu ? cpu.tdp : 0) + (gpu ? gpu.power : 0);
  if (psu) {
    if (totalPower > psu.power) {
      msgs.push({
        type: "error",
        text: `Недостаточная мощность БП: требуется ≈${totalPower} Вт, выбран ${psu.power} Вт`,
      });
      compatible = false;
    } else if (totalPower > 0) {
      msgs.push({
        type: "success",
        text: `Блок питания обеспечит достаточную мощность (запас ${
          psu.power - totalPower
        } Вт)`,
      });
    }
  } else if (totalPower > 0) {
    msgs.push({ type: "warning", text: "Выберите блок питания" });
  }

  if (mb && pcCase) {
    const mbForm = mb.form_factor?.toLowerCase();
    const caseForm = pcCase.form_factor?.toLowerCase();
    if (mbForm && caseForm) {
      if (
        (mbForm === "atx" && caseForm === "atx") ||
        (mbForm === "micro-atx" && (caseForm === "micro-atx" || caseForm === "atx")) ||
        (mbForm === "mini-itx" &&
          (caseForm === "mini-itx" || caseForm === "micro-atx" || caseForm === "atx"))
      ) {
        msgs.push({
          type: "success",
          text: `Материнская плата (${mb.form_factor}) помещается в корпус (${pcCase.form_factor})`,
        });
      } else {
        msgs.push({
          type: "error",
          text: `Форм-фактор материнской платы (${mb.form_factor}) несовместим с корпусом (${pcCase.form_factor})`,
        });
        compatible = false;
      }
    }
  }

  compatibilityMessages.value = msgs;
  isCompatible.value = compatible;
  saveBuildToLocalStorage();
};

const addToCart = async () => {
  if (!isCompatible.value || totalPrice.value === 0) return;

  try {
    for (const key in selected.value) {
      if (selected.value[key]) {
        await axios.post("/api/cart", {
          product_id: selected.value[key].id,
          quantity: 1,
        });
      }
    }
    selected.value = {
      cpu: null,
      gpu: null,
      motherboard: null,
      ram: null,
      psu: null,
      cooler: null,
      storage: null,
      case: null,
    };
    localStorage.removeItem("pcBuilderBuild");
    updateCompatibility();
    router.push("/cart");
  } catch (err) {
    console.error("Ошибка добавления в корзину:", err);
    alert("Не удалось добавить сборку в корзину");
  }
};

onMounted(async () => {
  const returnKey = sessionStorage.getItem("pcBuilderReturnKey");
  const selectedProductId = route.query.product_id;

  if (returnKey && selectedProductId) {
    let products = allProductsCache.get();
    if (!products) {
      const { data: cats } = await axios.get("/categories");
      const promises = cats.map(async (cat) => {
        const res = await axios.get(`/${cat.url_key}`);
        return res.data;
      });
      const results = await Promise.all(promises);
      products = results.flat();
      allProductsCache.set(products);
    }
    const product = products.find((p) => p.id === parseInt(selectedProductId));
    if (product) {
      selected.value[returnKey] = product;
      sessionStorage.removeItem("pcBuilderReturnKey");
      router.replace("/pc-builder");
      updateCompatibility();
    }
  }

  const savedBuild = loadSavedBuild();
  if (Object.keys(savedBuild).length > 0) {
    selected.value = savedBuild;
    updateCompatibility();
  }

  loading.value = false;
  updateCompatibility();
});
</script>

<style scoped>
.button-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 15px;
}

.clear-all-btn {
  width: 100%;
  padding: 12px;
  background-color: #e74c3c;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: background-color 0.3s;
}

.clear-all-btn:hover {
  background-color: #c0392b;
}

.add-to-cart-btn {
  width: 100%;
  padding: 12px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: background-color 0.3s;
}

.add-to-cart-btn:hover:not(:disabled) {
  background-color: #357abd;
}

.add-to-cart-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
.save-build-section {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #e0e0e0;
}

.save-build-btn {
  width: 100%;
  padding: 12px;
  background-color: #27ae60;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.save-build-btn:hover {
  background-color: #2ecc71;
}

.save-form {
  display: flex;
  gap: 8px;
  align-items: center;
}

.save-form input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 0.9rem;
}

.confirm-btn,
.cancel-btn {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
}

.confirm-btn {
  background-color: #27ae60;
  color: white;
}

.confirm-btn:hover {
  background-color: #2ecc71;
}

.cancel-btn {
  background-color: #e74c3c;
  color: white;
}

.cancel-btn:hover {
  background-color: #c0392b;
}
.pc-builder-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 80px 0;
}
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}
h1 {
  margin-bottom: 10px;
}
.subtitle {
  color: #666;
  margin-bottom: 30px;
}
.loading {
  text-align: center;
  padding: 40px;
}
.builder-content {
  display: flex;
  gap: 30px;
  flex-wrap: wrap;
}
.components-grid {
  flex: 2;
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
}
.component-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
.component-card h3 {
  margin-bottom: 12px;
  font-size: 1.2rem;
}
.selected-item {
  display: flex;
  align-items: center;
  gap: 15px;
}

.selected-item img {
  width: 60px;
  height: 60px;
  object-fit: contain;
}

.selected-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item-name {
  font-weight: 500;
  font-size: 0.95rem;
}

.item-price {
  color: #4a90e2;
  font-weight: bold;
  font-size: 0.9rem;
}
.item-name {
  display: block;
  font-weight: 500;
  margin-bottom: 5px;
}
.item-price {
  color: #4a90e2;
  font-weight: bold;
}
.remove-btn {
  background: none;
  border: none;
  color: #e74c3c;
  cursor: pointer;
  font-size: 1.1rem;
  padding: 8px;
  border-radius: 50%;
  transition: background 0.2s;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.remove-btn:hover {
  background: #fee;
  transform: scale(1.05);
}
.select-placeholder {
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}
.select-placeholder p {
  margin-bottom: 10px;
  color: #666;
}
.select-btn {
  padding: 8px 16px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}
.select-btn:hover {
  background-color: #357abd;
}
.compatibility-panel {
  flex: 1;
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 100px;
  height: fit-content;
}
.compatibility-panel h3 {
  margin-bottom: 15px;
}
.compatibility-panel ul {
  list-style: none;
  padding: 0;
  margin-bottom: 20px;
}
.compatibility-panel li {
  padding: 6px 0;
  border-bottom: 1px solid #eee;
  font-size: 0.9rem;
}
.compatibility-panel li.success {
  color: #27ae60;
}
.compatibility-panel li.error {
  color: #e74c3c;
}
.compatibility-panel li.warning {
  color: #f39c12;
}
.total-price {
  font-size: 1.3rem;
  margin: 15px 0;
}
.add-to-cart-btn {
  width: 100%;
  padding: 12px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
}
.add-to-cart-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
@media (max-width: 768px) {
  .builder-content {
    flex-direction: column;
  }
  .selected-item {
    flex-direction: column;
    text-align: center;
  }
}
</style>
