import client from './client'

export const getUserTrips = () =>
  client.get('/api/trips').then((r) => r.data)

export const deleteTrip = (tripId) =>
  client.delete(`/api/trips/${tripId}`).then((r) => r.data)

export const createTrip = (data) =>
  client.post('/api/trips/create', data).then((r) => r.data)

export const generateTripOptions = (tripId) =>
  client.post(`/api/trips/${tripId}/generate-options`).then((r) => r.data)

export const selectTripOption = (tripId, selection) =>
  client.post(`/api/trips/${tripId}/select-option`, selection).then((r) => r.data)

export const getTripItinerary = (tripId) =>
  client.get(`/api/trips/${tripId}/itinerary`).then((r) => r.data)

export const getItineraryActivities = (itineraryId) =>
  client.get(`/api/itineraries/${itineraryId}/activities`).then((r) => r.data)
