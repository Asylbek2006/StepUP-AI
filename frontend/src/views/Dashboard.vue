<script setup lang="ts">
import { onMounted, ref } from 'vue'
import apiClient from '../api/client'

const savedUniversitiesCount = ref(0)
const profileGPA = ref<number | null>(null)
const isLoading = ref(true)

onMounted(async () => {
  try {
    const [profileResponse, savedResponse] = await Promise.all([
      apiClient.get('/profile'),
      apiClient.get('/universities/saved'),
    ])
    profileGPA.value = profileResponse.data.gpa
    savedUniversitiesCount.value = savedResponse.data.universities?.length || 0
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-6 py-10">
    <div class="mb-10">
      <h1 class="text-4xl font-bold mb-2">Dashboard</h1>
      <p class="text-zinc-400">Track your progress and applications</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
      <div class="bg-bg-card border border-zinc-800 rounded-2xl p-6">
        <div class="flex items-center justify-between mb-4">
          <span class="text-sm text-zinc-400">Current GPA</span>
          <div class="w-10 h-10 rounded-lg bg-indigo-500/20 flex items-center justify-center">
            <span class="text-indigo-400">📊</span>
          </div>
        </div>
        <p class="text-3xl font-bold">{{ profileGPA?.toFixed(2) || '—' }}</p>
        <p class="text-sm text-zinc-500 mt-2">Out of 4.0</p>
      </div>

      <div class="bg-bg-card border border-zinc-800 rounded-2xl p-6">
        <div class="flex items-center justify-between mb-4">
          <span class="text-sm text-zinc-400">Saved Universities</span>
          <div class="w-10 h-10 rounded-lg bg-purple-500/20 flex items-center justify-center">
            <span class="text-purple-400">🎓</span>
          </div>
        </div>
        <p class="text-3xl font-bold">{{ savedUniversitiesCount }}</p>
        <p class="text-sm text-zinc-500 mt-2">In your list</p>
      </div>

      <div class="bg-bg-card border border-zinc-800 rounded-2xl p-6">
        <div class="flex items-center justify-between mb-4">
          <span class="text-sm text-zinc-400">Application Status</span>
          <div class="w-10 h-10 rounded-lg bg-green-500/20 flex items-center justify-center">
            <span class="text-green-400">✓</span>
          </div>
        </div>
        <p class="text-3xl font-bold">Active</p>
        <p class="text-sm text-zinc-500 mt-2">Profile complete</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <router-link
        to="/universities"
        class="bg-bg-card border border-zinc-800 rounded-2xl p-6 hover:border-indigo-500 transition group"
      >
        <h3 class="text-xl font-semibold mb-2 group-hover:text-indigo-400 transition">
          Browse Universities
        </h3>
        <p class="text-zinc-400 text-sm">
          Find your dream university and analyze admission chances
        </p>
      </router-link>

      <router-link
        to="/essay"
        class="bg-bg-card border border-zinc-800 rounded-2xl p-6 hover:border-purple-500 transition group"
      >
        <h3 class="text-xl font-semibold mb-2 group-hover:text-purple-400 transition">
          Essay Review
        </h3>
        <p class="text-zinc-400 text-sm">
          Get AI-powered feedback on your application essays
        </p>
      </router-link>
    </div>
  </div>
</template>