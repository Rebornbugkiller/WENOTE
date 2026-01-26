import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGamificationStatus,
  getAchievements,
  updateDailyGoal,
  markAchievementNotified,
  getReport
} from '../api/gamification'

export const useGamificationStore = defineStore('gamification', () => {
  const status = ref(null)
  const achievements = ref([])
  const pendingNotifications = ref([])
  const isLoading = ref(false)

  // Computed
  const streakDisplay = computed(() => {
    if (!status.value) return { current: 0, longest: 0, atRisk: false }
    return {
      current: status.value.current_streak,
      longest: status.value.longest_streak,
      atRisk: status.value.streak_at_risk
    }
  })

  const goalProgress = computed(() => {
    if (!status.value) return { current: 0, target: 500, percent: 0, completed: false }
    return {
      current: status.value.today_chars,
      target: status.value.daily_char_goal,
      percent: Math.min(status.value.goal_progress * 100, 100),
      completed: status.value.goal_completed
    }
  })

  const unlockedAchievements = computed(() =>
    achievements.value.filter(a => a.unlocked)
  )

  const lockedAchievements = computed(() =>
    achievements.value.filter(a => !a.unlocked)
  )

  // Actions
  async function fetchStatus() {
    isLoading.value = true
    try {
      const res = await getGamificationStatus()
      status.value = res
      if (res.new_achievements?.length > 0) {
        pendingNotifications.value = res.new_achievements
      }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchAchievements() {
    const res = await getAchievements()
    achievements.value = res.list || []
  }

  async function setDailyGoal(goal) {
    await updateDailyGoal(goal)
    await fetchStatus()
  }

  async function dismissNotification(achievementId) {
    await markAchievementNotified(achievementId)
    pendingNotifications.value = pendingNotifications.value.filter(
      a => a.id !== achievementId
    )
  }

  async function fetchReport(period = 'week') {
    return await getReport(period)
  }

  function clearPendingNotifications() {
    pendingNotifications.value = []
  }

  return {
    status,
    achievements,
    pendingNotifications,
    isLoading,
    streakDisplay,
    goalProgress,
    unlockedAchievements,
    lockedAchievements,
    fetchStatus,
    fetchAchievements,
    setDailyGoal,
    dismissNotification,
    fetchReport,
    clearPendingNotifications
  }
})
