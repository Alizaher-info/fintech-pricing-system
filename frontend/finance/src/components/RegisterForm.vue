<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <div class="max-w-md w-full">
      <!-- Header -->
      <div class="text-center mb-8">
        <div class="mx-auto h-16 w-16 bg-gradient-to-r from-emerald-500 via-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mb-6 shadow-2xl animate-pulse">
          <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
          </svg>
        </div>
        <h1 class="text-4xl font-bold text-white mb-3 tracking-tight">
          <span class="bg-gradient-to-r from-emerald-400 via-blue-400 to-purple-400 bg-clip-text text-transparent">
            FinTech Pro
          </span>
        </h1>
        <p class="text-blue-200 text-lg">Join Advanced Trading Platform</p>
      </div>

      <!-- Register Form -->
      <div class="bg-white/10 backdrop-blur-xl rounded-3xl shadow-2xl p-8 border border-white/20 hover:bg-white/15 transition-all duration-300">
        <h2 class="text-2xl font-bold text-center mb-6 text-white">Create Account</h2>
      
        <div v-if="error" class="mb-4 p-3 bg-red-500/20 border border-red-500/50 text-red-200 rounded-xl backdrop-blur-sm">
          {{ error }}
        </div>

        <div v-if="success" class="mb-4 p-3 bg-green-500/20 border border-green-500/50 text-green-200 rounded-xl backdrop-blur-sm">
          ✓ Registration successful! Redirecting...
        </div>

        <form @submit.prevent="handleRegister" class="space-y-5">
          <!-- First Name & Last Name -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-white mb-2">First Name</label>
              <input
                v-model="firstName"
                type="text"
                required
                class="block w-full px-4 py-3 bg-white/10 border border-white/30 rounded-xl text-white placeholder-blue-200 focus:ring-2 focus:ring-emerald-400 focus:border-emerald-400 transition-all duration-300 backdrop-blur-sm"
                placeholder="John"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-white mb-2">Last Name</label>
              <input
                v-model="lastName"
                type="text"
                required
                class="block w-full px-4 py-3 bg-white/10 border border-white/30 rounded-xl text-white placeholder-blue-200 focus:ring-2 focus:ring-emerald-400 focus:border-emerald-400 transition-all duration-300 backdrop-blur-sm"
                placeholder="Doe"
              />
            </div>
          </div>

          <!-- Email Field -->
          <div>
            <label class="block text-sm font-medium text-white mb-2">Email Address</label>
            <div class="relative">
              <input
                v-model="email"
                type="email"
                required
                class="block w-full px-4 py-3 pl-11 bg-white/10 border border-white/30 rounded-xl text-white placeholder-blue-200 focus:ring-2 focus:ring-emerald-400 focus:border-emerald-400 transition-all duration-300 backdrop-blur-sm"
                placeholder="john@example.com"
              />
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-blue-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
                </svg>
              </div>
            </div>
          </div>

          <!-- Password Field -->
          <div>
            <label class="block text-sm font-medium text-white mb-2">Password</label>
            <div class="relative">
              <input
                v-model="password"
                type="password"
                required
                minlength="6"
                class="block w-full px-4 py-3 pl-11 bg-white/10 border border-white/30 rounded-xl text-white placeholder-blue-200 focus:ring-2 focus:ring-emerald-400 focus:border-emerald-400 transition-all duration-300 backdrop-blur-sm"
                placeholder="Min. 6 characters"
              />
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-blue-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
            </div>
          </div>

          <!-- Register Button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-gradient-to-r from-emerald-500 via-blue-500 to-purple-600 text-white py-4 px-4 rounded-xl font-semibold hover:from-emerald-600 hover:via-blue-600 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-emerald-400 focus:ring-offset-2 focus:ring-offset-transparent transition-all duration-300 disabled:opacity-50 shadow-lg hover:shadow-xl transform hover:scale-105"
          >
            <span v-if="loading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Creating account...
            </span>
            <span v-else>Create Account</span>
          </button>
        </form>

        <!-- Sign In Link -->
        <p class="mt-8 text-center text-sm text-blue-200">
          Already have an account?
          <a href="#" @click.prevent="$emit('switch-to-login')" class="font-medium text-emerald-400 hover:text-emerald-300 transition-colors">
            Sign in
          </a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { RegisterService } from '../services/RegisterService';

const registerService = RegisterService.getInstance();
const emit = defineEmits(['register-success', 'switch-to-login']);

const firstName = ref('');
const lastName = ref('');
const email = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');
const success = ref(false);

const handleRegister = async () => {
  loading.value = true;
  error.value = '';
  success.value = false;

  try {
    // RegisterService - email registration only
    await registerService.register({
      email: email.value,
      password: password.value,
      firstName: firstName.value,
      lastName: lastName.value,
    });

    // Token already saved by RegisterService (auto-login)
    success.value = true;
    // Redirect after a short delay
    setTimeout(() => {
      emit('register-success');
    }, 1000);
  } catch (err: any) {
    error.value = err.message || err.error || 'Registration failed';
  } finally {
    loading.value = false;
  }
};
</script>
