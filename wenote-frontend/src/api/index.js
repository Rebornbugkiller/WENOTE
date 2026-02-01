import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'
import i18n from '../i18n'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 40000  // AI润色需要较长时间，设置为40秒
})

// 解析 JWT 获取过期时间
function getTokenExpiry(token) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.exp * 1000 // 转为毫秒
  } catch {
    return 0
  }
}

// 刷新令牌的状态管理
let isRefreshing = false
let refreshPromise = null

async function doRefreshToken(currentToken) {
  const response = await axios.post(
    (import.meta.env.VITE_API_BASE_URL || '/api/v1') + '/auth/refresh',
    {},
    { headers: { Authorization: `Bearer ${currentToken}` } }
  )
  if (response.data.code === 0) {
    return response.data.data.token
  }
  throw new Error('刷新令牌失败')
}

api.interceptors.request.use(async config => {
  const token = localStorage.getItem('token')
  if (token) {
    const expiry = getTokenExpiry(token)
    const now = Date.now()
    const sixHours = 6 * 60 * 60 * 1000

    // 距离过期不足6小时且令牌未过期，自动刷新
    if (expiry - now < sixHours && expiry > now) {
      if (!isRefreshing) {
        isRefreshing = true
        refreshPromise = doRefreshToken(token)
          .then(newToken => {
            localStorage.setItem('token', newToken)
            return newToken
          })
          .catch(() => token) // 刷新失败，使用原令牌
          .finally(() => {
            isRefreshing = false
            refreshPromise = null
          })
      }

      try {
        const newToken = await refreshPromise
        config.headers.Authorization = `Bearer ${newToken}`
      } catch {
        config.headers.Authorization = `Bearer ${token}`
      }
    } else {
      config.headers.Authorization = `Bearer ${token}`
    }
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
