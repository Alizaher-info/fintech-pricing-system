<script setup lang="ts">
interface MenuItem {
  name: string;
  path: string;
  icon: string;
  badge?: number;
}

defineProps<{
  isCollapsed?: boolean;
}>();

const emit = defineEmits<{
  (e: 'toggle'): void;
}>();

const menuItems: MenuItem[] = [
  { name: 'Dashboard', path: '/dashboard', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
  { name: 'Transactions', path: '/transactions', icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2', badge: 5 },
  { name: 'Pricing', path: '/pricing', icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z' },
  { name: 'Analytics', path: '/analytics', icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z' },
  { name: 'Users', path: '/users', icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' },
  { name: 'Settings', path: '/settings', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' },
];
</script>

<template>
  <aside
    :class="[
      'fixed left-0 top-0 h-screen bg-slate-900/95 backdrop-blur-xl text-white transition-all duration-300 ease-in-out border-r border-white/10 overflow-hidden',
      'z-50',
      isCollapsed ? 'w-20' : 'w-64',
      // Mobile: hide/show with transform
      isCollapsed ? '-translate-x-full md:translate-x-0' : 'translate-x-0'
    ]"
  >
    <!-- Logo & Toggle -->
    <div class="flex items-center justify-between px-3 h-20 border-b border-white/10 shrink-0 gap-2">
      <div v-if="!isCollapsed" class="flex items-center justify-center flex-1 min-w-0">
        <!-- Professional Logo -->
        <div class="relative group cursor-pointer">
          <!-- Main Logo -->
          <div class="relative flex items-center gap-1">
            <!-- Circle with Lightning Bolt -->
            <div class="relative w-10 h-10 bg-gradient-to-br from-blue-500 via-blue-600 to-cyan-500 rounded-full flex items-center justify-center shadow-lg group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">
              <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path d="M11 3a1 1 0 10-2 0v1a1 1 0 102 0V3zM15.657 5.757a1 1 0 00-1.414-1.414l-.707.707a1 1 0 001.414 1.414l.707-.707zM18 10a1 1 0 01-1 1h-1a1 1 0 110-2h1a1 1 0 011 1zM5.05 6.464A1 1 0 106.464 5.05l-.707-.707a1 1 0 00-1.414 1.414l.707.707zM5 10a1 1 0 01-1 1H3a1 1 0 110-2h1a1 1 0 011 1zM8 16v-1h4v1a2 2 0 11-4 0zM12 14c.015-.34.208-.646.477-.859a4 4 0 10-4.954 0c.27.213.462.519.476.859h4.002z"/>
              </svg>
              <!-- Animated Pulse Ring -->
              <div class="absolute inset-0 rounded-full bg-blue-400 opacity-0 group-hover:opacity-30 group-hover:scale-150 transition-all duration-500"></div>
            </div>
            
            <!-- Vertical Divider -->
            <div class="w-0.5 h-8 bg-gradient-to-b from-transparent via-blue-400 to-transparent"></div>
            
            <!-- Text Badge -->
            <div class="flex flex-col">
              <span class="text-xs font-bold text-white leading-none tracking-wider">FINTECH</span>
              <span class="text-[9px] text-blue-300 leading-none tracking-wide font-medium">PRICING</span>
            </div>
          </div>
        </div>
      </div>
      
      <button
        @click="emit('toggle')"
        class="group flex flex-col items-center justify-center gap-1.5 w-8 h-8 flex-shrink-0 transition-all duration-300"
        :class="isCollapsed ? 'mx-auto' : ''"
        title="Toggle Sidebar"
      >
        <span class="w-5 h-0.5 bg-slate-300 group-hover:bg-white rounded-full transition-all duration-300 group-hover:w-6"></span>
        <span class="w-5 h-0.5 bg-slate-300 group-hover:bg-white rounded-full transition-all duration-300"></span>
        <span class="w-5 h-0.5 bg-slate-300 group-hover:bg-white rounded-full transition-all duration-300 group-hover:w-4"></span>
      </button>
    </div>

    <!-- Navigation Menu -->
    <nav class="mt-6 px-3">
      <ul class="space-y-2">
        <li v-for="item in menuItems" :key="item.path">
          <router-link
            :to="item.path"
            :class="[
              'flex items-center px-3 py-3 rounded-lg transition-all duration-200 group relative',
              $route.path === item.path
                ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-500/50'
                : 'text-slate-300 hover:bg-white/10 hover:text-white'
            ]"
          >
            <!-- Icon -->
            <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="item.icon"/>
            </svg>

            <!-- Label -->
            <span
              v-if="!isCollapsed"
              class="ml-3 font-medium"
            >
              {{ item.name }}
            </span>

            <!-- Badge -->
            <span
              v-if="item.badge && !isCollapsed"
              class="ml-auto bg-red-500 text-white text-xs font-bold px-2 py-1 rounded-full"
            >
              {{ item.badge }}
            </span>

            <!-- Tooltip for collapsed state -->
            <div
              v-if="isCollapsed"
              class="absolute left-full ml-2 px-3 py-2 bg-slate-900 text-white text-sm rounded-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 whitespace-nowrap shadow-xl border border-slate-700 z-50"
            >
              {{ item.name }}
              <span v-if="item.badge" class="ml-2 text-red-400">({{ item.badge }})</span>
            </div>
          </router-link>
        </li>
      </ul>
    </nav>

    <!-- Bottom Section -->
    <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-white/10">
      <div
        :class="[
          'flex items-center text-slate-400 text-xs',
          isCollapsed ? 'justify-center' : 'justify-between'
        ]"
      >
        <span v-if="!isCollapsed">v1.0.0</span>
        <div class="flex space-x-2">
          <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
          <span v-if="!isCollapsed">Online</span>
        </div>
      </div>
    </div>
  </aside>
</template>

<style scoped>
/* Custom scrollbar for navigation */
nav {
  max-height: calc(100vh - 200px);
  overflow-y: auto;
}

nav::-webkit-scrollbar {
  width: 4px;
}

nav::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
}

nav::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

nav::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
