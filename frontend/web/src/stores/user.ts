import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, userApi } from '@/api'
import type { User } from '@/types'
import { useRouter } from 'vue-router'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const router = useRouter()

  const isLoggedIn = computed(() => !!user.value)

  const setUser = (userData: User) => {
    user.value = userData
  }

  const fetchProfile = async () => {
    const token = localStorage.getItem('token')
    if (!token) return

    try {
      const res = await userApi.getProfile()
      user.value = res.data
    } catch (error) {
      localStorage.removeItem('token')
      user.value = null
    }
  }

  const logout = async () => {
    try {
      await authApi.logout()
    } catch (error) {
      console.error(error)
    } finally {
      localStorage.removeItem('token')
      user.value = null
      router.push('/login')
    }
  }

  return {
    user,
    isLoggedIn,
    setUser,
    fetchProfile,
    logout,
  }
})
