<script setup lang="ts">
import SideNav from './SideNav.vue';
import HeaderNav from './HeaderNav.vue';
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

defineProps<{
  user: User;
}>();

const emit = defineEmits<{
  (e: 'logout'): void;
}>();

const sidebarCollapsed = ref(false);

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value;
};


</script>

<template>
  <div class="min-h-screen bg-transparent overflow-x-hidden">
    <!-- Mobile Overlay Backdrop -->
    <div
      v-if="!sidebarCollapsed"
      @click="toggleSidebar"
      class="fixed inset-0 bg-black/50 z-40 md:hidden transition-opacity"
    ></div>

    <!-- Sidebar -->
    <SideNav :is-collapsed="sidebarCollapsed" @toggle="toggleSidebar" />

    <!-- Main Content Area -->
    <div
      :class="[
        'min-h-screen transition-all duration-300',
        'md:ml-20',
        !sidebarCollapsed && 'md:ml-64'
      ]"
    >
      <!-- Header -->
      <HeaderNav :user="user" :sidebar-collapsed="sidebarCollapsed" @logout="emit('logout')" />

      <!-- Page Content -->
      <main class="pt-24 p-4 sm:p-6 md:p-8 overflow-x-hidden">
        <slot></slot>
      </main>
    </div>
  </div>
</template>
