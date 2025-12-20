<template>
  <div class="login-page">
    <div class="card login-box">
      <div class="header">
        <h1>📚 图书管理系统</h1>
        <p>{{ isLogin ? '欢迎回来，请登录您的账户' : '注册新账户以开始使用' }}</p>
      </div>
      
      <div class="form-group">
        <label>用户名</label>
        <input v-model="form.username" type="text" placeholder="输入用户名" class="input">
      </div>
      <div class="form-group">
        <label>密码</label>
        <input v-model="form.password" type="password" placeholder="输入密码" class="input">
      </div>
      <div v-if="!isLogin" class="form-group">
        <label>邮箱</label>
        <input v-model="form.email" type="email" placeholder="输入邮箱地址" class="input">
      </div>
      
      <button class="btn w-full mt-4" style="height: 44px; font-size: 1rem;" @click="handleSubmit">
        {{ isLogin ? '立即登录' : '创建账户' }}
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

const router = useRouter();
const isLogin = ref(true);
const form = reactive({ username: '', password: '', email: '' });

const toggleMode = () => {
  isLogin.value = !isLogin.value;
  form.username = ''; form.password = ''; form.email = '';
};

const handleSubmit = async () => {
  try {
    if (isLogin.value) {
      // 1. 发起登录
      const data = await login({ username: form.username, password: form.password });
      
      // 2. 检查返回数据是否有效
      if (!data) {
        console.error("登录未返回任何数据");
        return;
      }

      // 3. 存储 Token
      if (data.access_token) {
        setToken(data.access_token);
        setRefreshToken(data.refresh_token);
      }

      // 4. 存储用户信息 (防崩坏检查)
      if (data.user) {
        setUser(data.user);
        $message.success('登录成功'); // 加上提示
        router.push('/'); // 跳转
      } else {
        $message.error('登录异常：缺少用户信息');
      }
      
    } else {
      // 注册逻辑
      await register(form);
      $message.success('注册成功，请登录');
      isLogin.value = true;
    }
  } catch (error) {
    // 错误已经被 request.js 里的 $message.error 处理过了，这里只需要打印日志
    console.error("登录流程中断:", error);
  }
};
</script>

<style scoped>
.login-page {
  display: flex; justify-content: center; align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}
.login-box { width: 400px; padding: 40px; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1), 0 10px 10px -5px rgba(0,0,0,0.04); }
.header { text-align: center; margin-bottom: 30px; }
.header h1 { color: var(--primary); font-size: 1.8rem; margin-bottom: 8px; }
.header p { color: var(--text-secondary); font-size: 0.9rem; margin: 0; }
.toggle-text { text-align: center; margin-top: 20px; font-size: 0.875rem; color: var(--text-secondary); }
.toggle-text span { color: var(--primary); cursor: pointer; font-weight: 500; }
.toggle-text span:hover { text-decoration: underline; }
</style>