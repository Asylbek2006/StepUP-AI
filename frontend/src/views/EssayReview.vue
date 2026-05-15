<script setup lang="ts">
import { ref } from 'vue'
import apiClient from '../api/client'

const essayText = ref('')
const universityName = ref('')
const programName = ref('')
const wordLimit = ref(650)
const isLoading = ref(false)
const result = ref<any>(null)
const errorMessage = ref('')

async function handleSubmit() {
  errorMessage.value = ''
  isLoading.value = true
  result.value = null
  try {
    const response = await apiClient.post('/ai/essay', {
      essay_text: essayText.value,
      university_name: universityName.value,
      program_name: programName.value,
      word_limit: Number(wordLimit.value),
    })
    result.value = response.data
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Failed to review essay'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="max-w-5xl mx-auto px-6 py-10">
    <div class="mb-8">
      <h1 class="text-4xl font-bold mb-2">Essay Review</h1>
      <p class="text-zinc-400">Get AI-powered feedback on your application essay</p>
    </div>

    <div class="bg-bg-card border border-zinc-800 rounded-2xl p-8 mb-6">
      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm text-zinc-400 mb-2">University</label>
            <input
              v-model="universityName"
              type="text"
              required
              class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
              placeholder="MIT"
            />
          </div>

          <div>
            <label class="block text-sm text-zinc-400 mb-2">Program</label>
            <input
              v-model="programName"
              type="text"
              required
              class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
              placeholder="Computer Science"
            />
          </div>

          <div>
            <label class="block text-sm text-zinc-400 mb-2">Word limit</label>
            <input
              v-model="wordLimit"
              type="number"
              required
              class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm text-zinc-400 mb-2">Essay text</label>
          <textarea
            v-model="essayText"
            required
            rows="12"
            class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white resize-none"
            placeholder="Paste your essay here..."
          ></textarea>
          <p class="text-xs text-zinc-500 mt-1">
            {{ essayText.split(/\s+/).filter(w => w).length }} / {{ wordLimit }} words
          </p>
        </div>

        <p v-if="errorMessage" class="text-red-400 text-sm">{{ errorMessage }}</p>

        <button
          type="submit"
          :disabled="isLoading"
          class="px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 hover:from-indigo-600 hover:to-purple-700 rounded-lg font-medium transition disabled:opacity-50"
        >
          {{ isLoading ? 'Analyzing...' : 'Review essay' }}
        </button>
      </form>
    </div>

    <div v-if="result" class="bg-bg-card border border-zinc-800 rounded-2xl p-8">
      <h2 class="text-2xl font-bold mb-6">AI Feedback</h2>

      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <div class="bg-zinc-900 rounded-xl p-4 text-center">
          <p class="text-3xl font-bold text-indigo-400">{{ result.grammar_score?.toFixed(1) }}</p>
          <p class="text-xs text-zinc-400 mt-1">Grammar</p>
        </div>
        <div class="bg-zinc-900 rounded-xl p-4 text-center">
          <p class="text-3xl font-bold text-purple-400">{{ result.coherence_score?.toFixed(1) }}</p>
          <p class="text-xs text-zinc-400 mt-1">Coherence</p>
        </div>
        <div class="bg-zinc-900 rounded-xl p-4 text-center">
          <p class="text-3xl font-bold text-pink-400">{{ result.uniqueness_score?.toFixed(1) }}</p>
          <p class="text-xs text-zinc-400 mt-1">Uniqueness</p>
        </div>
        <div class="bg-zinc-900 rounded-xl p-4 text-center">
          <p class="text-3xl font-bold text-green-400">{{ result.relevance_score?.toFixed(1) }}</p>
          <p class="text-xs text-zinc-400 mt-1">Relevance</p>
        </div>
      </div>

      <div v-if="result.suggestions?.length">
        <h3 class="text-lg font-semibold mb-3">Improvement Suggestions</h3>
        <ul class="space-y-2">
          <li
            v-for="(suggestion, index) in result.suggestions"
            :key="index"
            class="flex gap-3 text-zinc-300 text-sm"
          >
            <span class="text-indigo-400">→</span>
            <span>{{ suggestion }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>