import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMe } from '../api/auth'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  const setToken = (t) => {
    token.value = t
    localStorage.setItem('token', t)
  }

  const setUser = (u) => {
    user.value = u
  }

  const fetchUser = async () => {
    if (token.value) {
      user.value = await getMe()
    }
  }

  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  return { user, token, setToken, setUser, fetchUser, logout }
})
