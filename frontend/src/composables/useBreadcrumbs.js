import { ref } from 'vue';

const breadcrumbItems = ref([]);

export function useBreadcrumbs() {
  const setBreadcrumbs = (items) => {
    breadcrumbItems.value = items;
  };

  return {
    breadcrumbItems,
    setBreadcrumbs,
  };
}