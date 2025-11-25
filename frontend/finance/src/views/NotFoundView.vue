<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-blue-900 to-indigo-900">
    <div class="text-center px-4">
      <div class="mb-8">
        <h1 class="text-9xl font-bold text-blue-500 mb-4">404</h1>
        <h2 class="text-3xl font-semibold text-white mb-2">Page Not Found</h2>
        <p class="text-blue-200 text-lg mb-8">
          The page you're looking for doesn't exist or has been moved.
        </p>
      </div>
      
      <div class="space-x-4">
        <button
          @click="goBack"
          class="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors"
        >
          Go Back
        </button>
        <button
          @click="goHome"
          class="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white font-medium rounded-lg transition-colors"
        >
          Go Home
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();

const goBack = () => {
  router.back();
};

const goHome = () => {
  router.push('/dashboard');
};
</script>
