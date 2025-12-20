import { createRouter, createWebHistory } from 'vue-router';
import { getToken } from '../utils/auth';
import Login from '../views/Login.vue';
import Layout from '../views/Layout.vue';
import Dashboard from '../views/Dashboard.vue';
import Books from '../views/Books.vue';
import Borrowing from '../views/Borrowing.vue';
import Profile from '../views/Profile.vue';
import Users from '../views/Users.vue';

const routes = [
  { path: '/login', component: Login },
  {
    path: '/',
    component: Layout,
    children: [
      { path: '', redirect: '/books' },
      { path: 'dashboard', component: Dashboard },
       { path: 'users', component: Users },
      { path: 'books', component: Books },
      { path: 'borrowing', component: Borrowing },
      { path: 'profile', component: Profile }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const hasToken = getToken();
  if (!hasToken && to.path !== '/login') {
    next('/login');
  } else {
    next();
  }
});

export default router;