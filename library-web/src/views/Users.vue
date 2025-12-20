<template>
  <div>
    <div class="card flex justify-between">
      <div class="flex gap-2">
        <input v-model="params.username" class="input" placeholder="搜索用户名..." />
        <button class="btn" @click="fetchUsers">搜索</button>
      </div>
      <button class="btn" @click="openModal('create')">➕ 新增用户</button>
    </div>

    <div class="card" style="padding: 0; overflow: hidden;">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>用户名</th>
            <th>邮箱 / 手机</th>
            <th>角色</th>
            <th>状态</th>
            <th>借阅额度</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td>{{ u.id }}</td>
            <td><strong>{{ u.username }}</strong></td>
            <td>
              <div>{{ u.email }}</div>
              <div class="text-secondary text-sm">{{ u.phone || '-' }}</div>
            </td>
            <td>
              <span class="badge" :class="u.role === 'admin' ? 'badge-dark' : 'badge-info'">
                {{ u.role === 'admin' ? '管理员' : '用户' }}
              </span>
            </td>
            <td>
              <span class="badge" :class="u.status === 'active' ? 'badge-success' : 'badge-danger'">
                {{ u.status === 'active' ? '正常' : '禁用' }}
              </span>
            </td>
            <td>{{ u.borrowing_count }} / {{ u.borrow_limit }}</td>
            <td class="text-right">
              <button class="btn btn-sm btn-secondary" @click="openModal('edit', u)">编辑</button>
              <button class="btn btn-sm btn-danger" style="margin-left: 5px;" @click="handleDelete(u.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 用户弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <h3>{{ modalType === 'create' ? '新增用户' : '编辑用户' }}</h3>
        
        <div class="form-group">
          <label>用户名 <span style="color:red">*</span></label>
          <input v-model="form.username" class="input" :disabled="modalType === 'edit'">
        </div>
        
        <div class="form-group" v-if="modalType === 'create'">
          <label>密码 <span style="color:red">*</span></label>
          <input v-model="form.password" class="input" type="password">
        </div>

        <div class="form-group">
          <label>邮箱 <span style="color:red">*</span></label>
          <input v-model="form.email" class="input">
        </div>

        <div class="form-group">
          <label>手机号</label>
          <input v-model="form.phone" class="input">
        </div>

        <div class="flex gap-4">
          <div class="form-group w-full">
            <label>角色</label>
            <select v-model="form.role" class="input">
              <option value="user">普通用户</option>
              <option value="admin">管理员</option>
            </select>
          </div>
          <div class="form-group w-full">
             <label>状态</label>
             <select v-model="form.status" class="input">
               <option value="active">正常</option>
               <option value="disabled">禁用</option>
             </select>
          </div>
        </div>

        <div class="form-group">
          <label>借阅上限</label>
          <input v-model.number="form.borrow_limit" type="number" class="input">
        </div>

        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showModal = false">取消</button>
          <button class="btn" @click="submitForm">确认保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { getUserList, createUser, updateUser, deleteUser } from '../api';
import { $message } from '../utils/toast';

const users = ref([]);
const params = reactive({ page: 1, limit: 100, username: '' });
const showModal = ref(false);
const modalType = ref('create');
const currentId = ref(null);

const form = reactive({
  username: '', password: '', email: '', phone: '', role: 'user', status: 'active', borrow_limit: 5
});

const fetchUsers = async () => {
  try {
    const res = await getUserList(params);
    users.value = res.users || [];
  } catch(e) { console.error(e); }
};

const openModal = (type, user = null) => {
  modalType.value = type;
  showModal.value = true;
  if (type === 'create') {
    Object.assign(form, { username: '', password: '', email: '', phone: '', role: 'user', status: 'active', borrow_limit: 5 });
  } else {
    currentId.value = user.id;
    // 复制数据到 form
    Object.assign(form, { 
      username: user.username, 
      email: user.email, 
      phone: user.phone, 
      role: user.role, 
      status: user.status, 
      borrow_limit: user.borrow_limit 
    });
  }
};

const submitForm = async () => {
  try {
    if (modalType.value === 'create') {
      await createUser(form);
      $message.success('用户创建成功');
    } else {
      await updateUser(currentId.value, form);
      $message.success('用户更新成功');
    }
    showModal.value = false;
    fetchUsers();
  } catch (e) {
    // 错误会被 request.js 捕获，但我们也可以在这里处理
  }
};

const handleDelete = async (id) => {
  if(confirm('确定删除该用户？此操作不可恢复。')) {
    try {
      await deleteUser(id);
      $message.success('用户删除成功');
      fetchUsers();
    } catch(e) {}
  }
};

onMounted(fetchUsers);
</script>