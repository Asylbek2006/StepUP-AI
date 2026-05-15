<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const errorMessage = ref('')
const isLoading = ref(false)

async function handleSubmit() {
  errorMessage.value = ''
  isLoading.value = true
  try {
    await authStore.login(email.value, password.value)
    router.push('/dashboard')
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Login failed'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="inline-block w-14 h-14 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 mb-4"></div>
        <h1 class="text-3xl font-bold bg-gradient-to-r from-indigo-400 to-purple-400 bg-clip-text text-transparent">
          Welcome back
        </h1>
        <p class="text-zinc-400 mt-2">Sign in to your StepUp AI account</p>
      </div>

      <div class="bg-bg-card border border-zinc-800 rounded-2xl p-8">
        <form @submit.prevent="handleSubmit" class="space-y-5">
          <div>
            <label class="block text-sm text-zinc-400 mb-2">Email</label>
            <input
              v-model="email"
              type="email"
              required
              class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
              placeholder="you@example.com"
            />
          </div>

          <div>
            <label class="block text-sm text-zinc-400 mb-2">Password</label>
            <input
              v-model="password"
              type="password"
              required
              class="w-full px-4 py-3 bg-zinc-900 border border-zinc-800 rounded-lg focus:border-indigo-500 focus:outline-none transition text-white"
              placeholder="••••••••"
            />
          </div>

          <p v-if="errorMessage" class="text-red-400 text-sm">{{ errorMessage }}</p>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-3 bg-gradient-to-r from-indigo-500 to-purple-600 hover:from-indigo-600 hover:to-purple-700 rounded-lg font-medium transition disabled:opacity-50"
          >
            {{ isLoading ? 'Signing in...' : 'Sign in' }}
          </button>
        </form>

        <p class="text-center text-sm text-zinc-400 mt-6">
          Don't have an account?
          <router-link to="/register" class="text-indigo-400 hover:text-indigo-300">Sign up</router-link>
        </p>
      </div>
    </div>
  </div>
</template>