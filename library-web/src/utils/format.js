// 格式化日期：将 2025-12-09T10:00:00Z 转为 2025-12-09
export const formatDate = (dateStr) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  }).replace(/\//g, '-'); // 确保显示为 YYYY-MM-DD
};

// 格式化日期时间
export const formatDateTime = (dateStr) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN', {
    year: 'numeric', 
    month: '2-digit', 
    day: '2-digit',
    hour: '2-digit', 
    minute: '2-digit'
  });
};

// 状态翻译与颜色映射
export const getStatusConfig = (status, isOverdue = false) => {
  if (isOverdue) return { label: '已逾期', class: 'badge-danger' };
  
  const map = {
    'borrowed': { label: '借阅中', class: 'badge-primary' },
    'returned': { label: '已归还', class: 'badge-success' },
    'waiting': { label: '预约等待', class: 'badge-warning' },
    'canceled': { label: '已取消', class: 'badge-secondary' },
    'active': { label: '正常', class: 'badge-success' },
    'disabled': { label: '禁用', class: 'badge-danger' },
    'admin': { label: '管理员', class: 'badge-dark' },
    'user': { label: '普通用户', class: 'badge-info' }
  };
  return map[status] || { label: status, class: 'badge-secondary' };
};