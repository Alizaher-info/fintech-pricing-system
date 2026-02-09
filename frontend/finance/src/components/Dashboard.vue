<script setup lang="ts">
import { ref } from 'vue';
import LivePrices from './LivePrices.vue';

interface User {
  user_id?: number;
  id?: number;
  email: string;
  firstName?: string;
  lastName?: string;
  role?: string;
  roles?: string[];
}

const props = defineProps<{
  user: User;
}>();

// Helper to get user display name
const getUserName = () => {
  if (props.user.firstName && props.user.lastName) {
    return `${props.user.firstName} ${props.user.lastName}`;
  }
  return props.user.email.split('@')[0];
};

// Helper to get user role
const getUserRole = () => {
  if (props.user.roles && props.user.roles.length > 0) {
    return props.user.roles.join(', ');
  }
  if (props.user.role) {
    return props.user.role.replace('ROLE_', '');
  }
  return 'User';
};

// Helper to get user ID
const getUserId = () => {
  return props.user.id || props.user.user_id || 0;
};

// Dashboard Stats
const stats = ref([
  {
    title: 'Total Revenue',
    value: '$124,586',
    change: '+12.5%',
    trend: 'up',
    icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    color: 'blue'
  },
  {
    title: 'Transactions',
    value: '1,247',
    change: '+8.2%',
    trend: 'up',
    icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2',
    color: 'green'
  },
  {
    title: 'Active Users',
    value: '892',
    change: '+3.1%',
    trend: 'up',
    icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z',
    color: 'purple'
  },
  {
    title: 'Pending Quotes',
    value: '24',
    change: '-2.4%',
    trend: 'down',
    icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
    color: 'yellow'
  }
]);

const recentTransactions = ref([
  { id: 'TRX001', customer: 'John Doe', amount: '$1,234.56', status: 'completed', date: '2025-11-20' },
  { id: 'TRX002', customer: 'Jane Smith', amount: '$856.00', status: 'pending', date: '2025-11-20' },
  { id: 'TRX003', customer: 'Bob Johnson', amount: '$2,100.75', status: 'completed', date: '2025-11-19' },
  { id: 'TRX004', customer: 'Alice Brown', amount: '$450.25', status: 'failed', date: '2025-11-19' },
  { id: 'TRX005', customer: 'Charlie Wilson', amount: '$3,200.00', status: 'completed', date: '2025-11-18' }
]);
</script>

<template>
  <!-- Page Header with Gradient -->
  <div class="mb-8 relative overflow-hidden rounded-2xl bg-gradient-to-r from-slate-900 via-blue-900 to-slate-900 p-8 text-white shadow-xl">
    <div class="relative z-10">
        <h1 class="text-3xl font-bold mb-2">Dashboard Overview</h1>
        <p class="text-blue-200 text-lg">Welcome back, {{ getUserName() }}! Here's your financial portfolio status.</p>
      </div>
      <!-- Decorative background elements -->
      <div class="absolute top-0 right-0 -mt-10 -mr-10 w-64 h-64 bg-blue-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob"></div>
      <div class="absolute bottom-0 left-0 -mb-10 -ml-10 w-64 h-64 bg-purple-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000"></div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div
        v-for="stat in stats"
        :key="stat.title"
        class="bg-white/10 backdrop-blur-md rounded-xl shadow-sm border border-white/10 p-6 hover:shadow-lg transition-all duration-300 transform hover:-translate-y-1"
      >
        <div class="flex items-center justify-between mb-4">
          <div
            :class="[
              'w-12 h-12 rounded-xl flex items-center justify-center shadow-sm',
              stat.color === 'blue' ? 'bg-blue-500/20 text-blue-300' : '',
              stat.color === 'green' ? 'bg-emerald-500/20 text-emerald-300' : '',
              stat.color === 'purple' ? 'bg-purple-500/20 text-purple-300' : '',
              stat.color === 'yellow' ? 'bg-amber-500/20 text-amber-300' : ''
            ]"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="stat.icon" />
            </svg>
          </div>
          <span
            :class="[
              'text-sm font-bold px-2.5 py-0.5 rounded-full',
              stat.trend === 'up' ? 'bg-green-500/20 text-green-300' : 'bg-red-500/20 text-red-300'
            ]"
          >
            {{ stat.change }}
          </span>
        </div>
        <h3 class="text-slate-300 text-sm font-medium mb-1 uppercase tracking-wider">{{ stat.title }}</h3>
        <p class="text-2xl font-bold text-white">{{ stat.value }}</p>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      
      <!-- Recent Transactions -->
      <div class="lg:col-span-2 bg-white/10 backdrop-blur-md rounded-xl shadow-sm border border-white/10 overflow-hidden">
        <div class="p-6 border-b border-white/10 flex items-center justify-between bg-white/5">
          <h2 class="text-lg font-bold text-white">Recent Transactions</h2>
          <button class="text-sm text-blue-400 hover:text-blue-300 font-semibold hover:underline">
            View All Transactions
          </button>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="bg-white/5 border-b border-white/10">
                <th class="text-left text-xs font-semibold text-slate-300 uppercase tracking-wider py-4 px-6">Transaction ID</th>
                <th class="text-left text-xs font-semibold text-slate-300 uppercase tracking-wider py-4 px-6">Customer</th>
                <th class="text-left text-xs font-semibold text-slate-300 uppercase tracking-wider py-4 px-6">Amount</th>
                <th class="text-left text-xs font-semibold text-slate-300 uppercase tracking-wider py-4 px-6">Status</th>
                <th class="text-left text-xs font-semibold text-slate-300 uppercase tracking-wider py-4 px-6">Date</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-white/10">
              <tr v-for="transaction in recentTransactions" :key="transaction.id" class="hover:bg-white/5 transition-colors duration-150">
                <td class="py-4 px-6 text-sm font-medium text-white font-mono">{{ transaction.id }}</td>
                <td class="py-4 px-6 text-sm text-slate-300">
                  <div class="flex items-center">
                    <div class="h-8 w-8 rounded-full bg-white/20 flex items-center justify-center text-xs font-bold text-white mr-3">
                      {{ transaction.customer.charAt(0) }}
                    </div>
                    {{ transaction.customer }}
                  </div>
                </td>
                <td class="py-4 px-6 text-sm font-bold text-white">{{ transaction.amount }}</td>
                <td class="py-4 px-6">
                  <span
                    :class="[
                      'px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full border',
                      transaction.status === 'completed' ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30' : '',
                      transaction.status === 'pending' ? 'bg-amber-500/20 text-amber-300 border-amber-500/30' : '',
                      transaction.status === 'failed' ? 'bg-red-500/20 text-red-300 border-red-500/30' : ''
                    ]"
                  >
                    <span class="w-1.5 h-1.5 rounded-full mr-1.5 my-auto" :class="[
                      transaction.status === 'completed' ? 'bg-emerald-400' : '',
                      transaction.status === 'pending' ? 'bg-amber-400' : '',
                      transaction.status === 'failed' ? 'bg-red-400' : ''
                    ]"></span>
                    {{ transaction.status.charAt(0).toUpperCase() + transaction.status.slice(1) }}
                  </span>
                </td>
                <td class="py-4 px-6 text-sm text-slate-400">{{ transaction.date }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Security & System Info -->
      <div class="space-y-6">
        
        <!-- Live Crypto Prices -->
        <LivePrices />
        
        <!-- Security Features -->
        <div class="bg-white/10 backdrop-blur-md rounded-xl shadow-sm border border-white/10 p-6">
          <h3 class="text-lg font-bold text-white mb-4 flex items-center">
            <svg class="w-5 h-5 text-emerald-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            Security Status
          </h3>
          <div class="space-y-4">
            <div class="flex items-center p-3 bg-emerald-500/10 rounded-lg border border-emerald-500/20">
              <div class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center shadow-sm mr-3">
                <svg class="w-5 h-5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                </svg>
              </div>
              <div>
                <p class="text-sm font-bold text-emerald-100">JWT Authentication</p>
                <p class="text-xs text-emerald-300">Active & Secure</p>
              </div>
              <div class="ml-auto">
                <div class="w-2 h-2 bg-emerald-400 rounded-full animate-pulse"></div>
              </div>
            </div>

            <div class="flex items-center p-3 bg-blue-500/10 rounded-lg border border-blue-500/20">
              <div class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center shadow-sm mr-3">
                <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                </svg>
              </div>
              <div>
                <p class="text-sm font-bold text-blue-100">Rate Limiting</p>
                <p class="text-xs text-blue-300">Redis Enabled</p>
              </div>
              <div class="ml-auto">
                <div class="w-2 h-2 bg-blue-400 rounded-full animate-pulse"></div>
              </div>
            </div>

            <div class="flex items-center p-3 bg-purple-500/10 rounded-lg border border-purple-500/20">
              <div class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center shadow-sm mr-3">
                <svg class="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <div>
                <p class="text-sm font-bold text-purple-100">Account Protection</p>
                <p class="text-xs text-purple-300">Monitoring Active</p>
              </div>
              <div class="ml-auto">
                <div class="w-2 h-2 bg-purple-400 rounded-full animate-pulse"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- User Info Card -->
        <div class="bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl shadow-sm p-6 text-white">
          <h3 class="text-lg font-bold mb-4 flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
            </svg>
            Your Account
          </h3>
          <div class="space-y-3 text-sm">
            <div class="flex justify-between items-center py-2 border-b border-white/20">
              <span class="opacity-90">User ID</span>
              <span class="font-semibold">{{ getUserId() }}</span>
            </div>
            <div class="flex justify-between items-center py-2 border-b border-white/20">
              <span class="opacity-90">Email</span>
              <span class="font-semibold truncate ml-2">{{ user.email }}</span>
            </div>
            <div class="flex justify-between items-center py-2 border-b border-white/20">
              <span class="opacity-90">Role</span>
              <span class="font-semibold">{{ getUserRole() }}</span>
            </div>
            <div class="flex justify-between items-center py-2">
              <span class="opacity-90">Status</span>
              <span class="font-semibold flex items-center">
                <span class="w-2 h-2 bg-green-400 rounded-full mr-2 animate-pulse"></span>
                Active
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
