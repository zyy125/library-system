<template>
  <div>
    <!-- Tab 切换 -->
    <div class="tabs-container">
      <div class="tabs">
        <button 
          v-for="t in tabs" :key="t.key"
          class="tab-btn" 
          :class="{ active: tab === t.key }"
          @click="tab = t.key"
        >
          {{ t.name }}
          <span v-if="t.count !== undefined && t.count > 0" class="tab-count">{{ t.count }}</span>
        </button>
      </div>
    </div>

    <!-- 当前借阅概览 -->
    <div v-if="tab === 'current' && summary" class="card summary-card">
      <div class="summary-item">
        <span class="label">当前借阅</span>
        <span class="value">{{ summary.borrowing_count }} / {{ summary.borrow_limit }}</span>
      </div>
      <div class="summary-item" v-if="summary.overdue_count > 0">
        <span class="label text-danger">⚠️ 逾期图书</span>
        <span class="value text-danger">{{ summary.overdue_count }} 本</span>
      </div>
      <div class="summary-item" v-if="summary.total_fine > 0">
        <span class="label text-warning">💰 待缴罚款</span>
        <span class="value text-warning">¥{{ summary.total_fine.toFixed(2) }}</span>
      </div>
    </div>

    <div class="card" style="padding: 0; min-height: 400px;">
      <div v-if="loading" class="loading-state">加载中...</div>
      
      <!-- 表格内容 -->
      <table v-else-if="dataList.length > 0">
        <thead>
          <tr>
            <th>图书信息</th>
            <th>{{ tab === 'reserved' ? '预约时间' : '借阅时间' }}</th>
            <th>{{ getTimeColumnHeader() }}</th>
            <th>状态</th>
            <th v-if="tab === 'current'">续借次数</th>
            <th v-if="tab === 'history'">罚款</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in dataList" :key="item.id">
            <td>
              <div class="book-info">
                <img :src="item.book?.cover_url || defaultCover" class="book-cover-small">
                <div>
                  <div class="book-title">{{ item.book?.title }}</div>
                  <div class="text-sm text-secondary">{{ item.book?.author }}</div>
                </div>
              </div>
            </td>
            <td>{{ formatDate(item.borrow_date || item.reserved_at) }}</td>
            <td>
              <template v-if="tab === 'reserved'">
                <div>排队位置:  <strong>#{{ item.queue_position }}</strong></div>
                <div class="text-sm text-secondary">有效期至: {{ formatDate(item.expires_at) }}</div>
              </template>
              <template v-else-if="tab === 'history'">
                {{ formatDate(item.return_date) }}
              </template>
              <template v-else>
                <div>{{ formatDate(item.due_date) }}</div>
                <div v-if="! item.is_overdue && item.days_until_due !== undefined" class="text-sm text-secondary">
                  剩余 {{ item.days_until_due }} 天
                </div>
                <div v-if="item.is_overdue" class="text-sm text-danger">
                  已逾期 {{ item.overdue_days }} 天
                </div>
              </template>
            </td>
            <td>
              <span class="badge" :class="getStatusInfo(item).class">
                {{ getStatusInfo(item).label }}
              </span>
            </td>
            <td v-if="tab === 'current'">
              <span class="text-sm">
                {{ item.renew_count || 0 }} / {{ item.max_renew_count || 2 }}
              </span>
              <span v-if="! item.can_renew" class="text-sm text-secondary"> (不可续)</span>
            </td>
            <td v-if="tab === 'history'">
              <span v-if="item.fine > 0" class="text-danger font-bold">¥{{ item.fine.toFixed(2) }}</span>
              <span v-else class="text-secondary">-</span>
            </td>
            <td class="text-right">
              <div v-if="tab === 'current'" class="flex justify-end gap-2">
                <button class="btn btn-sm" @click="openReturnModal(item)">归还</button>
                <button 
                  class="btn btn-sm btn-secondary" 
                  :disabled="! item.can_renew"
                  @click="handleRenew(item)"
                >
                  续借
                </button>
              </div>
              <div v-else-if="tab === 'reserved'" class="flex justify-end">
                <button class="btn btn-sm btn-danger" @click="handleCancel(item.id)">取消预约</button>
              </div>
              <span v-else class="text-secondary text-sm">已归档</span>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-else class="empty-state">
        <template v-if="tab === 'current'">
          <p>📖 暂无借阅中的图书</p>
          <p class="text-sm">去 <router-link to="/books" class="link">图书大厅</router-link> 看看吧</p>
        </template>
        <template v-else-if="tab === 'history'">
          <p>📜 暂无历史借阅记录</p>
        </template>
        <template v-else>
          <p>⏳ 暂无预约</p>
        </template>
      </div>

      <!-- 历史记录分页 -->
      <div class="pagination" v-if="tab === 'history' && dataList.length > 0">
        <div class="page-info">共 {{ historyTotal }} 条记录</div>
        <div class="page-controls">
          <button class="btn btn-secondary btn-sm" :disabled="historyPage <= 1" @click="loadHistory(historyPage - 1)">上一页</button>
          <span class="page-num">{{ historyPage }} / {{ historyTotalPages }}</span>
          <button class="btn btn-secondary btn-sm" :disabled="historyPage >= historyTotalPages" @click="loadHistory(historyPage + 1)">下一页</button>
        </div>
      </div>
    </div>

    <!-- 还书弹窗 -->
    <div v-if="showReturnModal" class="modal-overlay" @click.self="showReturnModal = false">
      <div class="modal">
        <h3>📚 归还图书</h3>
        <p style="margin-bottom: 16px;">确认归还《{{ returningBook?.book?.title }}》？</p>
        
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
          <textarea v-model="returnForm.remark" class="input" rows="2" placeholder="如有损坏请说明情况" style="height: auto;"></textarea>
        </div>

         <div v-if="returningBook?.is_overdue" class="overdue-warning">
          ⚠️ 此书已逾期 {{ returningBook.overdue_days }} 天，预计罚款 ¥{{ (returningBook.overdue_days * 1).toFixed(2) }}
        </div>
        
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showReturnModal = false">取消</button>
          <button class="btn" :disabled="returning" @click="confirmReturn">
            {{ returning ? '处理中...' :  '确认归还' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted, computed } from 'vue';
import { getCurrentBorrows, getBorrowRecords, returnBook, renewBook, getMyReservations, cancelReservation } from '../api';
import { formatDate, getStatusConfig } from '../utils/format';
import { $message } from '../utils/toast';

const defaultCover = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="48" height="64" viewBox="0 0 48 64"%3E%3Crect fill="%23f3f4f6" width="48" height="64"/%3E%3C/svg%3E';

const tab = ref('current');
const dataList = ref([]);
const loading = ref(false);
const summary = ref(null);

// 历史记录分页
const historyPage = ref(1);
const historyTotal = ref(0);
const historyTotalPages = ref(1);

// 还书相关
const showReturnModal = ref(false);
const returningBook = ref(null);
const returnForm = reactive({ condition: 'good', remark: '' });
const returning = ref(false);

const tabs = computed(() => [
  { key: 'current', name: '📖 当前借阅', count:  summary.value?.borrowing_count },
  { key:  'history', name: '📜 历史记录' },
  { key: 'reserved', name: '⏳ 我的预约' }
]);

const getTimeColumnHeader = () => {
  if (tab.value === 'reserved') return '排队信息';
  if (tab.value === 'history') return '归还时间';
  return '应还日期';
};

const getStatusInfo = (item) => {
  if (tab.value === 'reserved') return getStatusConfig(item.status);
  return getStatusConfig(item.status, item.is_overdue);
};

const loadData = async () => {
  loading.value = true;
  dataList.value = [];
  
  try {
    if (tab.value === 'current') {
      const res = await getCurrentBorrows();
      dataList.value = res.records || [];
      summary.value = {
        borrowing_count: res.borrowing_count || 0,
        borrow_limit: res.borrow_limit || 5,
        overdue_count: res.overdue_count || 0,
        total_fine:  res.total_fine || 0
      };
    } else if (tab.value === 'history') {
      await loadHistory(1);
    } else {
      const res = await getMyReservations();
      dataList.value = res.reservations || [];
    }
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
};

const loadHistory = async (page) => {
  loading.value = true;
  try {
    const res = await getBorrowRecords({ page, limit: 10, status: 'returned' });
    dataList.value = res.records || [];
    historyPage.value = res.page || 1;
    historyTotal.value = res.total || 0;
    historyTotalPages.value = res.total_pages || 1;
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
};

const openReturnModal = (item) => {
  returningBook.value = item;
  returnForm.condition = 'good';
  returnForm.remark = '';
  showReturnModal.value = true;
};

const confirmReturn = async () => {
  returning.value = true;
  try {
    const data = { condition: returnForm.condition };
    if (returnForm.remark) data.remark = returnForm.remark;
    
    const res = await returnBook(returningBook.value.id, data);
    
    if (res.is_overdue) {
      $message.warning(`归还成功！逾期 ${res.overdue_days} 天，罚款 ¥${res.fine.toFixed(2)}`);
    } else {
      $message.success('归还成功！');
    }
    
    showReturnModal.value = false;
    loadData();
  } catch (e) {
    console.error(e);
  } finally {
    returning.value = false;
  }
};

const handleRenew = async (item) => {
  if (! item.can_renew) {
    $message.warning('该图书无法续借');
    return;
  }
  
  try {
    const res = await renewBook(item.id, { renew_days: 30 });
    $message.success(`续借成功！新的到期日:  ${formatDate(res.new_due_date)}`);
    loadData();
  } catch (e) {
    console.error(e);
  }
};

const handleCancel = async (id) => {
  if (confirm('确定取消此预约？')) {
    try {
      await cancelReservation(id);
      $message.warning('预约已取消');
      loadData();
    } catch (e) {
      console.error(e);
    }
  }
};

watch(tab, loadData);
onMounted(loadData);
</script>

<style scoped>
.tabs-container {
  margin-bottom: 24px;
}
.tabs {
  display:  flex;
  gap: 8px;
  border-bottom: 2px solid var(--border);
}
.tab-btn {
  padding: 12px 24px;
  border: none;
  background: transparent;
  font-size: 1rem;
  color: var(--text-secondary);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap:  8px;
}
.tab-btn:hover { color: var(--text-main); }
.tab-btn.active { 
  color: var(--primary); 
  border-bottom-color:  var(--primary); 
  font-weight: 600; 
}
.tab-count {
  background: var(--primary);
  color: white;
  font-size: 0.75rem;
  padding: 2px 8px;
  border-radius: 10px;
}

.summary-card {
  display: flex;
  gap: 32px;
  padding: 16px 24px;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
}
.summary-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.summary-item .label { font-size: 0.875rem; color: var(--text-secondary); }
.summary-item .value { font-size: 1.25rem; font-weight: 600; }

.book-info {
  display: flex;
  align-items: center;
  gap: 12px;
}
.book-cover-small {
  width: 40px;
  height: 56px;
  object-fit: cover;
  border-radius: 4px;
  background: #f3f4f6;
}
.book-title {
  font-weight: 600;
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.text-sm { font-size: 0.8rem; }
.text-secondary { color: var(--text-secondary); }
.text-danger { color: var(--danger); }
.text-warning { color: var(--warning); }
.font-bold { font-weight:  600; }

.loading-state, .empty-state {
  padding: 60px;
  text-align: center;
  color: var(--text-secondary);
}
.empty-state p { margin: 8px 0; }
.link { color: var(--primary); text-decoration: none; }
.link:hover { text-decoration: underline; }

.pagination {
  padding: 16px;
  border-top: 1px solid var(--border);
  display: flex;
  justify-content:  space-between;
  align-items:  center;
}
.page-info { color: var(--text-secondary); font-size: 0.875rem; }
.page-controls { display: flex; gap: 8px; align-items: center; }
.page-num { color: var(--text-secondary); font-size: 0.875rem; padding: 0 8px; }

.overdue-warning {
  background: #fef2f2;
  border:  1px solid #fecaca;
  color: #991b1b;
  padding: 12px;
  border-radius: 6px;
  font-size: 0.875rem;
  margin-top: 16px;
}
</style>