<script setup lang="ts">
import { onMounted, ref } from 'vue'
import apiClient from '../api/client'

const gpa = ref(0)
const satScore = ref(0)
const ieltsScore = ref(0)
const country = ref('')
const isLoading = ref(true)
const isSaving = ref(false)
const successMessage = ref('')

onMounted(async () => {
  try {
    const response = await apiClient.get('/profile')
    gpa.value = response.data.gpa || 0
    satScore.value = response.data.sat_score || 0
    ieltsScore.value = response.data.ielts_score || 0
    country.value = response.data.country || ''
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
})

async function handleSave() {
  isSaving.value = true
  successMessage.value = ''
  try {
    await apiClient.put('/profile', {
      gpa: Number(gpa.value),
      sat_score: Number(satScore.value),
      ielts_score: Number(ieltsScore.value),
      country: country.value,
    })
    successMessage.value = 'Profile updated successfully'
    setTimeout(() => (successMessage.value = ''), 3000)
  } catch (err) {
    console.error(err)
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 py-10">
    <div class="mb-8">
      <h1 class="text-4xl font-bold mb-2">Profile</h1>
      <p class="text-zinc-400">Manage your academic information</p>
    </div>

    <div class="bg-bg-card border border-zinc-800 rounded-2xl p-8">
      <form @submit.prevent="handleSave" class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm text-zinc-400 mb-2">GPA</label>
            <input
              v-model="gpa"
              type="number"
              step="0.01"
              min="0"
              max="4"
              class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
            />
          </div>

          <div>
            <label class="block text-sm text-zinc-400 mb-2">SAT Score</label>
            <input
              v-model="satScore"
              type="number"
              min="400"
              max="1600"
              class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
            />
          </div>

          <div>
            <label class="block text-sm text-zinc-400 mb-2">IELTS Score</label>
            <input
              v-model="ieltsScore"
              type="number"
              step="0.5"
              min="0"
              max="9"
              class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
            />
          </div>

          <div>
            <label class="block text-sm text-zinc-400 mb-2">Country</label>
            <input
              v-model="country"
              type="text"
              class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
              placeholder="Kazakhstan"
            />
          </div>
        </div>

        <p v-if="successMessage" class="text-green-400 text-sm">{{ successMessage }}</p>

        <button
          type="submit"
          :disabled="isSaving"
          class="px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 hover:from-indigo-600 hover:to-purple-700 rounded-lg font-medium transition disabled:opacity-50"
        >
          {{ isSaving ? 'Saving...' : 'Save changes' }}
        </button>
      </form>
    </div>
  </div>
</template>