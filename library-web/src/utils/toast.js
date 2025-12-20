import { reactive } from 'vue';

export const toastState = reactive({
  items: []
});

let idCounter = 0;

export const showToast = (message, type = 'success') => {
  const id = idCounter++;
  const item = { id, message, type };
  toastState.items.push(item);

  // 3秒后自动消失
  setTimeout(() => {
    const index = toastState.items.findIndex(i => i.id === id);
    if (index !== -1) toastState.items.splice(index, 1);
  }, 3000);
};

export const $message = {
  success: (msg) => showToast(msg, 'success'),
  error: (msg) => showToast(msg, 'error'),
  warning: (msg) => showToast(msg, 'warning')
};