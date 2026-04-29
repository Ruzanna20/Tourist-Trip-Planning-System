import { createContext, useContext, useState } from 'react'
import { login as loginApi } from '../api/auth'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('token'))
  const [user, setUser] = useState(() => {
    try {
      const savedUser = localStorage.getItem('user');
      if (savedUser && savedUser !== "undefined") {
        return JSON.parse(savedUser);
      }
    } catch (error) {
      console.error("Error parsing user:", error);
    }
    return null;
  });

  const login = async (email, password) => {
    try {
      const data = await loginApi(email, password)
      if (data.token) {
        localStorage.setItem('token', data.token)
        localStorage.setItem('refresh_token', data.refresh_token)
        setToken(data.token)
      }
      if (data.user) {
        localStorage.setItem('user', JSON.stringify(data.user))
        setUser(data.user)
      }
      return data
    } catch (error) {
      throw error;
    }
  }

  const logout = () => {
    localStorage.clear(); 
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ token, user, login, logout, isAuthenticated: !!token }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => useContext(AuthContext)