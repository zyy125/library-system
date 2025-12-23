<template>
  <div class="layout">
    <header class="navbar">
      <div class="container flex justify-between h-full">
        <div class="flex gap-4">
          <h1 class="logo">📚 LibraryOS</h1>
          <nav class="nav-items">
            <!-- 普通用户导航 -->
            <template v-if="!isAdmin">
              <router-link to="/books" class="nav-link" active-class="active">📚 图书大厅</router-link>
              <router-link to="/popular" class="nav-link" active-class="active">🔥 热门排行</router-link>
              <router-link to="/borrowing" class="nav-link" active-class="active">📖 当前借阅</router-link>
              <router-link to="/borrow-history" class="nav-link" active-class="active">📜 借阅历史</router-link>
              <router-link to="/reservations" class="nav-link" active-class="active">⏰ 我的预约</router-link>
            </template>
            
            <!-- 管理员导航 -->
            <template v-if="isAdmin">
              <router-link to="/dashboard" class="nav-link" active-class="active">📊 数据看板</router-link>
              <router-link to="/admin/users" class="nav-link" active-class="active">👥 用户管理</router-link>
              <router-link to="/admin/books" class="nav-link" active-class="active">📚 图书管理</router-link>
              <router-link to="/admin/borrows" class="nav-link" active-class="active">📖 借阅管理</router-link>
              <router-link to="/admin/categories" class="nav-link" active-class="active">📋 分类管理</router-link>
            </template>
          </nav>
        </div>
        
        <div class="flex gap-4">
          <router-link to="/profile" class="user-profile flex gap-2">
            <div class="avatar">{{ username[0]?.toUpperCase() }}</div>
            <div class="user-info">
              <span class="username">{{ username }}</span>
              <span class="role-badge" :class="isAdmin ? 'admin' : 'user'">
                {{ isAdmin ? '管理员' : '用户' }}
              </span>
            </div>
          </router-link>
          <button class="btn btn-secondary btn-sm" @click="handleLogout">退出</button>
        </div>
      </div>
    </header>
    <main class="main-content container">
      <router-view />
    </main>
    <footer class="footer">
      <div class="container">
        <p>© 2025 LibraryOS 图书管理系统 | Powered by Vue 3</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { logout, getUserInfo } from '../api';
import { clearAuth, getUser, setUser } from '../utils/auth';
import { $message } from '../utils/toast';

const router = useRouter();
const user = ref(getUser());

const isAdmin = computed(() => user.value?.role === 'admin');
const username = computed(() => user.value?.username || 'User');

// 刷新用户信息
const refreshUserInfo = async () => {
  try {
    const info = await getUserInfo();
    setUser(info);
    user.value = info;
  } catch (e) {
    console.error('刷新用户信息失败');
  }
};

const handleLogout = async () => {
  try {
    await logout();
    $message.success('已安全登出');
  } catch (e) {
    // 即使登出API失败也清除本地状态
    $message.info('登出成功');
  }
  clearAuth();
  router.push('/login');
};

onMounted(() => {
  // 页面加载时刷新用户信息
  refreshUserInfo();
});
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-body);
}

.navbar {
  background: white;
  border-bottom: 1px solid var(--border-light);
  height: var(--header-height);
  position: sticky;
  top: 0;
  z-index: 50;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  backdrop-filter: blur(10px);
}

.h-full {
  height: 100%;
}

.logo {
  font-size: 1.5rem;
  font-weight: 800;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  white-space: nowrap;
  letter-spacing: -0.5px;
}

.nav-items {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: 32px;
}

.nav-link {
  text-decoration: none;
  color: var(--text-secondary);
  font-weight: 600;
  font-size: 0.9rem;
  padding: 10px 18px;
  border-radius: var(--radius-sm);
  transition: var(--transition);
  position: relative;
}

.nav-link:hover {
  background: var(--border-light);
  color: var(--text-main);
}

.nav-link.active {
  color: var(--primary);
  background: #eef2ff;
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 30px;
  height: 3px;
  background: var(--primary);
  border-radius: 3px 3px 0 0;
}

.user-profile {
  text-decoration: none;
  color: var(--text-main);
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  transition: var(--transition);
}

.user-profile:hover {
  background: var(--border-light);
}

.avatar {
  width: 42px;
  height: 42px;
  background: var(--gradient-primary);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  font-weight: bold;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.username {
  font-weight: 600;
  font-size: 0.9rem;
  line-height: 1.3;
  color: var(--text-main);
}

.role-badge {
  font-size: 0.7rem;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.role-badge.admin {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.role-badge.user {
  background: linear-gradient(135deg, #0ba360 0%, #3cba92 100%);
  color: white;
}

.main-content {
  flex: 1;
  padding-top: 40px;
  padding-bottom: 40px;
}

.footer {
  background: white;
  border-top: 1px solid var(--border-light);
  padding: 24px 0;
  text-align: center;
}

.footer p {
  margin: 0;
  font-size: 0.875rem;
  color: var(--text-secondary);
}
</style>