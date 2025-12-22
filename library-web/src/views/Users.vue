<template>
  <div>
    <div class="card toolbar">
      <div class="flex gap-2">
        <input v-model="params.username" class="input" placeholder="搜索用户名..." @keyup.enter="fetchUsers" />
        <select v-model="params.role" class="input" @change="fetchUsers">
          <option value="">全部角色</option>
          <option value="admin">管理员</option>
          <option value="user">普通用户</option>
        </select>
        <select v-model="params.status" class="input" @change="fetchUsers">
          <option value="">全部状态</option>
          <option value="active">正常</option>
          <option value="disabled">禁用</option>
        </select>
        <button class="btn" @click="fetchUsers">🔍 搜索</button>
      </div>
      <button class="btn" @click="openModal('create')">➕ 新增用户</button>
    </div>

    <div class="card" style="padding: 0; overflow: hidden;">
      <div v-if="loading" class="loading-state">加载中... </div>
      
      <table v-else-if="users.length > 0">
        <thead>
          <tr>
            <th width="60">ID</th>
            <th>用户名</th>
            <th>邮箱 / 手机</th>
            <th>角色</th>
            <th>状态</th>
            <th>借阅情况</th>
            <th>注册时间</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td class="text-secondary">{{ u.id }}</td>
            <td><strong>{{ u.username }}</strong></td>
            <td>
              <div>{{ u.email }}</div>
              <div class="text-secondary text-sm">{{ u.phone || '-' }}</div>
            </td>
            <td>
              <span class="badge" :class="u.role === 'admin' ? 'badge-dark' : 'badge-info'">
                {{ u.role === 'admin' ?  '管理员' : '用户' }}
              </span>
            </td>
            <td>
              <span class="badge" :class="u.status === 'active' ? 'badge-success' : 'badge-danger'">
                {{ u.status === 'active' ? '正常' : '禁用' }}
              </span>
            </td>
            <td>
              <div>借阅:  {{ u.borrowing_count || 0 }} / {{ u.borrow_limit }}</div>
              <div v-if="u.overdue_count > 0" class="text-danger text-sm">
                逾期: {{ u.overdue_count }} 本
              </div>
            </td>
            <td class="text-secondary text-sm">{{ formatDate(u.created_at) }}</td>
            <td class="text-right">
              <div class="flex justify-end gap-2">
                <button class="btn btn-sm btn-secondary" @click="openModal('edit', u)">编辑</button>
                <button 
                  class="btn btn-sm btn-danger" 
                  :disabled="u.borrowing_count > 0"
                  :title="u.borrowing_count > 0 ? '该用户有未归还图书' : '删除用户'"
                  @click="handleDelete(u)"
                >
                  删除
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-else class="empty-state">暂无用户数据</div>

      <!-- 分页 -->
      <div class="pagination" v-if="users.length > 0">
        <div class="page-info">
          共 <strong>{{ total }}</strong> 位用户，第 {{ params.page }} / {{ totalPages }} 页
        </div>
        <div class="page-controls">
          <button class="btn btn-secondary btn-sm" :disabled="params.page <= 1" @click="changePage(-1)">上一页</button>
          <button class="btn btn-secondary btn-sm" :disabled="params.page >= totalPages" @click="changePage(1)">下一页</button>
        </div>
      </div>
    </div>

    <!-- 用户弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <h3>{{ modalType === 'create' ? '➕ 新增用户' : '✏️ 编辑用户' }}</h3>
        
        <div class="form-group">
          <label>用户名 <span class="required">*</span></label>
          <input v-model="form.username" class="input" :disabled="modalType === 'edit'">
          <span v-if="modalType === 'edit'" class="hint">用户名不可修改</span>
        </div>
        
        <div class="form-group" v-if="modalType === 'create'">
          <label>密码 <span class="required">*</span></label>
          <input v-model="form.password" class="input" type="password" placeholder="8-32字符">
        </div>

        <div class="form-group">
          <label>邮箱 <span class="required">*</span></label>
          <input v-model="form.email" class="input" type="email">
        </div>

        <div class="form-group">
          <label>手机号</label>
          <input v-model="form.phone" class="input" placeholder="11位手机号（选填）">
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>角色</label>
            <select v-model="form.role" class="input">
              <option value="user">普通用户</option>
              <option value="admin">管理员</option>
            </select>
          </div>
          <div class="form-group">
             <label>状态</label>
             <select v-model="form.status" class="input">
               <option value="active">正常</option>
               <option value="disabled">禁用</option>
             </select>
          </div>
        </div>

        <div class="form-group">
          <label>借阅上限</label>
          <input v-model.number="form.borrow_limit" type="number" min="0" max="20" class="input">
        </div>

        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showModal = false">取消</button>
          <button class="btn" :disabled="submitting" @click="submitForm">
            {{ submitting ? '保存中...' :  '确认保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { getUserList, createUser, updateUser, deleteUser } from '../api';
import { formatDate } from '../utils/format';
import { $message } from '../utils/toast';

const users = ref([]);
const loading = ref(false);
const submitting = ref(false);
const total = ref(0);
const totalPages = ref(1);

const params = reactive({ page: 1, limit: 10, username: '', role: '', status: '' });
const showModal = ref(false);
const modalType = ref('create');
const currentId = ref(null);

const form = reactive({
  username:  '', password: '', email: '', phone: '', role:  'user', status: 'active', borrow_limit:  5
});

const fetchUsers = async () => {
  loading.value = true;
  try {
    const query = { page: params.page, limit: params.limit };
    if (params.username) query.username = params.username;
    if (params.role) query.role = params.role;
    if (params.status) query.status = params.status;

    const res = await getUserList(query);
    users.value = res.users || [];
    total.value = res.total || 0;
    totalPages.value = res.total_pages || 1;
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
};

const changePage = (delta) => {
  params.page += delta;
  fetchUsers();
};

const openModal = (type, user = null) => {
  modalType.value = type;
  showModal.value = true;
  if (type === 'create') {
    Object.assign(form, { username: '', password: '', email:  '', phone: '', role: 'user', status: 'active', borrow_limit: 5 });
  } else {
    currentId.value = user.id;
    Object.assign(form, {
      username: user.username,
      email: user.email,
      phone: user.phone || '',
      role:  user.role,
      status: user.status,
      borrow_limit: user.borrow_limit
    });
  }
};

const submitForm = async () => {
  // 基础验证
  if (!form.username || !form.email) {
    $message.warning('请填写必填字段');
    return;
  }
  if (modalType.value === 'create' && ! form.password) {
    $message.warning('请输入密码');
    return;
  }

  submitting.value = true;
  try {
    if (modalType.value === 'create') {
      await createUser(form);
      $message.success('用户创建成功');
    } else {
      // 编辑时构建更新数据
      const updateData = {
        email: form.email,
        phone: form.phone,
        role: form.role,
        status:  form.status,
        borrow_limit: form.borrow_limit
      };
      await updateUser(currentId.value, updateData);
      $message.success('用户更新成功');
    }
    showModal.value = false;
    fetchUsers();
  } catch (e) {
    console.error(e);
  } finally {
    submitting.value = false;
  }
};

const handleDelete = async (user) => {
  if (user.borrowing_count > 0) {
    $message.error(`无法删除：该用户有 ${user.borrowing_count} 本未归还图书`);
    return;
  }
  
  if (confirm(`确定删除用户「${user.username}」？此操作不可恢复。`)) {
    try {
      await deleteUser(user.id);
      $message.success('用户删除成功');
      fetchUsers();
    } catch (e) {}
  }
};

onMounted(fetchUsers);
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
}
.toolbar .input { width: 150px; }

.text-sm { font-size: 0.8rem; }
.text-secondary { color: var(--text-secondary); }
.text-danger { color: var(--danger); }

.loading-state, .empty-state {
  padding: 60px;
  text-align: center;
  color:  var(--text-secondary);
}

.pagination {
  padding:  16px;
  border-top:  1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.page-info { color: var(--text-secondary); font-size: 0.875rem; }
.page-controls { display: flex; gap: 8px; }

.form-row { display:  grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.required { color: var(--danger); }
.hint { font-size: 0.75rem; color: var(--text-secondary); margin-top: 4px; display: block; }
</style>