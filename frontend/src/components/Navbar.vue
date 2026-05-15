<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <nav class="border-b border-zinc-800 bg-bg-card/50 backdrop-blur sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
      <router-link to="/dashboard" class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600"></div>
        <span class="text-xl font-bold bg-gradient-to-r from-indigo-400 to-purple-400 bg-clip-text text-transparent">
          StepUp AI
        </span>
      </router-link>

      <div class="flex items-center gap-6">
        <router-link to="/dashboard" class="text-sm text-zinc-400 hover:text-white transition">
          Dashboard
        </router-link>
        <router-link to="/universities" class="text-sm text-zinc-400 hover:text-white transition">
          Universities
        </router-link>
        <router-link to="/essay" class="text-sm text-zinc-400 hover:text-white transition">
          Essay Review
        </router-link>
        <router-link to="/profile" class="text-sm text-zinc-400 hover:text-white transition">
          Profile
        </router-link>
        <button
          @click="handleLogout"
          class="text-sm text-zinc-400 hover:text-red-400 transition"
        >
          Logout
        </button>
      </div>
    </div>
  </nav>
</template>