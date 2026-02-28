import client from './client'

export const getCountries   = () => client.get('/api/countries').then((r) => r.data)
export const getCities      = () => client.get('/api/cities').then((r) => r.data)
export const getAttractions = () => client.get('/api/attractions').then((r) => r.data)
export const getHotels      = () => client.get('/api/hotels').then((r) => r.data)
export const getRestaurants = () => client.get('/api/restaurants').then((r) => r.data)
export const getFlights     = () => client.get('/api/flights').then((r) => r.data)
