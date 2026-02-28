import client from './client'

export const getUserReviews = () =>
  client.get('/api/reviews').then((r) => r.data)

export const createReview = (data) =>
  client.post('/api/reviews', data).then((r) => r.data)
