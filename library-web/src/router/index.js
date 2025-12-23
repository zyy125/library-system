import { createRouter, createWebHistory } from 'vue-router';
import { getToken, getUser } from '../utils/auth';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import Layout from '../views/Layout.vue';
import Dashboard from '../views/Dashboard.vue';
import Books from '../views/Books.vue';
import BookDetail from '../views/BookDetail.vue';
import PopularBooks from '../views/PopularBooks.vue';
import Borrowing from '../views/Borrowing.vue';
import BorrowHistory from '../views/BorrowHistory.vue';
import Reservations from '../views/Reservations.vue';
import Profile from '../views/Profile.vue';
import Users from '../views/Users.vue';
import Categories from '../views/Categories.vue';
import AdminBooks from '../views/AdminBooks.vue';
import AdminBorrows from '../views/AdminBorrows.vue';

const routes = [
  { path: '/login', component: Login, meta: { guest: true } },
  { path: '/register', component: Register, meta: { guest: true } },
  {
    path: '/',
    component: Layout,
    children: [
      { 
        path: '', 
        redirect: to => {
          const user = getUser();
          return user?.role === 'admin' ? '/dashboard' : '/books';
        }
      },
      // 普通用户路由
      { path: 'books', component: Books, name: 'Books' },
      { path: 'books/:id', component: BookDetail, name: 'BookDetail' },
      { path: 'popular', component: PopularBooks, name: 'PopularBooks' },
      { path: 'borrowing', component: Borrowing, name: 'Borrowing' },
      { path: 'borrow-history', component: BorrowHistory, name: 'BorrowHistory' },
      { path: 'reservations', component: Reservations, name: 'Reservations' },
      { path: 'profile', component: Profile, name: 'Profile' },
      // 管理员专属路由
      { path: 'dashboard', component: Dashboard, meta: { admin: true }, name: 'Dashboard' },
      { path: 'admin/users', component: Users, meta: { admin: true }, name: 'Users' },
      { path: 'admin/books', component: AdminBooks, meta: { admin: true }, name: 'AdminBooks' },
      { path: 'admin/borrows', component: AdminBorrows, meta: { admin: true }, name: 'AdminBorrows' },
      { path: 'admin/categories', component: Categories, meta: { admin: true }, name: 'Categories' }
    ]
  },
  // 404 重定向
  { path: '/:pathMatch(.*)*', redirect: to => {
    const user = getUser();
    return user?.role === 'admin' ? '/dashboard' : '/books';
  }}
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
  
  // 已登录访问登录页，根据角色跳转
  if (hasToken && to.meta.guest) {
    const user = getUser();
    next(user?.role === 'admin' ? '/dashboard' : '/books');
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