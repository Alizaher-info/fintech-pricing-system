<script setup lang="ts">
import { inject, type Ref } from 'vue';
import Dashboard from '../components/Dashboard.vue';

const user = inject<Ref<any>>('user');
</script>

<template>
  <Dashboard v-if="user" :user="user" />
</template>
