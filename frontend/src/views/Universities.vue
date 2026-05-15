<script setup lang="ts">
import { onMounted, ref } from 'vue'
import apiClient from '../api/client'

interface University {
  id: string
  name: string
  country: string
  acceptance_rate: number
  type: string
  category: string
}

const universities = ref<University[]>([])
const countryFilter = ref('')
const isLoading = ref(false)

async function loadUniversities() {
  isLoading.value = true
  try {
    const response = await apiClient.get('/universities', {
      params: { country: countryFilter.value },
    })
    universities.value = response.data.universities || []
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

async function saveUniversity(universityId: string) {
  try {
    await apiClient.post('/universities/save', { university_id: universityId })
  } catch (err) {
    console.error(err)
  }
}

onMounted(loadUniversities)
</script>

<template>
  <div class="max-w-7xl mx-auto px-6 py-10">
    <div class="mb-8">
      <h1 class="text-4xl font-bold mb-2">Universities</h1>
      <p class="text-zinc-400">Find universities that match your goals</p>
    </div>

    <div class="bg-bg-card border border-zinc-800 rounded-2xl p-6 mb-8">
      <div class="flex gap-4">
        <input
          v-model="countryFilter"
          type="text"
          placeholder="Filter by country..."
          class="flex-1 px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
          @keyup.enter="loadUniversities"
        />
        <button
          @click="loadUniversities"
          class="px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 hover:from-indigo-600 hover:to-purple-700 rounded-lg font-medium transition"
        >
          Search
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="text-center text-zinc-400 py-12">Loading...</div>

    <div v-else-if="universities.length === 0" class="text-center text-zinc-400 py-12">
      No universities found
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="university in universities"
        :key="university.id"
        class="bg-bg-card border border-zinc-800 rounded-2xl p-6 hover:border-indigo-500 transition"
      >
        <div class="flex items-start justify-between mb-4">
          <h3 class="text-xl font-semibold">{{ university.name }}</h3>
          <span
            class="px-2 py-1 text-xs rounded"
            :class="{
              'bg-red-500/20 text-red-400': university.category === 'reach',
              'bg-yellow-500/20 text-yellow-400': university.category === 'target',
              'bg-green-500/20 text-green-400': university.category === 'safety',
            }"
          >
            {{ university.category }}
          </span>
        </div>

        <p class="text-zinc-400 text-sm mb-1">{{ university.country }}</p>
        <p class="text-zinc-500 text-sm mb-4">
          Acceptance: {{ university.acceptance_rate }}%
        </p>

        <button
          @click="saveUniversity(university.id)"
          class="w-full py-2 bg-zinc-900 border border-zinc-800 hover:border-indigo-500 rounded-lg text-sm transition"
        >
          Save
        </button>
      </div>
    </div>
  </div>
</template>