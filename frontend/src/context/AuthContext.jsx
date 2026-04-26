import { createContext, useContext, useState } from 'react'
import { login as loginApi } from '../api/auth'
import { jwtDecode } from 'jwt-decode' // Ավելացրու սա

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('token'))
  
  // Հաշվարկում ենք user-ի տվյալները token-ից
  const user = token ? jwtDecode(token) : null;

  const login = async (username, password) => {
    const data = await loginApi(username, password)
    localStorage.setItem('token', data.token)
    localStorage.setItem('refresh_token', data.refresh_token)
    setToken(data.token)
    return data
  }

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    setToken(null)
  }

  return (
    <AuthContext.Provider value={{ token, user, login, logout, isAuthenticated: !!token }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => useContext(AuthContext)