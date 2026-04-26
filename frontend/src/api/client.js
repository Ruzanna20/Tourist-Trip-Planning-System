import axios from 'axios'

const API_URL = 'http://localhost:8080'

const client = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' },
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    const cleanToken = token.replace(/^"(.*)"$/, '$1'); 
    config.headers.Authorization = `Bearer ${cleanToken}`;
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config
    if (error.response?.status === 401 && !original._retry) {
      original._retry = true
      const refreshToken = localStorage.getItem('refresh_token')
      
      if (refreshToken) {
        try {
          const { data } = await axios.post(`${API_URL}/refresh`, { 
            refresh_token: refreshToken 
          })
          
          localStorage.setItem('token', data.token)
          const cleanToken = data.token.replace(/^"(.*)"$/, '$1');
          client.defaults.headers.common['Authorization'] = `Bearer ${cleanToken}`
          original.headers.Authorization = `Bearer ${cleanToken}`
          
          return client(original)
        } catch (refreshError) {
          localStorage.removeItem('token')
          localStorage.removeItem('refresh_token')
          window.location.href = '/login'
          return Promise.reject(refreshError)
        }
      }
    }
    return Promise.reject(error)
  },
)

export default client