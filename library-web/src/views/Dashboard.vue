<template>
  <div>
    <h2 class="mb-4">📊 运营数据概览</h2>
    
    <div v-if="stats" class="grid-4">
      <!-- 统计卡片 -->
      <div class="stat-card">
        <div class="label">藏书总量</div>
        <div class="value">{{ stats.total_books }}</div>
        <div class="trend text-secondary">库存充沛</div>
      </div>
      
      <div class="stat-card">
        <div class="label">累计借阅</div>
        <div class="value text-primary">{{ stats.total_borrow_count }}</div>
        <div class="trend text-secondary">次</div>
      </div>
      
      <div class="stat-card">
        <div class="label">当前借出</div>
        <div class="value">{{ stats.borrowed_books }}</div>
        <div class="trend text-secondary">本正在流转</div>
      </div>
      
      <div class="stat-card border-danger">
        <div class="label">逾期未还</div>
        <div class="value text-danger">{{ stats.overdue_books }}</div>
        <div class="trend text-danger">需催还</div>
      </div>
    </div>

    <!-- 简单图表模拟 -->
    <div class="card mt-4">
      <h3>📈 借阅趋势 (近30天)</h3>
      <div class="chart-placeholder">
        <div class="bar" style="height: 40%"></div>
        <div class="bar" style="height: 60%"></div>
        <div class="bar" style="height: 45%"></div>
        <div class="bar" style="height: 80%"></div>
        <div class="bar" style="height: 70%"></div>
        <div class="bar active" style="height: 90%"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getStatsOverview } from '../api';

const stats = ref(null);

onMounted(async () => {
  try { stats.value = await getStatsOverview(); } catch (e) {}
});
</script>

<style scoped>
.mb-4 { margin-bottom: 24px; }
.grid-4 { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 24px; }

.stat-card {
  background: white; padding: 24px; border-radius: var(--radius);
  box-shadow: var(--shadow); border: 1px solid transparent;
}
.stat-card.border-danger { border-color: #fee2e2; background: #fef2f2; }

.label { font-size: 0.875rem; color: var(--text-secondary); margin-bottom: 8px; }
.value { font-size: 2rem; font-weight: 700; color: var(--text-main); line-height: 1; margin-bottom: 8px; }
.trend { font-size: 0.875rem; }

.text-primary { color: var(--primary); }
.text-danger { color: var(--danger); }
.text-secondary { color: var(--text-secondary); }

.chart-placeholder {
  height: 200px; display: flex; align-items: flex-end; justify-content: space-around;
  padding-top: 20px; border-bottom: 1px solid var(--border);
}
.bar { width: 40px; background-color: #e5e7eb; border-radius: 4px 4px 0 0; }
.bar.active { background-color: var(--primary); }
</style>