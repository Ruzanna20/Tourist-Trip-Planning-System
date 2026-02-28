import client from './client'

export const registerUser = (data) =>
  client.post('/api/users/register', data).then((r) => r.data)

export const getUserPreferences = () =>
  client.get('/api/users/preferences').then((r) => r.data)

export const setPreferences = (data) =>
  client.post('/api/users/preferences', data).then((r) => r.data)
