<template>
  <div class="min-h-screen relative overflow-hidden">
    <!-- Global Animated Financial Background -->
    <div class="fixed inset-0 bg-gradient-to-br from-slate-900 via-blue-900 to-indigo-900 z-0">
      <!-- Floating Financial Elements -->
      <div class="absolute inset-0">
        <!-- Stock Chart Lines -->
        <svg class="absolute top-10 left-10 w-32 h-20 text-green-400 opacity-20 animate-pulse" viewBox="0 0 128 80" fill="none">
          <path d="M0 40 L20 30 L40 45 L60 25 L80 35 L100 15 L128 20" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
        <svg class="absolute top-32 right-20 w-40 h-24 text-red-400 opacity-20 animate-pulse delay-1000" viewBox="0 0 160 96" fill="none">
          <path d="M0 60 L25 40 L50 65 L75 35 L100 55 L125 25 L160 30" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
        <svg class="absolute bottom-32 left-20 w-36 h-22 text-blue-400 opacity-20 animate-pulse delay-2000" viewBox="0 0 144 88" fill="none">
          <path d="M0 50 L24 35 L48 60 L72 30 L96 45 L120 20 L144 25" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>

        <!-- More Chart Lines for Coverage -->
        <svg class="absolute top-60 left-1/3 w-28 h-16 text-yellow-400 opacity-15 animate-pulse delay-3000" viewBox="0 0 112 64" fill="none">
          <path d="M0 32 L16 20 L32 40 L48 15 L64 30 L80 10 L96 18 L112 12" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
        <svg class="absolute bottom-40 right-1/3 w-30 h-18 text-purple-400 opacity-15 animate-pulse delay-4000" viewBox="0 0 120 72" fill="none">
          <path d="M0 36 L20 25 L40 50 L60 20 L80 40 L100 15 L120 22" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>

        <!-- Floating Numbers -->
        <div class="absolute top-20 right-32 text-green-400 font-mono text-sm opacity-30 animate-bounce">
          +2.45%
        </div>
        <div class="absolute top-48 left-16 text-red-400 font-mono text-sm opacity-30 animate-bounce delay-500">
          -1.23%
        </div>
        <div class="absolute bottom-40 right-16 text-blue-400 font-mono text-sm opacity-30 animate-bounce delay-1000">
          $1,234.56
        </div>
        <div class="absolute bottom-20 left-32 text-yellow-400 font-mono text-sm opacity-30 animate-bounce delay-1500">
          ₿65,432
        </div>
        <div class="absolute top-1/3 left-1/2 text-emerald-400 font-mono text-xs opacity-25 animate-bounce delay-2000">
          TSLA +3.2%
        </div>
        <div class="absolute bottom-1/2 right-1/4 text-orange-400 font-mono text-xs opacity-25 animate-bounce delay-2500">
          NVDA -0.8%
        </div>

        <!-- Currency Symbols -->
        <div class="absolute top-40 left-1/4 text-6xl text-white opacity-5 animate-spin-slow">$</div>
        <div class="absolute bottom-1/3 right-1/4 text-5xl text-white opacity-5 animate-spin-slow delay-2000">€</div>
        <div class="absolute top-1/3 right-1/3 text-4xl text-white opacity-5 animate-spin-slow delay-4000">£</div>
        <div class="absolute top-2/3 left-1/3 text-5xl text-white opacity-5 animate-spin-slow delay-6000">¥</div>

        <!-- Candlestick Charts -->
        <div class="absolute top-60 left-1/2 transform -translate-x-1/2">
          <div class="flex space-x-1 opacity-20">
            <div class="w-2 bg-green-400 animate-pulse" style="height: 24px;"></div>
            <div class="w-2 bg-red-400 animate-pulse delay-200" style="height: 32px;"></div>
            <div class="w-2 bg-green-400 animate-pulse delay-400" style="height: 20px;"></div>
            <div class="w-2 bg-green-400 animate-pulse delay-600" style="height: 28px;"></div>
            <div class="w-2 bg-red-400 animate-pulse delay-800" style="height: 16px;"></div>
          </div>
        </div>

        <!-- More Candlesticks -->
        <div class="absolute bottom-1/4 left-1/4">
          <div class="flex space-x-1 opacity-15">
            <div class="w-1.5 bg-blue-400 animate-pulse delay-1000" style="height: 18px;"></div>
            <div class="w-1.5 bg-green-400 animate-pulse delay-1200" style="height: 22px;"></div>
            <div class="w-1.5 bg-red-400 animate-pulse delay-1400" style="height: 14px;"></div>
            <div class="w-1.5 bg-green-400 animate-pulse delay-1600" style="height: 26px;"></div>
          </div>
        </div>

        <!-- Particle Effect -->
        <div class="absolute inset-0">
          <div class="absolute top-1/4 left-1/4 w-1 h-1 bg-blue-400 rounded-full animate-ping opacity-50"></div>
          <div class="absolute top-3/4 right-1/4 w-1 h-1 bg-green-400 rounded-full animate-ping opacity-50 delay-1000"></div>
          <div class="absolute bottom-1/4 left-3/4 w-1 h-1 bg-yellow-400 rounded-full animate-ping opacity-50 delay-2000"></div>
          <div class="absolute top-1/2 right-1/2 w-1 h-1 bg-purple-400 rounded-full animate-ping opacity-50 delay-3000"></div>
          <div class="absolute bottom-1/3 left-1/2 w-1 h-1 bg-red-400 rounded-full animate-ping opacity-50 delay-4000"></div>
        </div>

        <!-- Market Grid Lines -->
        <div class="absolute inset-0 opacity-5">
          <div class="h-full w-full bg-gradient-to-r from-transparent via-white to-transparent" style="background-size: 100px 1px; background-image: repeating-linear-gradient(90deg, transparent, transparent 99px, rgba(255,255,255,0.1) 100px);"></div>
          <div class="h-full w-full bg-gradient-to-b from-transparent via-white to-transparent" style="background-size: 1px 100px; background-image: repeating-linear-gradient(0deg, transparent, transparent 99px, rgba(255,255,255,0.1) 100px);"></div>
        </div>
      </div>

      <!-- Overlay Gradient -->
      <div class="absolute inset-0 bg-gradient-to-br from-black/20 via-transparent to-black/40"></div>
    </div>

    <!-- Content -->
    <div class="relative z-10">
      <router-view />
    </div>
  </div>
</template>

<style scoped>
@keyframes spin-slow {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin-slow {
  animation: spin-slow 15s linear infinite;
}
</style>
