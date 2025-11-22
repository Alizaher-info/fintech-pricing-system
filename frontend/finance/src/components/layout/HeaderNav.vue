<script setup lang="ts">
import { ref } from 'vue';

interface User {
  user_id?: number;
  id?: number;
  firstName?: string;
  lastName?: string;
  email: string;
  role?: string;
  roles?: string[];
}

const props = defineProps<{
  user: User;
  sidebarCollapsed?: boolean;
}>();

const emit = defineEmits<{
  (e: 'logout'): void;
}>();

const showUserMenu = ref(false);
const showNotifications = ref(false);

const notifications = ref([
  { id: 1, title: 'New transaction', message: 'Payment of $1,234.56 received', time: '5 min ago', unread: true },
  { id: 2, title: 'System update', message: 'Security patch applied successfully', time: '1 hour ago', unread: true },
  { id: 3, title: 'Price alert', message: 'BTC reached $65,000', time: '2 hours ago', unread: false },
]);

const unreadCount = ref(notifications.value.filter(n => n.unread).length);

const toggleUserMenu = () => {
  showUserMenu.value = !showUserMenu.value;
  showNotifications.value = false;
};

const toggleNotifications = () => {
  showNotifications.value = !showNotifications.value;
  showUserMenu.value = false;
};

const markAllAsRead = () => {
  notifications.value.forEach(n => n.unread = false);
  unreadCount.value = 0;
};
</script>

<template>
  <header 
    :class="[
      'fixed top-0 right-0 bg-slate-900/80 backdrop-blur-md border-b border-white/10 z-30 shadow-sm transition-all duration-300',
      'left-0 md:left-20',
      !props.sidebarCollapsed && 'md:left-64'
    ]"
  >
    <div class="flex items-center justify-between px-4 sm:px-6 h-20">
      
      <!-- Left Section - Search -->
      <div class="flex-1 max-w-2xl">
        <div class="relative group">
          <svg class="absolute left-4 top-1/2 transform -translate-y-1/2 w-5 h-5 text-slate-400 group-focus-within:text-blue-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
          </svg>
          <input
            type="text"
            placeholder="Search transactions, users, analytics..."
            class="w-full pl-11 pr-4 py-2.5 bg-white/5 border border-white/10 rounded-full text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500/50 focus:bg-white/10 transition-all duration-300 text-sm shadow-sm backdrop-blur-sm"
          />
        </div>
      </div>

      <!-- Right Section - Actions & User -->
      <div class="flex items-center space-x-6 ml-6">
        
        <!-- Quick Actions -->
        <button
          class="group w-10 h-10 flex items-center justify-center rounded-full bg-gradient-to-br from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white transition-all duration-300 shadow-md hover:shadow-xl hover:scale-110 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 focus:ring-offset-slate-900"
          title="New Transaction"
        >
          <svg class="w-5 h-5 transition-transform duration-300 group-hover:rotate-90" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
        </button>

        <!-- Notifications -->
        <div class="relative">
          <button
            @click="toggleNotifications"
            class="group relative w-10 h-10 flex items-center justify-center rounded-full bg-gradient-to-br from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white transition-all duration-300 shadow-md hover:shadow-xl hover:scale-110 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 focus:ring-offset-slate-900"
            title="Notifications"
          >
            <svg class="w-5 h-5 transition-transform duration-300 group-hover:rotate-12 group-hover:-translate-y-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
            </svg>
            <span
              v-if="unreadCount > 0"
              class="absolute -top-1 -right-1 w-5 h-5 bg-gradient-to-br from-red-500 to-red-600 rounded-full flex items-center justify-center text-[10px] font-bold border-2 border-slate-900 shadow-lg animate-pulse"
            >
              {{ unreadCount }}
            </span>
          </button>

          <!-- Notifications Dropdown -->
          <div
            v-if="showNotifications"
            class="absolute right-0 mt-4 w-80 bg-[#0B1120]/95 backdrop-blur-2xl rounded-2xl shadow-2xl border border-white/5 overflow-hidden transform origin-top-right transition-all duration-200 z-50 ring-1 ring-black/50"
          >
            <div class="flex items-center justify-between px-5 py-4 border-b border-white/5 bg-white/[0.02]">
              <h3 class="font-semibold text-white text-sm tracking-wide">Notifications</h3>
              <button
                v-if="unreadCount > 0"
                @click="markAllAsRead"
                class="text-xs font-medium text-blue-400 hover:text-blue-300 transition-colors"
              >
                Mark all read
              </button>
            </div>
            <div class="max-h-[350px] overflow-y-auto custom-scrollbar">
              <div
                v-for="notification in notifications"
                :key="notification.id"
                :class="[
                  'px-5 py-4 border-b border-white/5 hover:bg-white/[0.02] transition-colors cursor-pointer group',
                  notification.unread ? 'bg-blue-500/[0.03]' : ''
                ]"
              >
                <div class="flex justify-between items-start mb-1.5">
                  <h4 :class="['text-sm font-medium', notification.unread ? 'text-blue-100' : 'text-slate-400 group-hover:text-slate-300']">
                    {{ notification.title }}
                  </h4>
                  <span class="text-[10px] text-slate-600 whitespace-nowrap ml-2">{{ notification.time }}</span>
                </div>
                <p class="text-xs text-slate-500 leading-relaxed group-hover:text-slate-400 transition-colors">{{ notification.message }}</p>
              </div>
            </div>
            <div class="px-5 py-3 bg-white/[0.02] border-t border-white/5 text-center">
              <button class="text-xs font-medium text-slate-400 hover:text-white transition-colors">
                View all history
              </button>
            </div>
          </div>
        </div>

        <!-- User Menu -->
        <div class="relative">
          <button
            @click="toggleUserMenu"
            class="group w-10 h-10 flex items-center justify-center rounded-full bg-gradient-to-br from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white transition-all duration-300 shadow-md hover:shadow-xl hover:scale-110 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 focus:ring-offset-slate-900"
            title="User Profile"
          >
            <svg class="w-5 h-5 transition-transform duration-300 group-hover:scale-110" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </button>

          <!-- User Dropdown -->
          <div
            v-if="showUserMenu"
            class="absolute right-0 mt-4 w-64 bg-[#0B1120]/95 backdrop-blur-2xl rounded-2xl shadow-2xl border border-white/5 overflow-hidden transform origin-top-right transition-all duration-200 z-50 ring-1 ring-black/50"
          >
            <div class="px-6 py-5 border-b border-white/5 bg-white/[0.02]">
              <p class="text-sm font-semibold text-white tracking-wide">{{ user.firstName && user.lastName ? `${user.firstName} ${user.lastName}` : user.email.split('@')[0] }}</p>
              <p class="text-xs text-slate-500 mt-1 truncate font-mono">{{ user.email }}</p>
            </div>
            
            <div class="py-2">
              <a href="#" class="flex items-center px-6 py-3 text-sm text-slate-400 hover:bg-white/[0.02] hover:text-blue-400 transition-colors group">
                <svg class="w-4 h-4 mr-3 text-slate-600 group-hover:text-blue-500 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                </svg>
                Your Profile
              </a>
              <a href="#" class="flex items-center px-6 py-3 text-sm text-slate-400 hover:bg-white/[0.02] hover:text-blue-400 transition-colors group">
                <svg class="w-4 h-4 mr-3 text-slate-600 group-hover:text-blue-500 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                </svg>
                Settings
              </a>
            </div>

            <div class="border-t border-white/5 py-2">
              <button
                @click="emit('logout')"
                class="flex w-full items-center px-6 py-3 text-sm text-red-400 hover:bg-red-500/10 transition-colors group"
              >
                <svg class="w-4 h-4 mr-3 text-red-400/50 group-hover:text-red-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
                </svg>
                Sign out
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
/* Close dropdowns when clicking outside */
</style>
