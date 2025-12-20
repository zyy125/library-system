<template>
  <div>
    <!-- Tab 切换 -->
    <div class="tabs flex gap-2 mb-4">
      <button 
        v-for="t in tabs" :key="t.key"
        class="tab-btn" 
        :class="{ active: tab === t.key }"
        @click="tab = t.key"
      >
        {{ t.name }}
      </button>
    </div>

    <div class="card" style="padding: 0; min-height: 400px;">
      <!-- 表格内容 -->
      <table v-if="dataList.length > 0">
        <thead>
          <tr>
            <th>图书信息</th>
            <th>借阅时间</th>
            <th>{{ tab === 'history' ? '归还时间' : '截止/逾期' }}</th>
            <th>状态</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in dataList" :key="item.id">
            <td>
              <div style="font-weight: 600;">{{ item.book?.title }}</div>
              <div style="font-size: 0.8rem; color: var(--text-secondary);">ISBN: {{ item.book?.isbn || '-' }}</div>
            </td>
            <!-- 使用工具函数格式化时间 -->
            <td>{{ formatDate(item.borrow_date || item.reserved_at) }}</td>
            <td>
              <span v-if="tab === 'reserved'">位置: {{ item.queue_position }}</span>
              <span v-else>{{ formatDate(item.return_date || item.due_date) }}</span>
            </td>
            <td>
              <!-- 业务要求：状态标签化 -->
              <span class="badge" :class="getStatusInfo(item).class">
                {{ getStatusInfo(item).label }}
              </span>
              <span v-if="item.is_overdue" class="badge badge-danger" style="margin-left: 5px;">
                逾期 {{ item.overdue_days }} 天
              </span>
            </td>
            <td class="text-right">
              <div v-if="tab === 'current'" class="flex justify-end gap-2">
                <button class="btn btn-sm" @click="handleReturn(item.id)">归还</button>
                <button class="btn btn-sm btn-secondary" @click="handleRenew(item.id)">续借</button>
              </div>
              <div v-if="tab === 'reserved'" class="flex justify-end">
                <button class="btn btn-sm btn-danger" @click="handleCancel(item.id)">取消预约</button>
              </div>
              <span v-if="tab === 'history'" style="color: var(--text-secondary); font-size: 0.8rem;">已归档</span>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-else class="empty-state">
        暂无记录
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { getCurrentBorrows, getBorrowRecords, returnBook, renewBook, getMyReservations, cancelReservation } from '../api';
import { formatDate, getStatusConfig } from '../utils/format';
import { $message } from '../utils/toast';

const tab = ref('current');
const dataList = ref([]);
const tabs = [
  { key: 'current', name: '📖 当前借阅' },
  { key: 'history', name: '📜 历史记录' },
  { key: 'reserved', name: '⏳ 我的预约' }
];

const getStatusInfo = (item) => {
  if (tab.value === 'reserved') return getStatusConfig(item.status);
  return getStatusConfig(item.status, item.is_overdue);
};

const loadData = async () => {
  dataList.value = [];
  try {
    if (tab.value === 'current') {
      const res = await getCurrentBorrows();
      dataList.value = res.records || [];
    } else if (tab.value === 'history') {
      const res = await getBorrowRecords({ page: 1, limit: 20 });
      dataList.value = res.records || [];
    } else {
      const res = await getMyReservations();
      dataList.value = res.reservations || [];
    }
  } catch(e) {}
};

const handleReturn = async (id) => {
  if(confirm('确认归还？')) { 
    try {
        await returnBook(id, { condition: 'good' }); 
        $message.success('归还成功');
        loadData(); 
    } catch(e) {}
  }
};

const handleRenew = async (id) => {
  try { 
    await renewBook(id, {}); 
    $message.success('续借成功'); 
    loadData(); 
  } catch(e) {}
};

const handleCancel = async (id) => {
  if(confirm('取消此预约？')) { 
      try {
        await cancelReservation(id); 
        $message.warning('预约已取消');
        loadData(); 
      } catch(e) {}
  }
};

watch(tab, loadData);
onMounted(loadData);
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
.tab-btn {
  padding: 10px 20px; border: none; background: transparent;
  font-size: 1rem; color: var(--text-secondary); cursor: pointer;
  border-bottom: 2px solid transparent; transition: all 0.2s;
}
.tab-btn.active { color: var(--primary); border-bottom-color: var(--primary); font-weight: 600; }
.empty-state { padding: 60px; text-align: center; color: var(--text-secondary); }
</style>