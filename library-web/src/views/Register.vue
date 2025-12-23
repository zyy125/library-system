<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-card">
        <h1>📚 图书管理系统</h1>
        <h2>用户注册</h2>
        
        <form @submit.prevent="handleRegister">
          <div class="form-group">
            <label>用户名</label>
            <input 
              v-model="form.username" 
              class="input" 
              placeholder="4-20字符，字母数字下划线"
              required
            />
          </div>

          <div class="form-group">
            <label>邮箱</label>
            <input 
              v-model="form.email" 
              type="email"
              class="input" 
              placeholder="请输入有效邮箱"
              required
            />
          </div>

          <div class="form-group">
            <label>手机号（可选）</label>
            <input 
              v-model="form.phone" 
              class="input" 
              placeholder="11位手机号"
            />
          </div>

          <div class="form-group">
            <label>密码</label>
            <input 
              v-model="form.password" 
              type="password"
              class="input" 
              placeholder="8-32字符，包含字母和数字"
              required
            />
          </div>

          <div class="form-group">
            <label>确认密码</label>
            <input 
              v-model="form.confirmPassword" 
              type="password"
              class="input" 
              placeholder="请再次输入密码"
              required
            />
          </div>

          <button type="submit" class="btn w-full" :disabled="loading">
            {{ loading ? '注册中...' : '注册' }}
          </button>
        </form>

        <div class="auth-footer">
          已有账号？<router-link to="/login">立即登录</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { register } from '../api';
import { $message } from '../utils/toast';

const router = useRouter();
const loading = ref(false);

const form = ref({
  username: '',
  email: '',
  phone: '',
  password: '',
  confirmPassword: ''
});

const handleRegister = async () => {
  if (form.value.password !== form.value.confirmPassword) {
    $message.error('两次输入的密码不一致');
    return;
  }

  loading.value = true;
  try {
    const data = {
      username: form.value.username,
      email: form.value.email,
      password: form.value.password
    };
    if (form.value.phone) {
      data.phone = form.value.phone;
    }

    await register(data);
    $message.success('注册成功，请登录');
    router.push('/login');
  } catch (error) {
    console.error('注册失败:', error);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.auth-container {
  width: 100%;
  max-width: 440px;
}

.auth-card {
  background: white;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.auth-card h1 {
  text-align: center;
  margin-bottom: 8px;
  font-size: 1.75rem;
}

.auth-card h2 {
  text-align: center;
  color: var(--text-secondary);
  font-weight: 500;
  margin-bottom: 32px;
  font-size: 1rem;
}

.auth-footer {
  text-align: center;
  margin-top: 24px;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.auth-footer a {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
  margin-left: 4px;
}

.auth-footer a:hover {
  text-decoration: underline;
}
</style>
