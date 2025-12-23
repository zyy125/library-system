<template>
  <div>
    <div class="card">
      <h2>📖 借阅历史</h2>
      
      <!-- 筛选器 -->
      <div class="filters">
        <select v-model="filters.status" class="input" @change="loadRecords">
          <option value="">全部状态</option>
          <option value="borrowed">借阅中</option>
          <option value="returned">已归还</option>
          <option value="overdue">已逾期</option>
        </select>

        <input 
          v-model="filters.start_date" 
          type="date" 
          class="input" 
          placeholder="开始日期"
          @change="loadRecords"
        />

        <input 
          v-model="filters.end_date" 
          type="date" 
          class="input" 
          placeholder="结束日期"
          @change="loadRecords"
        />

        <button class="btn btn-secondary" @click="resetFilters">重置</button>
      </div>
    </div>

    <!-- 记录列表 -->
    <div class="card" style="padding: 0;">
      <div v-if="loading" class="loading-state">加载中...</div>
      
      <table v-else-if="records.length > 0">
        <thead>
          <tr>
            <th>图书信息</th>
            <th>借阅日期</th>
            <th>应还日期</th>
            <th>归还日期</th>
            <th>状态</th>
            <th>逾期</th>
            <th>罚金</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in records" :key="record.id">
            <td>
              <div class="book-info-cell">
                <img 
                  v-if="record.book?.cover_url" 
                  :src="record.book.cover_url" 
                  class="book-cover"
                />
                <div>
                  <div class="book-title">{{ record.book?.title }}</div>
                  <div class="text-sm text-secondary">{{ record.book?.author }}</div>
                </div>
              </div>
            </td>
            <td>{{ formatDate(record.borrow_date) }}</td>
            <td>{{ formatDate(record.due_date) }}</td>
            <td>
              <span v-if="record.return_date">{{ formatDate(record.return_date) }}</span>
              <span v-else class="text-secondary">-</span>
            </td>
            <td>
              <span class="badge" :class="getStatusClass(record.status)">
                {{ getStatusText(record.status) }}
              </span>
            </td>
            <td>
              <span v-if="record.is_overdue" class="text-danger">
                {{ record.overdue_days }} 天
              </span>
              <span v-else class="text-secondary">-</span>
            </td>
            <td>
              <span v-if="record.fine > 0" class="text-danger">
                ¥{{ record.fine.toFixed(2) }}
              </span>
              <span v-else class="text-secondary">-</span>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else class="empty-state">暂无借阅记录</div>

      <!-- 分页 -->
      <div class="pagination" v-if="total > 0">
        <div class="page-info">
          共 {{ total }} 条记录，第 {{ page }} / {{ totalPages }} 页
        </div>
        <div class="page-controls">
          <button class="btn btn-secondary btn-sm" :disabled="page <= 1" @click="changePage(-1)">
            上一页
          </button>
          <button class="btn btn-secondary btn-sm" :disabled="page >= totalPages" @click="changePage(1)">
            下一页
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { getBorrowRecords } from '../api';
import { formatDate } from '../utils/format';

const records = ref([]);
const loading = ref(false);
const total = ref(0);
const totalPages = ref(1);
const page = ref(1);
const limit = 10;

const filters = reactive({
  status: '',
  start_date: '',
  end_date: ''
});

const loadRecords = async () => {
  loading.value = true;
  try {
    const params = {
      page: page.value,
      limit,
      sort_by: 'borrow_date',
      order: 'desc'
    };

    if (filters.status) params.status = filters.status;
    if (filters.start_date) params.start_date = filters.start_date;
    if (filters.end_date) params.end_date = filters.end_date;

    const res = await getBorrowRecords(params);
    records.value = res.records || [];
    total.value = res.total || 0;
    totalPages.value = res.total_pages || 1;
  } catch (error) {
    console.error('加载借阅记录失败:', error);
  } finally {
    loading.value = false;
  }
};

const resetFilters = () => {
  filters.status = '';
  filters.start_date = '';
  filters.end_date = '';
  page.value = 1;
  loadRecords();
};

const changePage = (delta) => {
  page.value += delta;
  loadRecords();
};

const getStatusClass = (status) => {
  const map = {
    borrowed: 'badge-primary',
    returned: 'badge-success',
    overdue: 'badge-danger'
  };
  return map[status] || 'badge-secondary';
};

const getStatusText = (status) => {
  const map = {
    borrowed: '借阅中',
    returned: '已归还',
    overdue: '已逾期'
  };
  return map[status] || status;
};

onMounted(() => {
  loadRecords();
});
</script>

<style scoped>
.filters {
  display: flex;
  gap: 12px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.filters .input {
  max-width: 180px;
}

.book-info-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.book-title {
  font-weight: 600;
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.text-sm {
  font-size: 0.8rem;
}

.text-secondary {
  color: var(--text-secondary);
}

.text-danger {
  color: var(--danger);
  font-weight: 600;
}

.pagination {
  padding: 16px;
  border-top: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-info {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.page-controls {
  display: flex;
  gap: 8px;
}
</style>
