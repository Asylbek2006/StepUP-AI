import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const isAuthenticated = ref(!!localStorage.getItem('access_token'))

  async function register(email: string, password: string, fullName: string) {
    const response = await apiClient.post('/auth/register', {
      email,
      password,
      full_name: fullName,
    })
    localStorage.setItem('access_token', response.data.access_token)
    localStorage.setItem('refresh_token', response.data.refresh_token)
    isAuthenticated.value = true
  }

  async function login(email: string, password: string) {
    const response = await apiClient.post('/auth/login', { email, password })
    localStorage.setItem('access_token', response.data.access_token)
    localStorage.setItem('refresh_token', response.data.refresh_token)
    isAuthenticated.value = true
  }

  async function logout() {
    try {
      await apiClient.post('/auth/logout')
    } catch (e) {
      console.error(e)
    }
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    isAuthenticated.value = false
  }

  return { isAuthenticated, register, login, logout }
})