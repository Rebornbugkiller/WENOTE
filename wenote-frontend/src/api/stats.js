import api from './index'

export const getStatsOverview = () => api.get('/stats/overview')
export const getStatsTrend = (days = 7) => api.get('/stats/trend', { params: { days } })
export const getStatsTags = (limit = 10) => api.get('/stats/tags', { params: { limit } })
export const getStatsNotebooks = () => api.get('/stats/notebooks')
