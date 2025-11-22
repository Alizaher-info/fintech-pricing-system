import { createRouter, createWebHistory } from 'vue-router';
import { api } from '../services/api';
import AuthenticatedLayout from '../layouts/AuthenticatedLayout.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/LoginView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../views/RegisterView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: AuthenticatedLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('../views/DashboardView.vue')
        },
        {
          path: 'transactions',
          name: 'Transactions',
          component: () => import('../views/TransactionsView.vue')
        },
        {
          path: 'pricing',
          name: 'Pricing',
          component: () => import('../views/PricingView.vue')
        },
        {
          path: 'analytics',
          name: 'Analytics',
          component: () => import('../views/AnalyticsView.vue')
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('../views/UsersView.vue')
        },
        {
          path: 'settings',
          name: 'Settings',
          component: () => import('../views/SettingsView.vue')
        }
      ]
    }
  ]
});

// Navigation guard for authentication
router.beforeEach(async (to, from, next) => {
  const requiresAuth = to.meta.requiresAuth;
  const token = api.getToken();

  if (requiresAuth && !token) {
    // Protected route, no token - redirect to login
    next('/login');
  } else if (!requiresAuth && token) {
    // Login/register page, but already logged in - redirect to dashboard
    next('/dashboard');
  } else {
    // All good, proceed
    next();
  }
});

export default router;
