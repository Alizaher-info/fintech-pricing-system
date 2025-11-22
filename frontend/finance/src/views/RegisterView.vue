<script setup lang="ts">
import { useRouter } from 'vue-router';
import RegisterForm from '../components/RegisterForm.vue';

const router = useRouter();

const handleRegisterSuccess = () => {
  router.push('/dashboard');
};

const switchToLogin = () => {
  router.push('/login');
};
</script>

<template>
  <div class="min-h-screen relative overflow-hidden">
    <div class="fixed inset-0 bg-gradient-to-br from-slate-900 via-blue-900 to-indigo-900 z-0">
      <div class="absolute inset-0">
        <svg class="absolute top-10 left-10 w-32 h-20 text-green-400 opacity-20 animate-pulse" viewBox="0 0 128 80" fill="none">
          <path d="M0 40 L20 30 L40 45 L60 25 L80 35 L100 15 L128 20" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
      </div>
      <div class="absolute inset-0 bg-gradient-to-br from-black/20 via-transparent to-black/40"></div>
    </div>

    <div class="relative z-10">
      <RegisterForm @register-success="handleRegisterSuccess" @switch-to-login="switchToLogin" />
    </div>
  </div>
</template>
