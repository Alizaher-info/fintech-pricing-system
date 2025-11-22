<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import LoginForm from '../components/LoginForm.vue';

const router = useRouter();
const showRegister = ref(false);

const handleLoginSuccess = (userData: any) => {
  // Store user data temporarily to avoid re-validation
  sessionStorage.setItem('temp_user', JSON.stringify(userData));
  router.push('/dashboard');
};

const switchToRegister = () => {
  router.push('/register');
};
</script>

<template>
  <div class="min-h-screen relative overflow-hidden">
    <!-- Global Animated Financial Background -->
    <div class="fixed inset-0 bg-gradient-to-br from-slate-900 via-blue-900 to-indigo-900 z-0">
      <!-- Same background animation as App.vue -->
      <div class="absolute inset-0">
        <svg class="absolute top-10 left-10 w-32 h-20 text-green-400 opacity-20 animate-pulse" viewBox="0 0 128 80" fill="none">
          <path d="M0 40 L20 30 L40 45 L60 25 L80 35 L100 15 L128 20" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
        <svg class="absolute top-32 right-20 w-40 h-24 text-red-400 opacity-20 animate-pulse" viewBox="0 0 160 96" fill="none">
          <path d="M0 60 L25 40 L50 65 L75 35 L100 55 L125 25 L160 30" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
      </div>
      <div class="absolute inset-0 bg-gradient-to-br from-black/20 via-transparent to-black/40"></div>
    </div>

    <!-- Content -->
    <div class="relative z-10">
      <LoginForm @login-success="handleLoginSuccess" @switch-to-register="switchToRegister" />
    </div>
  </div>
</template>
