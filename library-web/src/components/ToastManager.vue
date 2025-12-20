<template>
  <div class="toast-container">
    <transition-group name="toast">
      <div 
        v-for="item in toastState.items" 
        :key="item.id" 
        class="toast-item"
        :class="`toast-${item.type}`"
      >
        <span class="icon" v-if="item.type === 'success'">✅</span>
        <span class="icon" v-if="item.type === 'error'">❌</span>
        <span class="icon" v-if="item.type === 'warning'">⚠️</span>
        {{ item.message }}
      </div>
    </transition-group>
  </div>
</template>

<script setup>
import { toastState } from '../utils/toast';
</script>

<style scoped>
.toast-container {
  position: fixed; top: 20px; left: 50%; transform: translateX(-50%);
  z-index: 9999; display: flex; flex-direction: column; gap: 10px;
  pointer-events: none; /* 让鼠标能透过提示点击下面的内容 */
}
.toast-item {
  padding: 10px 20px; border-radius: 4px; background: white;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15); font-size: 14px;
  display: flex; align-items: center; gap: 8px; font-weight: 500;
  min-width: 300px; justify-content: center; pointer-events: auto;
}
.toast-success { background: #f0fdf4; color: #166534; border: 1px solid #bbf7d0; }
.toast-error { background: #fef2f2; color: #991b1b; border: 1px solid #fecaca; }
.toast-warning { background: #fffbeb; color: #92400e; border: 1px solid #fde68a; }

/* 动画效果 */
.toast-enter-active, .toast-leave-active { transition: all 0.3s ease; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(-20px); }
</style>