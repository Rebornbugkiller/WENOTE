import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getProfile, updateProfile as apiUpdateProfile } from '../api/user'

// 预设头像映射
export const AVATAR_STYLES = {
  cat: '🐱', dog: '🐶', panda: '🐼', fox: '🦊',
  frog: '🐸', monkey: '🐵', rabbit: '🐰', bear: '🐻',
  gamepad: '🎮', palette: '🎨', books: '📚', music: '🎵',
  star: '🌟', rocket: '🚀', gem: '💎', fire: '🔥'
}

// 预设颜色
export const AVATAR_COLORS = [
  '#fbbf24', '#22c55e', '#3b82f6', '#8b5cf6', '#ef4444', '#ec4899'
]

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
      try {
        user.value = await getProfile()
      } catch (e) {
        console.error('Failed to fetch user:', e)
      }
    }
  }

  const updateProfile = async (data) => {
    const updated = await apiUpdateProfile(data)
    user.value = updated
    return updated
  }

  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  // 获取显示名称（昵称或用户名）
  const displayName = computed(() => {
    return user.value?.nickname || user.value?.username || 'Guest'
  })

  // 获取头像 emoji
  const avatarEmoji = computed(() => {
    const style = user.value?.avatar_style || 'cat'
    return AVATAR_STYLES[style] || '🐱'
  })

  // 获取头像背景色
  const avatarColor = computed(() => {
    return user.value?.avatar_color || '#fbbf24'
  })

  return {
    user,
    token,
    setToken,
    setUser,
    fetchUser,
    updateProfile,
    logout,
    displayName,
    avatarEmoji,
    avatarColor
  }
})
