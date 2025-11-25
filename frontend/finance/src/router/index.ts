import { createRouter, createWebHistory } from 'vue-router';
import AuthenticatedLayout from '../layouts/AuthenticatedLayout.vue';
import { LoginService } from '../services/LoginService';

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
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('../views/NotFoundView.vue'),
      meta: { requiresAuth: false }
    }
  ]
});

// Get LoginService instance for authentication checks
const loginService = LoginService.getInstance();

/**
 * Navigation Guard
 * Implements authentication check before each route navigation
 */
router.beforeEach(async (to, _from, next) => {
  const requiresAuth = to.meta.requiresAuth;

  if (requiresAuth) {
    // Protected route - check if token exists first (fast)
    if (!loginService.isAuthenticated()) {
      console.log('❌ No token found, redirecting to login');
      next('/login');
      return;
    }

    // Token exists - validate session (may call backend if 5+ min passed)
    const isValid = await loginService.validateSession();
    
    if (isValid) {
      next(); // ✅ Authentication successful
    } else {
      console.log('❌ Authentication failed, redirecting to login');
      next('/login'); // ❌ Authentication failed
    }
  } else if (!requiresAuth && loginService.isAuthenticated()) {
    // Login/register page, but user already authenticated
    console.log('✅ Already authenticated, redirecting to dashboard');
    next('/dashboard');
  } else {
    // Public route, no authentication needed
    next();
  }
});

export default router;
