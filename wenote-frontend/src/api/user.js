import api from './index'

// 获取当前用户资料
export const getProfile = () => api.get('/users/me')

// 更新用户资料
export const updateProfile = (data) => api.patch('/users/me', data)

// 修改密码
export const changePassword = (data) => api.post('/users/me/password', data)

// 注销账号
export const deleteAccount = (data) => api.delete('/users/me', { data })
