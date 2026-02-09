<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

interface PriceData {
  symbol: string;
  price: number;
  volume_24h: number;
  high_24h: number;
  low_24h: number;
  timestamp: string;
}

interface ApiResponse {
  status: string;
  data: {
    BTC?: PriceData;
    ETH?: PriceData;
    SOL?: PriceData;
  };
  performance?: {
    total_ms: number;
    data_age: {
      BTC?: string;
      ETH?: string;
      SOL?: string;
    };
  };
}

const prices = ref<ApiResponse | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);
const lastUpdate = ref<Date>(new Date());

// Fetch prices from API
const fetchPrices = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/market-data/test');
    if (!response.ok) {
      throw new Error('Failed to fetch prices');
    }
    const data = await response.json();
    prices.value = data;
    lastUpdate.value = new Date();
    loading.value = false;
    error.value = null;
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Unknown error';
    loading.value = false;
  }
};

// Auto-refresh every 4 seconds
let intervalId: number | null = null;

onMounted(() => {
  fetchPrices();
  intervalId = window.setInterval(fetchPrices, 4000); // Refresh every 4 seconds
});

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId);
  }
});

// Helper to format price
const formatPrice = (price: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(price);
};

// Helper to format volume
const formatVolume = (volume: number) => {
  return volume.toFixed(2);
};

// Helper to get time ago
const getTimeAgo = () => {
  const seconds = Math.floor((Date.now() - lastUpdate.value.getTime()) / 1000);
  if (seconds < 5) return 'just now';
  return `${seconds}s ago`;
};
</script>

<template>
  <!-- Live Crypto Prices -->
  <div class="bg-white/10 backdrop-blur-md rounded-xl shadow-sm border border-white/10 overflow-hidden">
    <div class="p-6 border-b border-white/10 flex items-center justify-between bg-white/5">
      <div class="flex items-center gap-2">
        <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
        <h2 class="text-lg font-bold text-white">Live Crypto Prices</h2>
      </div>
      <span class="text-xs text-slate-400">{{ getTimeAgo() }}</span>
    </div>
    
    <div v-if="loading" class="p-8 text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-400"></div>
      <p class="text-slate-400 mt-2">Loading prices...</p>
    </div>

    <div v-else-if="error" class="p-8 text-center">
      <p class="text-red-400">{{ error }}</p>
      <button 
        @click="fetchPrices" 
        class="mt-4 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
      >
        Retry
      </button>
    </div>

    <div v-else-if="prices?.data" class="divide-y divide-white/10">
      <!-- Bitcoin -->
      <div v-if="prices.data.BTC" class="p-6 hover:bg-white/5 transition-colors">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 bg-orange-500/20 rounded-full flex items-center justify-center">
              <span class="text-orange-400 font-bold text-lg">₿</span>
            </div>
            <div>
              <h3 class="font-bold text-white text-lg">Bitcoin</h3>
              <p class="text-sm text-slate-400">BTC</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-2xl font-bold text-white">{{ formatPrice(prices.data.BTC.price) }}</p>
            <p class="text-xs text-slate-400 mt-1">Vol: {{ formatVolume(prices.data.BTC.volume_24h) }} BTC</p>
          </div>
        </div>
        <div class="flex gap-4 text-sm">
          <div>
            <span class="text-slate-400">High:</span>
            <span class="text-green-400 ml-1">{{ formatPrice(prices.data.BTC.high_24h) }}</span>
          </div>
          <div>
            <span class="text-slate-400">Low:</span>
            <span class="text-red-400 ml-1">{{ formatPrice(prices.data.BTC.low_24h) }}</span>
          </div>
        </div>
      </div>

      <!-- Ethereum -->
      <div v-if="prices.data.ETH" class="p-6 hover:bg-white/5 transition-colors">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 bg-blue-500/20 rounded-full flex items-center justify-center">
              <span class="text-blue-400 font-bold text-lg">Ξ</span>
            </div>
            <div>
              <h3 class="font-bold text-white text-lg">Ethereum</h3>
              <p class="text-sm text-slate-400">ETH</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-2xl font-bold text-white">{{ formatPrice(prices.data.ETH.price) }}</p>
            <p class="text-xs text-slate-400 mt-1">Vol: {{ formatVolume(prices.data.ETH.volume_24h) }} ETH</p>
          </div>
        </div>
        <div class="flex gap-4 text-sm">
          <div>
            <span class="text-slate-400">High:</span>
            <span class="text-green-400 ml-1">{{ formatPrice(prices.data.ETH.high_24h) }}</span>
          </div>
          <div>
            <span class="text-slate-400">Low:</span>
            <span class="text-red-400 ml-1">{{ formatPrice(prices.data.ETH.low_24h) }}</span>
          </div>
        </div>
      </div>

      <!-- Solana -->
      <div v-if="prices.data.SOL" class="p-6 hover:bg-white/5 transition-colors">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 bg-purple-500/20 rounded-full flex items-center justify-center">
              <span class="text-purple-400 font-bold text-lg">◎</span>
            </div>
            <div>
              <h3 class="font-bold text-white text-lg">Solana</h3>
              <p class="text-sm text-slate-400">SOL</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-2xl font-bold text-white">{{ formatPrice(prices.data.SOL.price) }}</p>
            <p class="text-xs text-slate-400 mt-1">Vol: {{ formatVolume(prices.data.SOL.volume_24h) }} SOL</p>
          </div>
        </div>
        <div class="flex gap-4 text-sm">
          <div>
            <span class="text-slate-400">High:</span>
            <span class="text-green-400 ml-1">{{ formatPrice(prices.data.SOL.high_24h) }}</span>
          </div>
          <div>
            <span class="text-slate-400">Low:</span>
            <span class="text-red-400 ml-1">{{ formatPrice(prices.data.SOL.low_24h) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Performance Info -->
    <div v-if="prices?.performance" class="p-4 bg-white/5 border-t border-white/10">
      <p class="text-xs text-slate-400 text-center">
        Response time: {{ prices.performance.total_ms }}ms | 
        Data age: {{ prices.performance.data_age?.BTC }}
      </p>
    </div>
  </div>
</template>
