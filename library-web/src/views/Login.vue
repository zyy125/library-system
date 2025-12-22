<template>
  <div class="login-page">
    <div class="card login-box">
      <div class="header">
        <h1>📚 图书管理系统</h1>
        <p>{{ isLogin ? '欢迎回来，请登录您的账户' : '注册新账户以开始使用' }}</p>
      </div>
      
      <div class="form-group">
        <label>用户名 <span class="required">*</span></label>
        <input v-model="form.username" type="text" placeholder="4-20字符，字母数字下划线" class="input">
        <span v-if="errors.username" class="error-text">{{ errors.username }}</span>
      </div>
      
      <div class="form-group">
        <label>密码 <span class="required">*</span></label>
        <input v-model="form.password" type="password" placeholder="8-32字符，包含字母和数字" class="input">
        <span v-if="errors.password" class="error-text">{{ errors.password }}</span>
      </div>
      
      <template v-if="! isLogin">
        <div class="form-group">
          <label>邮箱 <span class="required">*</span></label>
          <input v-model="form.email" type="email" placeholder="输入有效邮箱地址" class="input">
          <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
        </div>
        <div class="form-group">
          <label>手机号</label>
          <input v-model="form.phone" type="tel" placeholder="11位手机号（选填）" class="input">
          <span v-if="errors.phone" class="error-text">{{ errors.phone }}</span>
        </div>
      </template>
      
      <button 
        class="btn w-full mt-4" 
        style="height: 44px; font-size: 1rem;" 
        :disabled="loading"
        @click="handleSubmit"
      >
        <span v-if="loading">处理中...</span>
        <span v-else>{{ isLogin ? '立即登录' : '创建账户' }}</span>
      </button>
      
      <div class="toggle-text">
        <span @click="toggleMode">
          {{ isLogin ? '还没有账号？点击注册' : '已有账号？返回登录' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { login, register } from '../api';
import { setToken, setRefreshToken, setUser } from '../utils/auth';
import { $message } from '../utils/toast';

const router = useRouter();
const isLogin = ref(true);
const loading = ref(false);
const form = reactive({ username: '', password:  '', email: '', phone: '' });
const errors = reactive({ username: '', password:  '', email: '', phone: '' });

const toggleMode = () => {
  isLogin.value = !isLogin.value;
  Object.assign(form, { username: '', password: '', email:  '', phone: '' });
  Object.assign(errors, { username: '', password: '', email: '', phone: '' });
};

// 表单验证
const validate = () => {
  let valid = true;
  Object.assign(errors, { username: '', password: '', email: '', phone:  '' });

  // 用户名验证：4-20字符，字母数字下划线
  if (! form.username) {
    errors.username = '请输入用户名';
    valid = false;
  } else if (!/^[a-zA-Z0-9_]{4,20}$/.test(form.username)) {
    errors.username = '用户名需4-20字符，仅支持字母数字下划线';
    valid = false;
  }

  // 密码验证：8-32字符，包含字母和数字
  if (!form.password) {
    errors.password = '请输入密码';
    valid = false;
  } else if (form.password.length < 8 || form.password.length > 32) {
    errors.password = '密码需8-32字符';
    valid = false;
  } else if (!/(?=.*[a-zA-Z])(?=.*\d)/.test(form.password)) {
    errors.password = '密码需包含字母和数字';
    valid = false;
  }

  // 注册时额外验证
  if (!isLogin.value) {
    // 邮箱验证
    if (!form.email) {
      errors.email = '请输入邮箱';
      valid = false;
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
      errors.email = '邮箱格式不正确';
      valid = false;
    }

    // 手机号验证（选填）
    if (form.phone && !/^1\d{10}$/.test(form.phone)) {
      errors.phone = '手机号格式不正确';
      valid = false;
    }
  }

  return valid;
};

const handleSubmit = async () => {
  if (!validate()) return;

  loading.value = true;
  try {
    if (isLogin.value) {
      const data = await login({ username: form.username, password: form.password });
      
      if (! data) {
        $message.error('登录异常：服务器未返回数据');
        return;
      }

      if (data.access_token) {
        setToken(data.access_token);
        setRefreshToken(data.refresh_token);
      }

      if (data.user) {
        setUser(data.user);
        $message.success(`欢迎回来，${data.user.username}！`);
        router.push('/');
      } else {
        $message.error('登录异常：缺少用户信息');
      }
    } else {
      // 注册逻辑
      const registerData = {
        username: form.username,
        password: form.password,
        email:  form.email
      };
      if (form.phone) registerData.phone = form.phone;

      await register(registerData);
      $message.success('注册成功，请登录');
      isLogin.value = true;
      form.password = '';
    }
  } catch (error) {
    console.error('操作失败:', error);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-page {
  display: flex; 
  justify-content: center; 
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-box { 
  width: 420px; 
  padding: 40px; 
  box-shadow: 0 25px 50px -12px rgba(0,0,0,0.25);
  border-radius: 16px;
}
.header { text-align: center; margin-bottom: 32px; }
.header h1 { color: var(--primary); font-size: 1.8rem; margin-bottom: 8px; }
.header p { color: var(--text-secondary); font-size: 0.9rem; margin:  0; }
.toggle-text { text-align: center; margin-top: 24px; font-size: 0.875rem; color: var(--text-secondary); }
.toggle-text span { color: var(--primary); cursor: pointer; font-weight: 500; }
.toggle-text span:hover { text-decoration: underline; }
.required { color: var(--danger); }
.error-text { color: var(--danger); font-size: 0.75rem; margin-top: 4px; display: block; }
</style>