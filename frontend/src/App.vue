<script setup lang="ts">
import { useRoute } from 'vue-router'
import Navbar from './components/Navbar.vue'

const route = useRoute()
</script>

<template>
  <div class="min-h-screen bg-bg-dark">
    <Navbar v-if="route.path !== '/login' && route.path !== '/register'" />
    <main>
      <router-view />
    </main>
  </div>
</template>