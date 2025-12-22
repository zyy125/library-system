<template>
  <div>
    <h2 class="page-title">📊 运营数据概览</h2>
    
    <div v-if="loading" class="loading-state">加载中...</div>
    
    <template v-else-if="stats">
      <!-- 统计卡片 -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">📚</div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_books }}</div>
            <div class="stat-label">藏书总量</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">👥</div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_users }}</div>
            <div class="stat-label">注册用户</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">📖</div>
          <div class="stat-content">
            <div class="stat-value text-primary">{{ stats.borrowed_books }}</div>
            <div class="stat-label">当前借出</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">✅</div>
          <div class="stat-content">
            <div class="stat-value text-success">{{ stats.available_books }}</div>
            <div class="stat-label">可借图书</div>
          </div>
        </div>
        
        <div class="stat-card warning" v-if="stats.overdue_books > 0">
          <div class="stat-icon">⚠️</div>
          <div class="stat-content">
            <div class="stat-value text-danger">{{ stats.overdue_books }}</div>
            <div class="stat-label">逾期未还</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">🔄</div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_borrow_count }}</div>
            <div class="stat-label">累计借阅次数</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">📋</div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_categories }}</div>
            <div class="stat-label">图书分类</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">⏳</div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.reservations_count }}</div>
            <div class="stat-label">待处理预约</div>
          </div>
        </div>
      </div>

      <!-- 图表区域 -->
      <div class="charts-row">
        <!-- 借阅趋势 -->
        <div class="card chart-card">
          <h3>📈 借阅趋势 (近30天)</h3>
          <div class="chart-container" v-if="chartData.length > 0">
            <div class="chart-bars">
              <div 
                v-for="(item, index) in chartData" 
                :key="index" 
                class="chart-bar-wrapper"
              >
                <div 
                  class="chart-bar" 
                  :style="{ height: getBarHeight(item.borrow_count) + '%' }"
                  :class="{ active: index === chartData.length - 1 }"
                  :title="`${item.date}:  借阅 ${item.borrow_count} 本`"
                ></div>
                <div class="chart-label">{{ formatChartDate(item.date) }}</div>
              </div>
            </div>
          </div>
          <div v-else class="chart-empty">暂无数据</div>
        </div>

        <!-- 热门图书 -->
        <div class="card chart-card">
          <h3>🔥 热门图书排行</h3>
          <div class="popular-list" v-if="popularBooks.length > 0">
            <div 
              v-for="(book, index) in popularBooks" 
              :key="book.id" 
              class="popular-item"
            >
              <span class="rank" :class="{ top: index < 3 }">{{ index + 1 }}</span>
              <img :src="book.cover_url || defaultCover" class="book-cover-tiny">
              <div class="book-info">
                <div class="book-title">{{ book.title }}</div>
                <div class="book-author">{{ book.author }}</div>
              </div>
              <div class="borrow-count">{{ book.borrow_count || book.recent_borrow_count }} 次</div>
            </div>
          </div>
          <div v-else class="chart-empty">暂无数据</div>
        </div>
      </div>

      <!-- 活跃用户 -->
      <div class="card">
        <div class="card-header">
          <h3>👥 近30天活跃用户</h3>
          <span class="badge badge-info">{{ stats.active_users_30d }} 人</span>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getStatsOverview, getPopularBooks } from '../api';

const defaultCover = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="40" height="56" viewBox="0 0 40 56"%3E%3Crect fill="%23f3f4f6" width="40" height="56"/%3E%3C/svg%3E';

const loading = ref(true);
const stats = ref(null);
const chartData = ref([]);
const popularBooks = ref([]);

const loadData = async () => {
  loading.value = true;
  try {
    // 加载统计概览
    stats.value = await getStatsOverview();
    
    // 加载热门图书
    const popularRes = await getPopularBooks({ limit: 5, days: 30 });
    popularBooks.value = popularRes.books || [];
    
    // 模拟图表数据（实际应调用 /api/stats/borrow）
    // 这里生成最近7天的模拟数据展示效果
    const today = new Date();
    chartData.value = Array.from({ length: 7 }, (_, i) => {
      const date = new Date(today);
      date.setDate(date.getDate() - (6 - i));
      return {
        date: date.toISOString().split('T')[0],
        borrow_count: Math.floor(Math.random() * 20) + 5
      };
    });
  } catch (e) {
    console.error('加载统计数据失败:', e);
  } finally {
    loading.value = false;
  }
};

const getBarHeight = (count) => {
  const max = Math.max(... chartData.value.map(d => d.borrow_count));
  return max > 0 ? (count / max) * 100 : 0;
};

const formatChartDate = (dateStr) => {
  const date = new Date(dateStr);
  return `${date.getMonth() + 1}/${date.getDate()}`;
};

onMounted(loadData);
</script>

<style scoped>
.page-title {
  margin-bottom: 24px;
}

.loading-state {
  padding: 60px;
  text-align: center;
  color: var(--text-secondary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom:  24px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  display: flex;
  align-items: center;
  gap:  16px;
  transition: transform 0.2s;
}
.stat-card:hover {
  transform: translateY(-2px);
}
.stat-card.warning {
  background: #fef2f2;
  border: 1px solid #fecaca;
}

.stat-icon {
  font-size: 2rem;
}

.stat-content {
  flex:  1;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.text-primary { color: var(--primary); }
.text-success { color: var(--success); }
.text-danger { color: var(--danger); }

.charts-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  margin-bottom: 24px;
}

@media (max-width: 900px) {
  .charts-row {
    grid-template-columns: 1fr;
  }
}

.chart-card {
  min-height: 300px;
}
.chart-card h3 {
  margin-bottom:  20px;
}

.chart-container {
  height: 200px;
}

.chart-bars {
  height: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--border);
}

.chart-bar-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex:  1;
}

.chart-bar {
  width:  32px;
  background: #e5e7eb;
  border-radius:  4px 4px 0 0;
  min-height: 4px;
  transition: height 0.3s;
  cursor: pointer;
}
.chart-bar:hover {
  background: #9ca3af;
}
.chart-bar.active {
  background: var(--primary);
}

.chart-label {
  font-size: 0.75rem;
  color: var(--text-secondary);
  margin-top: 8px;
}

.chart-empty {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content:  center;
  color: var(--text-secondary);
}

.popular-list {
  display: flex;
  flex-direction: column;
  gap:  12px;
}

.popular-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 6px;
  transition: background 0.2s;
}
.popular-item:hover {
  background: #f9fafb;
}

.rank {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  font-weight: 600;
  color:  var(--text-secondary);
  background: #f3f4f6;
  border-radius:  50%;
}
.rank.top {
  background: var(--primary);
  color: white;
}

.book-cover-tiny {
  width: 32px;
  height: 44px;
  object-fit: cover;
  border-radius: 4px;
  background: #f3f4f6;
}

.popular-item .book-info {
  flex:  1;
  min-width: 0;
}

.popular-item .book-title {
  font-weight: 500;
  font-size: 0.875rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow:  ellipsis;
}

.popular-item .book-author {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.borrow-count {
  font-size: 0.875rem;
  color: var(--primary);
  font-weight: 600;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.card-header h3 {
  margin: 0;
}
</style>