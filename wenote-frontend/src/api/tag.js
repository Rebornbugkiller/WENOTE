import api from './index'

export const getTags = () => api.get('/tags')
export const createTag = (data) => api.post('/tags', data)
export const updateTag = (id, data) => api.patch(`/tags/${id}`, data)
export const deleteTag = (id) => api.delete(`/tags/${id}`)
