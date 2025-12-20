<template>
  <div class="card">
    <h2>个人中心</h2>
    <div v-if="info">
      <p><strong>用户名:</strong> {{ info.username }}</p>
      <p><strong>邮箱:</strong> {{ info.email }}</p>
      <p><strong>角色:</strong> {{ info.role }}</p>
      <p><strong>借阅上限:</strong> {{ info.borrow_limit }} 本</p>
    </div>
    
    <hr class="mt-4">
    
    <h3>修改密码</h3>
    <div class="form-group">
      <label>旧密码</label>
      <input v-model="pwd.old" type="password" class="input">
    </div>
    <div class="form-group">
      <label>新密码</label>
      <input v-model="pwd.new" type="password" class="input">
    </div>
    <button class="btn" @click="handlePwdChange">修改密码</button>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { getUserInfo, changePassword, logout } from '../api';
import { useRouter } from 'vue-router';
import { clearAuth } from '../utils/auth';

const info = ref(null);
const pwd = reactive({ old: '', new: '' });
const router = useRouter();

onMounted(async () => {
  const res = await getUserInfo();
  info.value = res;
});

const handlePwdChange = async () => {
  await changePassword({ old_password: pwd.old, new_password: pwd.new });
  alert('密码修改成功，请重新登录');
  await logout();
  clearAuth();
  router.push('/login');
};
</script>