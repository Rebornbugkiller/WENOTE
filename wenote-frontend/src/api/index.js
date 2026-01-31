import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'
import i18n from '../i18n'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 40000  // AI润色需要较长时间，设置为40秒
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  response => {
    const { code, message, data } = response.data
    if (code !== 0) {
      ElMessage.error(message || i18n.global.t('common.requestFailed'))
      if (code === 401) {
        localStorage.removeItem('token')
        router.push('/login')
      }
      return Promise.reject(new Error(message))
    }
    return data
  },
  error => {
    ElMessage.error(error.message || i18n.global.t('common.networkError'))
    return Promise.reject(error)
  }
)

export default api
