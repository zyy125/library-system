<template>
  <div class="layout">
    <header class="navbar">
      <div class="container flex justify-between h-full">
        <div class="flex gap-4">
          <h1 class="logo">📚 LibraryOS</h1>
          <nav class="nav-items">
            <router-link to="/books" class="nav-link" active-class="active">图书大厅</router-link>
            <router-link to="/borrowing" class="nav-link" active-class="active">我的借阅</router-link>
            <template v-if="isAdmin">
              <router-link to="/dashboard" class="nav-link" active-class="active">数据看板</router-link>
              <router-link to="/users" class="nav-link" active-class="active">用户管理</router-link>
              <router-link to="/categories" class="nav-link" active-class="active">分类管理</router-link>
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
  } catch (e) {
    // 即使登出API失败也清除本地状态
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
}

.navbar {
  background: white;
  border-bottom: 1px solid var(--border);
  height: var(--header-height);
  position: sticky;
  top: 0;
  z-index: 50;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.h-full {
  height: 100%;
}

.logo {
  font-size: 1.25rem;
  font-weight: 700;
  color:  var(--primary);
  margin:  0;
  white-space: nowrap;
}

.nav-items {
  display:  flex;
  align-items: center;
  gap: 8px;
  margin-left: 32px;
}

.nav-link {
  text-decoration: none;
  color: var(--text-secondary);
  font-weight:  500;
  font-size: 0.875rem;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.2s;
}
.nav-link:hover {
  background: #f3f4f6;
  color: var(--text-main);
}
.nav-link.active {
  color: var(--primary);
  background: #eef2ff;
}

.user-profile {
  text-decoration: none;
  color: var(--text-main);
  display: flex;
  align-items: center;
}
.user-profile:hover {
  opacity: 0.8;
}

.avatar {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, var(--primary) 0%, #7c3aed 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size:  0.875rem;
  font-weight:  bold;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.username {
  font-weight:  500;
  font-size: 0.875rem;
  line-height: 1.2;
}

.role-badge {
  font-size: 0.7rem;
  padding: 1px 6px;
  border-radius: 4px;
}
.role-badge.admin {
  background:  #374151;
  color: #f9fafb;
}
.role-badge.user {
  background: #dbeafe;
  color: #1e40af;
}

.main-content {
  flex: 1;
  padding-top: 32px;
  padding-bottom: 32px;
}

.footer {
  background:  #f9fafb;
  border-top: 1px solid var(--border);
  padding: 16px 0;
  text-align: center;
}
.footer p {
  margin: 0;
  font-size: 0.875rem;
  color: var(--text-secondary);
}
</style>