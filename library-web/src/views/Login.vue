<template>
  <div class="login-page">
    <div class="login-container">
      <!-- 左侧装饰区 -->
      <div class="login-decoration">
        <div class="decoration-content">
          <div class="logo-big">📚</div>
          <h1 class="welcome-title">LibraryOS</h1>
          <p class="welcome-subtitle">现代化智能图书管理系统</p>
          <div class="feature-list">
            <div class="feature-item">✨ 高效的借阅管理</div>
            <div class="feature-item">🔒 安全的权限控制</div>
            <div class="feature-item">📊 实时数据统计</div>
            <div class="feature-item">🎯 智能预约功能</div>
          </div>
        </div>
      </div>
      
      <!-- 右侧表单区 -->
      <div class="login-form-wrapper">
        <div class="form-card">
          <div class="form-header">
            <h2 class="form-title">{{ isLogin ? '欢迎回来' : '创建账户' }}</h2>
            <p class="form-subtitle">{{ isLogin ? '登录您的账户继续使用' : '注册新账户开始使用' }}</p>
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
        class="btn-submit" 
        :disabled="loading"
        @click="handleSubmit"
      >
        <span v-if="loading">⏳ 处理中...</span>
        <span v-else>{{ isLogin ? '🚀 立即登录' : '✨ 创建账户' }}</span>
      </button>
      
      <div class="form-footer">
        <span class="toggle-link" @click="toggleMode">
          {{ isLogin ? '还没有账号？点击注册' : '已有账号？返回登录' }}
        </span>
      </div>
        </div>
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
  if (!validate()) {
    $message.warning('请检查表单填写是否正确');
    return;
  }

  loading.value = true;
  try {
    if (isLogin.value) {
      const data = await login({ username: form.username, password: form.password });
      
      if (!data) {
        $message.error('登录失败：服务器未返回数据');
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
        $message.error('登录失败：缺少用户信息');
      }
    } else {
      // 注册逻辑
      const registerData = {
        username: form.username,
        password: form.password,
        email: form.email
      };
      if (form.phone) registerData.phone = form.phone;

      await register(registerData);
      $message.success('注册成功，请登录');
      isLogin.value = true;
      form.password = '';
    }
  } catch (error) {
    console.error('操作失败:', error);
    // 错误信息已在request.js中处理，这里不需要重复提示
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: var(--gradient-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  max-width: 1100px;
  width: 100%;
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-xl);
  animation: slideIn 0.6s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 左侧装饰区 */
.login-decoration {
  background: var(--gradient-primary);
  padding: 60px 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.login-decoration::before {
  content: '';
  position: absolute;
  width: 300px;
  height: 300px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  top: -100px;
  right: -100px;
}

.decoration-content {
  position: relative;
  z-index: 1;
  color: white;
  text-align: center;
}

.logo-big {
  font-size: 5rem;
  margin-bottom: 20px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

.welcome-title {
  font-size: 3rem;
  font-weight: 800;
  margin: 0 0 10px 0;
  letter-spacing: -1px;
}

.welcome-subtitle {
  font-size: 1.125rem;
  opacity: 0.95;
  margin: 0 0 40px 0;
}

.feature-list {
  text-align: left;
  display: inline-block;
}

.feature-item {
  font-size: 1rem;
  margin: 12px 0;
  display: flex;
  align-items: center;
  gap: 12px;
  opacity: 0.9;
}

/* 右侧表单区 */
.login-form-wrapper {
  padding: 60px 50px;
  display: flex;
  align-items: center;
}

.form-card {
  width: 100%;
}

.form-header {
  margin-bottom: 32px;
}

.form-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-main);
  margin: 0 0 8px 0;
}

.form-subtitle {
  font-size: 0.95rem;
  color: var(--text-secondary);
  margin: 0;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--text-main);
  margin-bottom: 8px;
}

.required {
  color: var(--danger);
}

.input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 0.95rem;
  transition: var(--transition);
  background: white;
}

.input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.error-text {
  color: var(--danger);
  font-size: 0.8rem;
  margin-top: 6px;
  display: block;
  animation: shake 0.3s;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}

.btn-submit {
  width: 100%;
  padding: 14px;
  background: var(--gradient-primary);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition);
  margin-top: 8px;
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.3);
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4);
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.form-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.toggle-link {
  color: var(--primary);
  cursor: pointer;
  font-weight: 600;
  transition: var(--transition);
}

.toggle-link:hover {
  color: var(--primary-hover);
  text-decoration: underline;
}

/* 响应式 */
@media (max-width: 968px) {
  .login-container {
    grid-template-columns: 1fr;
  }
  
  .login-decoration {
    display: none;
  }
  
  .login-form-wrapper {
    padding: 40px 30px;
  }
}
</style>