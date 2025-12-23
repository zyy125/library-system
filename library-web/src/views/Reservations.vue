<template>
  <div>
    <div class="card">
      <h2>📅 我的预约</h2>
      <p class="hint">预约有效期48小时，图书归还后将通知您，请及时借阅</p>
    </div>

    <div class="card" style="padding: 0;">
      <div v-if="loading" class="loading-state">加载中...</div>
      
      <table v-else-if="reservations.length > 0">
        <thead>
          <tr>
            <th>图书信息</th>
            <th>预约时间</th>
            <th>过期时间</th>
            <th>排队位置</th>
            <th>状态</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="reservation in reservations" :key="reservation.id">
            <td>
              <div class="book-info-cell">
                <img 
                  v-if="reservation.book?.cover_url" 
                  :src="reservation.book.cover_url" 
                  class="book-cover"
                />
                <div>
                  <div class="book-title">{{ reservation.book?.title }}</div>
                  <div class="text-sm text-secondary">
                    ISBN: {{ reservation.book?.isbn }}
                  </div>
                </div>
              </div>
            </td>
            <td>{{ formatDate(reservation.reserved_at) }}</td>
            <td>
              <span :class="isExpiringSoon(reservation.expires_at) ? 'text-danger' : ''">
                {{ formatDate(reservation.expires_at) }}
              </span>
            </td>
            <td>
              <span class="queue-position">第 {{ reservation.queue_position }} 位</span>
            </td>
            <td>
              <span class="badge" :class="getStatusClass(reservation.status)">
                {{ getStatusText(reservation.status) }}
              </span>
            </td>
            <td class="text-right">
              <button 
                v-if="reservation.status === 'waiting'"
                class="btn btn-sm btn-danger"
                @click="handleCancel(reservation.id)"
              >
                取消预约
              </button>
              <span v-else class="text-secondary">-</span>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else class="empty-state">
        <p>📅 暂无预约记录</p>
        <p class="text-sm text-secondary">当图书库存为0时可进行预约</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getMyReservations, cancelReservation } from '../api';
import { formatDate } from '../utils/format';
import { $message } from '../utils/toast';

const reservations = ref([]);
const loading = ref(false);

const loadReservations = async () => {
  loading.value = true;
  try {
    const res = await getMyReservations();
    reservations.value = res.reservations || [];
  } catch (error) {
    console.error('加载预约列表失败:', error);
  } finally {
    loading.value = false;
  }
};

const handleCancel = async (id) => {
  if (!confirm('确定要取消这个预约吗？')) return;

  try {
    await cancelReservation(id);
    $message.success('预约已取消');
    await loadReservations();
  } catch (error) {
    console.error('取消预约失败:', error);
  }
};

const isExpiringSoon = (expiresAt) => {
  const expires = new Date(expiresAt);
  const now = new Date();
  const hours = (expires - now) / (1000 * 60 * 60);
  return hours < 24 && hours > 0;
};

const getStatusClass = (status) => {
  const map = {
    waiting: 'badge-primary',
    available: 'badge-success',
    expired: 'badge-secondary',
    cancelled: 'badge-danger'
  };
  return map[status] || 'badge-secondary';
};

const getStatusText = (status) => {
  const map = {
    waiting: '等待中',
    available: '可借阅',
    expired: '已过期',
    cancelled: '已取消'
  };
  return map[status] || status;
};

onMounted(() => {
  loadReservations();
});
</script>

<style scoped>
.hint {
  color: var(--text-secondary);
  font-size: 0.875rem;
  margin: 12px 0 0 0;
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

.queue-position {
  font-weight: 600;
  color: var(--primary);
}
</style>
