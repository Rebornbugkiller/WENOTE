import api from './index'

export const getNotebooks = () => api.get('/notebooks')
export const getDefaultNotebook = () => api.get('/notebooks/default')
export const createNotebook = (data) => api.post('/notebooks', data)
export const updateNotebook = (id, data) => api.patch(`/notebooks/${id}`, data)
export const deleteNotebook = (id) => api.delete(`/notebooks/${id}`)
