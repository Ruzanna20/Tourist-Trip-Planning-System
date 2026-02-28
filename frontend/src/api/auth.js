import client from './client'

export const login = (username, password) =>
  client.post('/login', { username, password }).then((r) => r.data)

export const refresh = (refreshToken) =>
  client.post('/refresh', { refresh_token: refreshToken }).then((r) => r.data)
