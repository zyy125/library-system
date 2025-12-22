<template>
  <div class="profile-page">
    <div class="profile-grid">
      <!-- 左侧：用户信息卡片 -->
      <div class="card profile-card">
        <div class="profile-header">
          <div class="avatar-large">{{info?.username?.[0]?.toUpperCase() || 'U' }}</div>
          <div class="profile-info">
            <h2>{{ info?.username }}</h2>
            <span class="badge" :class="info?.role === 'admin' ?  'badge-dark' : 'badge-info'">
              {{ info?.role === 'admin' ?  '管理员' : '普通用户' }}
            </span>
            <span class="badge" :class="info?.status === 'active' ? 'badge-success' : 'badge-danger'" style="margin-left: 8px;">
              {{ info?.status === 'active' ? '正常' : '已禁用' }}
            </span>
          </div>
        </div>
        
        <div class="info-list" v-if="info">
          <div class="info-item">
            <span class="label">📧 邮箱</span>
            <span class="value">{{ info.email }}</span>
          </div>
          <div class="info-item">
            <span class="label">📱 手机</span>
            <span class="value">{{ info.phone || '未绑定' }}</span>
          </div>
          <div class="info-item">
            <span class="label">📚 借阅额度</span>
            <span class="value">
              <strong>{{ info.borrowing_count || 0 }}</strong> / {{ info.borrow_limit }} 本
            </span>
          </div>
          <div class="info-item">
            <span class="label">⚠️ 逾期图书</span>
            <span class="value" :class="info.overdue_count > 0 ?  'text-danger' : ''">
              {{ info.overdue_count || 0 }} 本
            </span>
          </div>
          <div class="info-item">
            <span class="label">📅 注册时间</span>
            <span class="value">{{ formatDate(info.created_at) }}</span>
          </div>
        </div>
        
        <button class="btn w-full mt-4" @click="showEditModal = true">
          ✏️ 编辑个人信息
        </button>
      </div>

      <!-- 右侧：修改密码 -->
      <div class="card">
        <h3>🔐 修改密码</h3>
        <p class="text-secondary" style="font-size: 0.875rem; margin-bottom: 20px;">
          修改密码后需要重新登录
        </p>
        
        <div class="form-group">
          <label>当前密码 <span class="required">*</span></label>
          <input v-model="pwd.old" type="password" class="input" placeholder="输入当前密码">
        </div>
        <div class="form-group">
          <label>新密码 <span class="required">*</span></label>
          <input v-model="pwd.new" type="password" class="input" placeholder="8-32字符，包含字母和数字">
        </div>
        <div class="form-group">
          <label>确认新密码 <span class="required">*</span></label>
          <input v-model="pwd.confirm" type="password" class="input" placeholder="再次输入新密码">
        </div>
        <button class="btn btn-danger" :disabled="pwdLoading" @click="handlePwdChange">
          {{ pwdLoading ? '处理中...' : '确认修改' }}
        </button>
      </div>
    </div>

    <!-- 编辑个人信息弹窗 -->
    <div v-if="showEditModal" class="modal-overlay" @click.self="showEditModal = false">
      <div class="modal">
        <h3>✏️ 编辑个人信息</h3>
        
        <div class="form-group">
          <label>用户名</label>
          <input v-model="editForm.username" class="input">
        </div>
        <div class="form-group">
          <label>邮箱</label>
          <input v-model="editForm.email" type="email" class="input">
        </div>
        <div class="form-group">
          <label>手机号</label>
          <input v-model="editForm.phone" type="tel" class="input" placeholder="11位手机号">
        </div>
        
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showEditModal = false">取消</button>
          <button class="btn" :disabled="editLoading" @click="handleUpdateInfo">
            {{ editLoading ? '保存中...' :  '保存修改' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { getUserInfo, updateUserInfo, changePassword, logout } from '../api';
import { useRouter } from 'vue-router';
import { clearAuth, setUser } from '../utils/auth';
import { formatDate } from '../utils/format';
import { $message } from '../utils/toast';

const router = useRouter();
const info = ref(null);
const pwd = reactive({ old: '', new: '', confirm: '' });
const pwdLoading = ref(false);
const showEditModal = ref(false);
const editForm = reactive({ username: '', email:  '', phone: '' });
const editLoading = ref(false);

const loadUserInfo = async () => {
  try {
    const res = await getUserInfo();
    info.value = res;
    // 同步更新本地存储
    setUser(res);
  } catch (e) {
    $message.error('获取用户信息失败');
  }
};

onMounted(loadUserInfo);

const handlePwdChange = async () => {
  if (! pwd.old || !pwd.new) {
    $message.warning('请填写完整密码信息');
    return;
  }
  
  if (pwd.new !== pwd.confirm) {
    $message.error('两次输入的新密码不一致');
    return;
  }

  if (pwd.new.length < 8 || pwd.new.length > 32) {
    $message.error('新密码需8-32字符');
    return;
  }

  if (!/(?=.*[a-zA-Z])(?=.*\d)/.test(pwd.new)) {
    $message.error('新密码需包含字母和数字');
    return;
  }

  pwdLoading.value = true;
  try {
    await changePassword({ old_password: pwd.old, new_password:  pwd.new });
    $message.success('密码修改成功，请重新登录');
    await logout();
    clearAuth();
    router.push('/login');
  } catch (e) {
    console.error(e);
  } finally {
    pwdLoading.value = false;
  }
};

const handleUpdateInfo = async () => {
  // 构建更新数据，只包含有变化的字段
  const updateData = {};
  if (editForm.username && editForm.username !== info.value.username) {
    updateData.username = editForm.username;
  }
  if (editForm.email && editForm.email !== info.value.email) {
    updateData.email = editForm.email;
  }
  if (editForm.phone !== info.value.phone) {
    updateData.phone = editForm.phone;
  }

  if (Object.keys(updateData).length === 0) {
    $message.warning('没有修改任何信息');
    return;
  }

  editLoading.value = true;
  try {
    await updateUserInfo(updateData);
    $message.success('个人信息更新成功');
    showEditModal.value = false;
    loadUserInfo(); // 重新加载
  } catch (e) {
    console.error(e);
  } finally {
    editLoading.value = false;
  }
};

// 打开编辑弹窗时填充当前数据
const openEditModal = () => {
  editForm.username = info.value?.username || '';
  editForm.email = info.value?.email || '';
  editForm.phone = info.value?.phone || '';
  showEditModal.value = true;
};
</script>

<style scoped>
.profile-page { max-width: 900px; }
.profile-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; }
@media (max-width: 768px) {
  .profile-grid { grid-template-columns: 1fr; }
}

.profile-card { text-align: center; }
.profile-header { margin-bottom: 24px; }
.avatar-large {
  width: 80px; height: 80px;
  background:  linear-gradient(135deg, var(--primary) 0%, #7c3aed 100%);
  color: white; border-radius: 50%;
  display:  flex; align-items: center; justify-content: center;
  font-size:  2rem; font-weight: bold;
  margin:  0 auto 16px;
  box-shadow: 0 4px 14px rgba(79, 70, 229, 0.4);
}
.profile-info h2 { margin: 0 0 8px; }

.info-list { text-align: left; }
.info-item {
  display: flex; justify-content: space-between;
  padding: 12px 0; border-bottom: 1px solid var(--border);
}
.info-item:last-child { border-bottom: none; }
.info-item .label { color: var(--text-secondary); }
.info-item .value { font-weight: 500; }

.text-danger { color: var(--danger); font-weight: 600; }
.text-secondary { color: var(--text-secondary); }
.required { color: var(--danger); }
</style>