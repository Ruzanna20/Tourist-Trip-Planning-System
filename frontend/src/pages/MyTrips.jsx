import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { getUserTrips, deleteTrip } from '../api/trips'
import PageHeader from '../components/PageHeader'

const STATUS_STYLE = {
  planned:    { badge: 'bg-blue-100 text-blue-800',    label: 'Planned' },
  pending:    { badge: 'bg-yellow-100 text-yellow-800', label: 'Pending' },
  processing: { badge: 'bg-indigo-100 text-indigo-800', label: 'Processing' },
  completed:  { badge: 'bg-green-100 text-green-800',  label: 'Completed' },
  cancelled:  { badge: 'bg-red-100 text-red-800',      label: 'Cancelled' },
}

function fmt(iso) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

function ConfirmDialog({ trip, onConfirm, onCancel, loading }) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" onClick={onCancel} />

      {/* Dialog */}
      <div className="relative bg-white rounded-2xl shadow-xl p-6 w-full max-w-sm">
        <div className="flex items-center justify-center w-12 h-12 rounded-full bg-red-100 mx-auto mb-4">
          <span className="text-2xl">🗑️</span>
        </div>
        <h3 className="text-center font-semibold text-gray-900 text-lg mb-1">Delete Trip</h3>
        <p className="text-center text-sm text-gray-500 mb-6">
          Are you sure you want to delete{' '}
          <span className="font-medium text-gray-700">"{trip.title}"</span>?
          This will permanently remove the trip, itinerary, and all activities.
        </p>
        <div className="flex gap-3">
          <button onClick={onCancel} disabled={loading} className="btn-secondary flex-1 justify-center">
            Cancel
          </button>
          <button onClick={onConfirm} disabled={loading} className="btn-danger flex-1 justify-center">
            {loading ? 'Deleting…' : 'Delete'}
          </button>
        </div>
      </div>
    </div>
  )
}

export default function MyTrips() {
  const [trips, setTrips] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [confirmTrip, setConfirmTrip] = useState(null)
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    getUserTrips()
      .then((data) => setTrips(Array.isArray(data) ? data : []))
      .catch(() => setError('Failed to load trips.'))
      .finally(() => setLoading(false))
  }, [])

  const handleDelete = async () => {
    if (!confirmTrip) return
    setDeleting(true)
    try {
      await deleteTrip(confirmTrip.trip_id)
      setTrips((prev) => prev.filter((t) => t.trip_id !== confirmTrip.trip_id))
      setConfirmTrip(null)
    } catch {
      setError('Failed to delete trip. Please try again.')
      setConfirmTrip(null)
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div>
      {confirmTrip && (
        <ConfirmDialog
          trip={confirmTrip}
          onConfirm={handleDelete}
          onCancel={() => setConfirmTrip(null)}
          loading={deleting}
        />
      )}

      <PageHeader
        icon="🗂️"
        title="My Trips"
        subtitle={loading ? 'Loading…' : `${trips.length} trip${trips.length !== 1 ? 's' : ''}`}
        action={
          <Link to="/trips/create" className="btn-primary">
            + Plan a Trip
          </Link>
        }
      />

      {error && (
        <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>
      )}

      {loading && (
        <div className="flex items-center gap-2 text-gray-400 py-8">
          <div className="w-4 h-4 border-2 border-gray-300 border-t-brand-600 rounded-full animate-spin" />
          Loading trips…
        </div>
      )}

      {!loading && trips.length === 0 && !error && (
        <div className="card text-center py-16">
          <p className="text-4xl mb-3">✈️</p>
          <p className="font-medium text-gray-700">No trips yet</p>
          <p className="text-sm text-gray-400 mt-1 mb-5">Create your first trip to get started.</p>
          <Link to="/trips/create" className="btn-primary inline-flex">
            Plan your first trip
          </Link>
        </div>
      )}

      {!loading && trips.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
          {trips.map((trip) => {
            const style = STATUS_STYLE[trip.status?.toLowerCase()] ?? STATUS_STYLE.planned
            return (
              <div key={trip.trip_id} className="card hover:shadow-md transition-shadow flex flex-col">
                {/* Header */}
                <div className="flex items-start justify-between mb-3">
                  <h3 className="font-semibold text-gray-900 text-base leading-tight flex-1 mr-2 truncate">
                    {trip.title}
                  </h3>
                  <span className={`badge ${style.badge} flex-shrink-0 capitalize`}>
                    {style.label}
                  </span>
                </div>

                {/* Dates */}
                <div className="flex items-center gap-1.5 text-sm text-gray-500 mb-2">
                  <span>📅</span>
                  <span>{fmt(trip.start_date)}</span>
                  <span className="text-gray-300">→</span>
                  <span>{fmt(trip.end_date)}</span>
                </div>

                {/* Budget */}
                {trip.total_price > 0 && (
                  <div className="flex items-center gap-1.5 text-sm text-gray-500 mb-4">
                    <span>💰</span>
                    <span>${Number(trip.total_price).toLocaleString('en-US')} budget</span>
                  </div>
                )}

                {/* Actions */}
                <div className="mt-auto pt-3 border-t border-gray-100 flex gap-2">
                  <Link
                    to={`/trips/${trip.trip_id}/itinerary`}
                    className="btn-primary flex-1 justify-center text-xs"
                  >
                    View Itinerary
                  </Link>
                  <button
                    onClick={() => setConfirmTrip(trip)}
                    className="flex-shrink-0 p-2 rounded-lg text-gray-400 hover:text-red-600 hover:bg-red-50 transition-colors"
                    title="Delete trip"
                  >
                    🗑️
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
