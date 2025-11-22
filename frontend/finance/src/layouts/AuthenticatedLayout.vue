<script setup lang="ts">
import { ref, onMounted, provide } from 'vue';
import { useRouter } from 'vue-router';
import DashboardLayout from '../components/layout/DashboardLayout.vue';
import { api } from '../services/api';

const router = useRouter();
const user = ref<any>(null);
const loading = ref(true);

// Provide user to all child components
provide('user', user);

onMounted(async () => {
  // Check if we have user data from recent login
  const tempUser = sessionStorage.getItem('temp_user');
  if (tempUser) {
    try {
      user.value = JSON.parse(tempUser);
      sessionStorage.removeItem('temp_user');
      loading.value = false;
      return;
    } catch (e) {
      // Invalid data, continue with validation
    }
  }

  // Validate token if no temp user data
  try {
    const response = await api.validateToken();
    if (response.success) {
      user.value = response.user;
    } else {
      router.push('/login');
    }
  } catch (error) {
    api.clearToken();
    router.push('/login');
  } finally {
    loading.value = false;
  }
});

const handleLogout = () => {
  api.clearToken();
  router.push('/login');
};
</script>

<template>
  <DashboardLayout v-if="user && !loading" :user="user" @logout="handleLogout">
    <router-view />
  </DashboardLayout>
  <div v-else-if="loading" class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-blue-900 to-indigo-900">
    <div class="text-center">
      <div class="relative w-20 h-20 mx-auto mb-6">
        <div class="absolute inset-0 border-4 border-blue-500/30 rounded-full"></div>
        <div class="absolute inset-0 border-4 border-t-blue-500 border-r-transparent border-b-transparent border-l-transparent rounded-full animate-spin"></div>
        <div class="absolute inset-2 border-4 border-t-emerald-500 border-r-transparent border-b-transparent border-l-transparent rounded-full animate-spin animation-delay-150" style="animation-direction: reverse;"></div>
      </div>
      <p class="text-blue-200 text-lg font-medium">Loading your dashboard...</p>
      <p class="text-blue-400 text-sm mt-2">Please wait</p>
    </div>
  </div>
</template>
