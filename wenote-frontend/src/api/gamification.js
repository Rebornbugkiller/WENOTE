import api from './index'

export const getGamificationStatus = () => api.get('/gamification/status')
export const getAchievements = () => api.get('/gamification/achievements')
export const updateDailyGoal = (dailyCharGoal) => api.post('/gamification/goal', { daily_char_goal: dailyCharGoal })
export const getReport = (period = 'week') => api.get('/gamification/report', { params: { period } })
export const markAchievementNotified = (id) => api.post(`/gamification/achievements/${id}/notify`)
