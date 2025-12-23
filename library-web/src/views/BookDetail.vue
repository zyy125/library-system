<template>
  <div class="book-detail-page">
    <div v-if="loading" class="loading-state">加载中...</div>
    
    <div v-else-if="book" class="book-detail">
      <div class="card">
        <button class="btn btn-secondary btn-sm" @click="$router.back()" style="margin-bottom: 16px;">
          ← 返回
        </button>

        <div class="book-header">
          <img 
            v-if="book.cover_url" 
            :src="book.cover_url" 
            :alt="book.title"
            class="book-cover-large"
          />
          <div v-else class="book-cover-large no-cover">
            📖
          </div>

          <div class="book-info">
            <h1>{{ book.title }}</h1>
            <div class="meta-row">
              <span class="meta-label">作者:</span>
              <span>{{ book.author }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">ISBN:</span>
              <span>{{ book.isbn }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">出版社:</span>
              <span>{{ book.publisher }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">出版日期:</span>
              <span>{{ book.publish_date }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">分类:</span>
              <span class="badge badge-primary">{{ book.category?.name }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">价格:</span>
              <span class="price">¥{{ book.price }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">库存状态:</span>
              <span v-if="book.available > 0" class="badge badge-success">
                有库存 ({{ book.available }}/{{ book.stock }})
              </span>
              <span v-else class="badge badge-danger">
                已借完 (0/{{ book.stock }})
              </span>
            </div>
            <div class="meta-row">
              <span class="meta-label">借阅次数:</span>
              <span>{{ book.borrow_count }} 次</span>
            </div>
          </div>
        </div>

        <div class="book-description" v-if="book.description">
          <h3>图书简介</h3>
          <p>{{ book.description }}</p>
        </div>

        <div class="actions" v-if="!isAdmin">
          <button 
            v-if="book.available > 0"
            class="btn" 
            @click="handleBorrow"
            :disabled="borrowing"
          >
            {{ borrowing ? '借阅中...' : '立即借阅' }}
          </button>
          <button 
            v-else
            class="btn btn-secondary" 
            @click="handleReserve"
            :disabled="reserving"
          >
            {{ reserving ? '预约中...' : '预约图书' }}
          </button>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">图书不存在</div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import { getBookDetail, borrowBook, reserveBook } from '../api';
import { getUser } from '../utils/auth';
import { $message } from '../utils/toast';

const route = useRoute();
const book = ref(null);
const loading = ref(true);
const borrowing = ref(false);
const reserving = ref(false);

const user = computed(() => getUser());
const isAdmin = computed(() => user.value?.role === 'admin');

const loadBookDetail = async () => {
  loading.value = true;
  try {
    book.value = await getBookDetail(route.params.id);
  } catch (error) {
    console.error('加载图书详情失败:', error);
  } finally {
    loading.value = false;
  }
};

const handleBorrow = async () => {
  borrowing.value = true;
  try {
    await borrowBook({ book_id: parseInt(route.params.id) });
    $message.success('借阅成功！');
    await loadBookDetail(); // 刷新图书信息
  } catch (error) {
    console.error('借阅失败:', error);
  } finally {
    borrowing.value = false;
  }
};

const handleReserve = async () => {
  reserving.value = true;
  try {
    await reserveBook({ book_id: parseInt(route.params.id) });
    $message.success('预约成功！');
  } catch (error) {
    console.error('预约失败:', error);
  } finally {
    reserving.value = false;
  }
};

onMounted(() => {
  loadBookDetail();
});
</script>

<style scoped>
.book-detail-page {
  max-width: 900px;
  margin: 0 auto;
}

.book-header {
  display: flex;
  gap: 32px;
  margin-bottom: 32px;
}

.book-cover-large {
  width: 200px;
  height: 280px;
  object-fit: cover;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  flex-shrink: 0;
}

.book-cover-large.no-cover {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 5rem;
}

.book-info {
  flex: 1;
}

.book-info h1 {
  font-size: 1.75rem;
  margin-bottom: 20px;
  color: var(--text-main);
}

.meta-row {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  font-size: 0.9rem;
}

.meta-label {
  font-weight: 600;
  color: var(--text-secondary);
  width: 100px;
  flex-shrink: 0;
}

.price {
  color: var(--danger);
  font-weight: 600;
  font-size: 1.1rem;
}

.book-description {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
}

.book-description h3 {
  margin-bottom: 12px;
  color: var(--text-main);
}

.book-description p {
  color: var(--text-secondary);
  line-height: 1.8;
}

.actions {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
  display: flex;
  gap: 12px;
}

.actions .btn {
  padding: 12px 32px;
  font-size: 1rem;
}

@media (max-width: 768px) {
  .book-header {
    flex-direction: column;
  }

  .book-cover-large {
    width: 100%;
    max-width: 200px;
    margin: 0 auto;
  }
}
</style>
