<template>
  <div>
    <div class="card">
      <h2>📊 借阅管理</h2>
      
      <!-- 筛选器 -->
      <div class="filters">
        <input 
          v-model.number="filters.user_id" 
          type="number" 
          class="input" 
          placeholder="用户ID"
          @keyup.enter="loadRecords"
        />
        
        <input 
          v-model.number="filters.book_id" 
          type="number" 
          class="input" 
          placeholder="图书ID"
          @keyup.enter="loadRecords"
        />

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
          @change="loadRecords"
        />

        <input 
          v-model="filters.end_date" 
          type="date" 
          class="input" 
          @change="loadRecords"
        />

        <button class="btn" @click="loadRecords">🔍 搜索</button>
        <button class="btn btn-secondary" @click="resetFilters">重置</button>
      </div>
    </div>

    <!-- 借阅记录列表 -->
    <div class="card" style="padding: 0;">
      <div v-if="loading" class="loading-state">加载中...</div>
      
      <table v-else-if="records.length > 0">
        <thead>
          <tr>
            <th>ID</th>
            <th>用户</th>
            <th>图书信息</th>
            <th>借阅日期</th>
            <th>应还日期</th>
            <th>归还日期</th>
            <th>状态</th>
            <th>逾期/罚金</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in records" :key="record.id">
            <td class="text-secondary">{{ record.id }}</td>
            <td>
              <div class="user-info">
                <div class="font-bold">{{ record.user?.username }}</div>
                <div class="text-sm text-secondary">ID: {{ record.user?.id }}</div>
              </div>
            </td>
            <td>
              <div class="book-info-cell">
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
              <span class="badge" :class="getStatusClass(record)">
                {{ getStatusText(record) }}
              </span>
            </td>
            <td>
              <div v-if="record.is_overdue || record.overdue_days > 0">
                <div class="text-danger">{{ record.overdue_days }} 天</div>
                <div v-if="record.fine > 0" class="text-sm text-danger">
                  ¥{{ record.fine.toFixed(2) }}
                </div>
              </div>
              <span v-else class="text-secondary">-</span>
            </td>
            <td class="text-right">
              <div v-if="record.status === 'borrowed'" class="flex justify-end gap-2">
                <button class="btn btn-sm" @click="openReturnModal(record)">
                  归还
                </button>
                <button 
                  class="btn btn-sm btn-secondary"
                  :disabled="!record.can_renew"
                  @click="handleRenew(record)"
                >
                  续借
                </button>
              </div>
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

    <!-- 归还弹窗 -->
    <div v-if="showReturnModal" class="modal-overlay" @click.self="showReturnModal = false">
      <div class="modal">
        <h3>📚 处理归还</h3>
        <p style="margin-bottom: 16px;">
          用户：{{ returningRecord?.user?.username }}<br/>
          图书：《{{ returningRecord?.book?.title }}》
        </p>
        
        <div class="form-group">
          <label>图书状况</label>
          <select v-model="returnForm.condition" class="input">
            <option value="good">完好</option>
            <option value="damaged">损坏</option>
            <option value="lost">丢失</option>
          </select>
        </div>
        
        <div class="form-group">
          <label>备注（选填）</label>
          <textarea 
            v-model="returnForm.remark" 
            class="input" 
            rows="2" 
            style="height: auto;"
          ></textarea>
        </div>

        <div v-if="returningRecord?.is_overdue" class="overdue-warning">
          ⚠️ 已逾期 {{ returningRecord.overdue_days }} 天，预计罚款 ¥{{ (returningRecord.overdue_days * 1).toFixed(2) }}
        </div>
        
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showReturnModal = false">取消</button>
          <button class="btn" :disabled="returning" @click="confirmReturn">
            {{ returning ? '处理中...' : '确认归还' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { getBorrowRecords, returnBook, renewBook } from '../api';
import { formatDate } from '../utils/format';
import { $message } from '../utils/toast';

const records = ref([]);
const loading = ref(false);
const total = ref(0);
const totalPages = ref(1);
const page = ref(1);
const limit = 15;

const filters = reactive({
  user_id: null,
  book_id: null,
  status: '',
  start_date: '',
  end_date: ''
});

const showReturnModal = ref(false);
const returningRecord = ref(null);
const returnForm = reactive({
  condition: 'good',
  remark: ''
});
const returning = ref(false);

const loadRecords = async () => {
  loading.value = true;
  try {
    const params = {
      page: page.value,
      limit,
      sort_by: 'borrow_date',
      order: 'desc'
    };

    if (filters.user_id) params.user_id = filters.user_id;
    if (filters.book_id) params.book_id = filters.book_id;
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
  filters.user_id = null;
  filters.book_id = null;
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

const getStatusClass = (record) => {
  if (record.is_overdue) return 'badge-danger';
  if (record.status === 'borrowed') return 'badge-primary';
  if (record.status === 'returned') return 'badge-success';
  return 'badge-secondary';
};

const getStatusText = (record) => {
  if (record.is_overdue) return '已逾期';
  const map = {
    borrowed: '借阅中',
    returned: '已归还',
    overdue: '已逾期'
  };
  return map[record.status] || record.status;
};

const openReturnModal = (record) => {
  returningRecord.value = record;
  returnForm.condition = 'good';
  returnForm.remark = '';
  showReturnModal.value = true;
};

const confirmReturn = async () => {
  returning.value = true;
  try {
    const data = { condition: returnForm.condition };
    if (returnForm.remark) data.remark = returnForm.remark;
    
    const res = await returnBook(returningRecord.value.id, data);
    
    if (res.is_overdue) {
      $message.warning(`归还成功！逾期 ${res.overdue_days} 天，罚款 ¥${res.fine.toFixed(2)}`);
    } else {
      $message.success('归还成功！');
    }
    
    showReturnModal.value = false;
    loadRecords();
  } catch (error) {
    console.error('归还失败:', error);
  } finally {
    returning.value = false;
  }
};

const handleRenew = async (record) => {
  if (!record.can_renew) {
    $message.warning('该借阅记录无法续借');
    return;
  }

  try {
    const res = await renewBook(record.id, { renew_days: 30 });
    $message.success(`续借成功！新的到期日: ${formatDate(res.new_due_date)}`);
    loadRecords();
  } catch (error) {
    console.error('续借失败:', error);
  }
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
  max-width: 160px;
}

.user-info,
.book-info-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.book-title {
  font-weight: 600;
  max-width: 200px;
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

.font-bold {
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

.overdue-warning {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
  padding: 12px;
  border-radius: 6px;
  font-size: 0.875rem;
  margin-top: 16px;
}
</style>
