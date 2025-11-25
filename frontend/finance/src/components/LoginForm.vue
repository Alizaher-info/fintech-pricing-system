<template>
  <div v-if="!success" class="min-h-screen flex items-center justify-center p-4">
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
        <p class="text-blue-200 text-lg">Advanced Trading & Portfolio Management</p>
        <div class="flex justify-center space-x-6 mt-4 text-sm text-blue-300">
          <div class="flex items-center">
            <div class="w-2 h-2 bg-green-400 rounded-full mr-2 animate-pulse"></div>
            Markets Open
          </div>
          <div class="flex items-center">
            <div class="w-2 h-2 bg-blue-400 rounded-full mr-2 animate-pulse delay-1000"></div>
            Live Data
          </div>
        </div>
      </div>

      <!-- Login Form -->
      <div class="bg-white/10 backdrop-blur-xl rounded-3xl shadow-2xl p-8 border border-white/20 hover:bg-white/15 transition-all duration-300">
        <h2 class="text-2xl font-bold text-center mb-6 text-white">Sign In</h2>
        
        <div v-if="error" class="mb-4 p-3 bg-red-500/20 border border-red-500/50 text-red-200 rounded-xl backdrop-blur-sm">
          {{ error }}
          <div v-if="remainingAttempts" class="text-sm mt-2">
            Remaining attempts - IP: {{ remainingAttempts.ip }}, User: {{ remainingAttempts.user }}
          </div>
          <div v-if="lockoutRemaining" class="text-sm mt-2 font-semibold">
            🔒 Account locked for {{ Math.ceil(lockoutRemaining / 60) }} minutes
          </div>
        </div>

        <div v-if="success" class="mb-4 p-3 bg-green-500/20 border border-green-500/50 text-green-200 rounded-xl backdrop-blur-sm">
          ✓ Login successful!
        </div>

        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- Google Sign-In Button -->
          <div>
            <button
              type="button"
              @click="initGoogleSignIn"
              :disabled="loading"
              class="w-full flex items-center justify-center gap-3 bg-white text-gray-700 py-3 px-4 rounded-xl font-medium hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 focus:ring-offset-transparent transition-all duration-300 disabled:opacity-50 shadow-lg hover:shadow-xl border border-gray-200"
            >
              <svg class="w-5 h-5" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              <span>Continue with Google</span>
            </button>
            <!-- Hidden Google button container for SDK -->
            <div id="google-signin-button" class="hidden"></div>
          </div>

          <!-- Divider -->
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-white/30"></div>
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-4 bg-white/10 text-blue-200 backdrop-blur-sm rounded-full">Or continue with email</span>
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
                placeholder="Enter your email"
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
                class="block w-full px-4 py-3 pl-11 bg-white/10 border border-white/30 rounded-xl text-white placeholder-blue-200 focus:ring-2 focus:ring-emerald-400 focus:border-emerald-400 transition-all duration-300 backdrop-blur-sm"
                placeholder="Enter your password"
              />
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-blue-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
            </div>
          </div>

          <!-- Login Button -->
          <button
            type="submit"
            :disabled="loading || isLocked"
            class="w-full bg-gradient-to-r from-emerald-500 via-blue-500 to-purple-600 text-white py-4 px-4 rounded-xl font-semibold hover:from-emerald-600 hover:via-blue-600 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-emerald-400 focus:ring-offset-2 focus:ring-offset-transparent transition-all duration-300 disabled:opacity-50 shadow-lg hover:shadow-xl transform hover:scale-105"
          >
            <span v-if="loading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Signing in...
            </span>
            <span v-else-if="isLocked">🔒 Account Locked</span>
            <span v-else>Sign In</span>
          </button>
        </form>

        <!-- Sign Up Link -->
        <p class="mt-8 text-center text-sm text-blue-200">
          Don't have an account?
          <a href="#" @click.prevent="$emit('switch-to-register')" class="font-medium text-emerald-400 hover:text-emerald-300 transition-colors">
            Sign up for free
          </a>
        </p>
      </div>
    </div>
  </div>

  <!-- Success Animation -->
  <div v-else class="min-h-screen flex items-center justify-center p-4">
    <div class="text-center">
      <!-- Success Checkmark Animation -->
      <div class="relative w-32 h-32 mx-auto mb-6">
        <div class="absolute inset-0 bg-gradient-to-r from-emerald-500 to-blue-500 rounded-full animate-ping"></div>
        <div class="absolute inset-0 bg-gradient-to-r from-emerald-500 to-blue-500 rounded-full flex items-center justify-center shadow-2xl">
          <svg class="w-16 h-16 text-white animate-bounce" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
          </svg>
        </div>
      </div>
      
      <!-- Success Message -->
      <h2 class="text-3xl font-bold text-white mb-3">
        <span class="bg-gradient-to-r from-emerald-400 to-blue-400 bg-clip-text text-transparent">
          Login Successful!
        </span>
      </h2>
      <p class="text-blue-200 text-lg mb-4">Welcome back to FinTech Pro</p>
      
      <!-- Loading Bar -->
      <div class="max-w-xs mx-auto">
        <div class="h-2 bg-white/20 rounded-full overflow-hidden backdrop-blur-sm">
          <div class="h-full bg-gradient-to-r from-emerald-500 to-blue-500 rounded-full animate-pulse" style="width: 100%; animation-duration: 1s;"></div>
        </div>
        <p class="text-blue-300 text-sm mt-3">Redirecting to dashboard...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { LoginService } from '../services/LoginService';

const loginService = LoginService.getInstance();
const emit = defineEmits(['login-success', 'switch-to-register']);

const email = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');
const success = ref(false);
const remainingAttempts = ref<{ ip: number; user: number } | null>(null);
const lockoutRemaining = ref<number | null>(null);
const isLocked = ref(false);

const handleLogin = async () => {
  loading.value = true;
  error.value = '';
  success.value = false;
  remainingAttempts.value = null;
  lockoutRemaining.value = null;

  try {
    // LoginService auto-detects strategy from credentials (factory pattern)
    const response = await loginService.login({
      email: email.value,
      password: password.value,
    });

    // Token already saved by LoginService
    success.value = true;
    setTimeout(() => {
      emit('login-success', response.user);
    }, 1500);
  } catch (err: any) {
    error.value = err.message || err.error || 'Login failed';
    
    if (err.remaining_attempts) {
      remainingAttempts.value = err.remaining_attempts;
    }
    
    if (err.lockout_remaining) {
      lockoutRemaining.value = err.lockout_remaining;
      isLocked.value = true;
      
      // Auto-unlock after lockout period
      setTimeout(() => {
        isLocked.value = false;
        lockoutRemaining.value = null;
      }, err.lockout_remaining * 1000);
    }
  } finally {
    loading.value = false;
  }
};

const handleGoogleLogin = async (credential: string) => {
  loading.value = true;
  error.value = '';
  success.value = false;
  remainingAttempts.value = null;
  lockoutRemaining.value = null;

  try {
    console.log('🔐 Google Sign-In: Processing credential...');
    
    // LoginService will route to /api/login/google with credential
    const response = await loginService.login({
      credential: credential,
      oauthProvider: 'google'
    });

    console.log('✅ Google login successful');
    success.value = true;
    setTimeout(() => {
      emit('login-success', response.user);
    }, 1500);
  } catch (err: any) {
    console.error('❌ Google login failed:', err);
    error.value = err.message || err.error || 'Google login failed';
    
    if (err.remaining_attempts) {
      remainingAttempts.value = err.remaining_attempts;
    }
    
    if (err.lockout_remaining) {
      lockoutRemaining.value = err.lockout_remaining;
      isLocked.value = true;
      
      setTimeout(() => {
        isLocked.value = false;
        lockoutRemaining.value = null;
      }, err.lockout_remaining * 1000);
    }
  } finally {
    loading.value = false;
  }
};

const initGoogleSignIn = async () => {
  console.log('🔐 Initializing Google Sign-In...');
  error.value = '';
  
  // Check if Google Identity Services is loaded
  if (typeof (window as any).google === 'undefined') {
    error.value = 'Google Sign-In is not available. Please refresh the page.';
    console.error('❌ Google Identity Services not loaded');
    return;
  }

  try {
    const google = (window as any).google;
    console.log('✅ Google object found');
    
    // Create a temporary container for the button
    const container = document.createElement('div');
    container.id = 'temp-google-button';
    container.style.position = 'absolute';
    container.style.top = '-9999px';
    document.body.appendChild(container);
    
    // Initialize and render the button
    google.accounts.id.initialize({
      client_id: '128167116288-65occcq0bfoeigvh3icbchjcg8iqpuc0.apps.googleusercontent.com',
      callback: (response: any) => {
        console.log('🔐 Google Sign-In callback received');
        if (response.credential) {
          handleGoogleLogin(response.credential);
        } else {
          error.value = 'Google sign-in failed: No credential received';
        }
        // Clean up
        if (container && container.parentNode) {
          container.parentNode.removeChild(container);
        }
      },
      ux_mode: 'popup',
      context: 'signin'
    });

    // Render the button
    google.accounts.id.renderButton(container, {
      theme: 'outline',
      size: 'large',
      type: 'standard',
      text: 'signin_with'
    });

    console.log('✅ Google button rendered, triggering click...');
    
    // Wait a bit for rendering, then click
    setTimeout(() => {
      const btn = container.querySelector('div[role="button"]') as HTMLElement;
      if (btn) {
        console.log('✅ Found Google button, clicking...');
        btn.click();
      } else {
        console.error('❌ Could not find Google button');
        error.value = 'Failed to trigger Google Sign-In';
        if (container && container.parentNode) {
          container.parentNode.removeChild(container);
        }
      }
    }, 100);

  } catch (err) {
    console.error('❌ Google Sign-In initialization failed:', err);
    error.value = 'Failed to initialize Google Sign-In: ' + (err as Error).message;
  }
};

// Mount Google Sign-In callback on window object
onMounted(() => {
  console.log('🔐 LoginForm mounted, waiting for Google SDK...');
  
  // Wait for Google SDK to load
  const checkGoogleLoaded = setInterval(() => {
    if (typeof (window as any).google !== 'undefined') {
      console.log('✅ Google Identity Services loaded');
      clearInterval(checkGoogleLoaded);
    }
  }, 100);
  
  // Clear interval after 5 seconds if not loaded
  setTimeout(() => clearInterval(checkGoogleLoaded), 5000);
});
</script>
