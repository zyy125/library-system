<template>
  <div class="popular-page">
    <div class="page-header">
      <div class="header-content">
        <h2 class="page-title">🔥 热门图书排行</h2>
        <p class="page-subtitle">发现最受欢迎的图书</p>
      </div>
      
      <!-- 筛选器 -->
      <div class="filters">
        <select v-model="filters.period" @change="loadBooks" class="filter-select">
          <option value="7d">近7天</option>
          <option value="30d">近30天</option>
          <option value="90d">近90天</option>
          <option value="all">全部时间</option>
        </select>
        
        <select v-model="filters.limit" @change="loadBooks" class="filter-select">
          <option value="10">Top 10</option>
          <option value="20">Top 20</option>
          <option value="50">Top 50</option>
        </select>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 热门图书列表 -->
    <div v-else-if="books.length > 0" class="books-grid">
      <div 
        v-for="(book, index) in books" 
        :key="book.book_id || book.id"
        class="book-card"
        @click="viewDetails(book.book_id || book.id)"
      >
        <!-- 排名徽章 -->
        <div class="rank-badge" :class="getRankClass(index + 1)">
          <span class="rank-number">#{{ index + 1 }}</span>
        </div>
        
        <!-- 图书封面 -->
        <div class="book-cover-wrapper">
          <img 
            :src="book.cover_url || defaultCover" 
            :alt="book.title"
            class="book-cover"
          />
          <div class="book-overlay">
            <button class="view-btn">查看详情</button>
          </div>
        </div>

        <!-- 图书信息 -->
        <div class="book-info">
          <h3 class="book-title">{{ book.title }}</h3>
          <p class="book-author">{{ book.author }}</p>
          
          <div class="book-stats">
            <div class="stat-item">
              <span class="stat-icon">📖</span>
              <span class="stat-value">{{ book.borrow_count || book.recent_borrow_count }} 次借阅</span>
            </div>
            <div class="stat-item" v-if="book.rating">
              <span class="stat-icon">⭐</span>
              <span class="stat-value">{{ book.rating }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <div class="empty-icon">📚</div>
      <p>暂无热门图书数据</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getPopularBooks } from '../api';
import { $message } from '../utils/toast';

const router = useRouter();
const defaultCover = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="200" height="280" viewBox="0 0 200 280"%3E%3Crect fill="%23f3f4f6" width="200" height="280"/%3E%3Ctext x="50%25" y="50%25" font-size="16" fill="%239ca3af" text-anchor="middle" dy=".3em"%3E暂无封面%3C/text%3E%3C/svg%3E';

const loading = ref(true);
const books = ref([]);
const filters = reactive({
  period: '30d',
  limit: 20
});

// 加载热门图书
const loadBooks = async () => {
  loading.value = true;
  try {
    // API参数映射
    const params = {
      limit: parseInt(filters.limit),
      days: filters.period === '7d' ? 7 : filters.period === '30d' ? 30 : filters.period === '90d' ? 90 : undefined
    };
    
    // 如果是all，不传days参数
    if (filters.period === 'all') {
      delete params.days;
    }

    const res = await getPopularBooks(params);
    books.value = res.books || [];
  } catch (error) {
    console.error('加载热门图书失败:', error);
    $message.error('加载失败，请稍后重试');
    books.value = [];
  } finally {
    loading.value = false;
  }
};

// 获取排名样式类
const getRankClass = (rank) => {
  if (rank === 1) return 'rank-gold';
  if (rank === 2) return 'rank-silver';
  if (rank === 3) return 'rank-bronze';
  return '';
};

// 查看图书详情
const viewDetails = (bookId) => {
  router.push(`/books/${bookId}`);
};

onMounted(loadBooks);
</script>

<style scoped>
.popular-page {
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
  flex-wrap: wrap;
  gap: 20px;
}

.header-content {
  flex: 1;
}

.page-title {
  font-size: 2rem;
  font-weight: 800;
  color: var(--text-main);
  margin: 0 0 8px 0;
}

.page-subtitle {
  color: var(--text-secondary);
  font-size: 1rem;
  margin: 0;
}

.filters {
  display: flex;
  gap: 12px;
}

.filter-select {
  padding: 10px 16px;
  border: 2px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 0.9rem;
  font-weight: 600;
  background: white;
  cursor: pointer;
  transition: var(--transition);
}

.filter-select:hover {
  border-color: var(--primary);
}

.filter-select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--text-secondary);
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid var(--border-light);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.books-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 24px;
}

.book-card {
  background: white;
  border-radius: var(--radius);
  overflow: hidden;
  box-shadow: var(--shadow);
  transition: var(--transition);
  cursor: pointer;
  position: relative;
  border: 1px solid var(--border-light);
}

.book-card:hover {
  transform: translateY(-8px);
  box-shadow: var(--shadow-lg);
}

.rank-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 10;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-weight: 700;
  font-size: 0.9rem;
  backdrop-filter: blur(4px);
}

.rank-badge.rank-gold {
  background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
  color: #333;
  box-shadow: 0 4px 12px rgba(255, 215, 0, 0.4);
}

.rank-badge.rank-silver {
  background: linear-gradient(135deg, #c0c0c0 0%, #e8e8e8 100%);
  color: #333;
  box-shadow: 0 4px 12px rgba(192, 192, 192, 0.4);
}

.rank-badge.rank-bronze {
  background: linear-gradient(135deg, #cd7f32 0%, #e8a87c 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(205, 127, 50, 0.4);
}

.book-cover-wrapper {
  position: relative;
  width: 100%;
  height: 320px;
  overflow: hidden;
  background: var(--border-light);
}

.book-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.book-card:hover .book-cover {
  transform: scale(1.05);
}

.book-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.book-card:hover .book-overlay {
  opacity: 1;
}

.view-btn {
  padding: 10px 24px;
  background: white;
  color: var(--primary);
  border: none;
  border-radius: var(--radius-sm);
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  transition: var(--transition);
}

.view-btn:hover {
  background: var(--primary);
  color: white;
  transform: scale(1.05);
}

.book-info {
  padding: 20px;
}

.book-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-main);
  margin: 0 0 8px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

.book-author {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin: 0 0 16px 0;
}

.book-stats {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.stat-icon {
  font-size: 1rem;
}

.stat-value {
  font-weight: 600;
  color: var(--primary);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 16px;
  opacity: 0.5;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
  }
  
  .filters {
    width: 100%;
  }
  
  .filter-select {
    flex: 1;
  }
  
  .books-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 16px;
  }
  
  .book-cover-wrapper {
    height: 220px;
  }
}
</style>
