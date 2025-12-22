import { createRouter, createWebHistory } from 'vue-router';
import { getToken, getUser } from '../utils/auth';
import Login from '../views/Login.vue';
import Layout from '../views/Layout.vue';
import Dashboard from '../views/Dashboard.vue';
import Books from '../views/Books.vue';
import Borrowing from '../views/Borrowing.vue';
import Profile from '../views/Profile.vue';
import Users from '../views/Users.vue';
import Categories from '../views/Categories.vue';

const routes = [
  { path: '/login', component: Login, meta: { guest: true } },
  {
    path: '/',
    component:  Layout,
    children: [
      { path: '', redirect: '/books' },
      { path: 'books', component: Books },
      { path: 'borrowing', component: Borrowing },
      { path: 'profile', component: Profile },
      // 管理员专属路由
      { path: 'dashboard', component: Dashboard, meta: { admin:  true } },
      { path: 'users', component: Users, meta: { admin: true } },
      { path: 'categories', component:  Categories, meta: { admin: true } }
    ]
  },
  // 404 重定向
  { path: '/:pathMatch(.*)*', redirect: '/books' }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const hasToken = getToken();
  const user = getUser();
  
  // 未登录且不是访客页面，跳转登录
  if (! hasToken && ! to.meta.guest) {
    next('/login');
    return;
  }
  
  // 已登录访问登录页，跳转首页
  if (hasToken && to.meta.guest) {
    next('/');
    return;
  }
  
  // 管理员页面权限检查
  if (to.meta.admin && user?.role !== 'admin') {
    next('/books');
    return;
  }
  
  next();
});

export default router;