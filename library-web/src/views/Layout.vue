<template>
  <div class="layout">
    <header class="navbar">
      <div class="container flex justify-between h-full">
        <div class="flex gap-4">
          <h1 class="logo">LibraryOS</h1>
          <nav class="nav-items">
            <router-link to="/books" class="nav-link" active-class="active">图书大厅</router-link>
            <router-link to="/borrowing" class="nav-link" active-class="active">我的借阅</router-link>
            <router-link v-if="isAdmin" to="/dashboard" class="nav-link" active-class="active">数据看板</router-link>
            <router-link v-if="isAdmin" to="/users" class="nav-link" active-class="active">用户管理</router-link>
          </nav>
        </div>
        
        <div class="flex gap-4">
          <router-link to="/profile" class="user-profile flex gap-2">
            <div class="avatar">{{ username[0]?.toUpperCase() }}</div>
            <span>{{ username }}</span>
          </router-link>
          <button class="btn btn-secondary btn-sm" @click="handleLogout">退出</button>
        </div>
      </div>
    </header>
    <main class="main-content container">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { logout } from '../api';
import { clearAuth, getUser } from '../utils/auth';

const router = useRouter();
const user = getUser();
const isAdmin = computed(() => user && user.role === 'admin');
const username = computed(() => user ? user.username : 'User');

const handleLogout = async () => {
  await logout();
  clearAuth();
  router.push('/login');
};
</script>

<style scoped>
.layout { min-height: 100vh; display: flex; flex-direction: column; }
.navbar {
  background: white; border-bottom: 1px solid var(--border);
  height: var(--header-height); position: sticky; top: 0; z-index: 50;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}
.h-full { height: 100%; }
.logo { font-size: 1.25rem; color: var(--primary); margin: 0; margin-right: 40px; display: flex; align-items: center; }
.nav-items { display: flex; gap: 8px; height: 100%; }
.nav-link {
  text-decoration: none; color: var(--text-secondary); padding: 0 16px;
  display: flex; align-items: center; font-weight: 500; font-size: 0.9rem;
  border-bottom: 2px solid transparent; transition: all 0.2s;
  height: 100%;
}
.nav-link:hover { color: var(--primary); background-color: #f8fafc; }
.nav-link.active { color: var(--primary); border-bottom-color: var(--primary); }
.user-profile { text-decoration: none; color: var(--text-main); font-weight: 500; display: flex; align-items: center; }
.user-profile:hover { opacity: 0.8; }
.avatar {
  width: 32px; height: 32px; background: var(--primary); color: white;
  border-radius: 50%; display: flex; align-items: center; justify-content: center;
  font-size: 0.875rem; font-weight: bold;
}
.main-content { flex: 1; padding-top: 32px; }
</style>