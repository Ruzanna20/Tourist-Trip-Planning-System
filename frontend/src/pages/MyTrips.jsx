import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { getUserTrips, deleteTrip } from '../api/trips'
import PageHeader from '../components/PageHeader'

export default function MyTrips() {
  const { t, i18n } = useTranslation()
  const [trips, setTrips] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [confirmTrip, setConfirmTrip] = useState(null)
  const [deleting, setDeleting] = useState(false)

  const STATUS_STYLE = {
    planned:    { badge: 'bg-blue-100 text-blue-800',    label: t('trips.status.planned') },
    pending:    { badge: 'bg-yellow-100 text-yellow-800', label: t('trips.status.pending') },
    processing: { badge: 'bg-indigo-100 text-indigo-800', label: t('trips.status.processing') },
    completed:  { badge: 'bg-green-100 text-green-800',  label: t('trips.status.completed') },
    cancelled:  { badge: 'bg-red-100 text-red-800',      label: t('trips.status.cancelled') },
  }

  function fmt(iso) {
  if (!iso) return '—';
  const d = new Date(iso);
  const day = d.getDate();
  const year = d.getFullYear();

  if (i18n.language === 'hy') {
    const armMonths = ["հնվ", "փտվ", "մրտ", "ապր", "մյս", "հնս", "հլս", "օգս", "սեպ", "հոկ", "նոյ", "դեկ"];
    return `${day} ${armMonths[d.getMonth()]} ${year}`;
  }

  return new Intl.DateTimeFormat('en-US', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  }).format(d);
}
  useEffect(() => {
    getUserTrips()
      .then((data) => setTrips(Array.isArray(data) ? data : []))
      .catch(() => setError(t('trips.error_load')))
      .finally(() => setLoading(false))
  }, [t])

  const handleDelete = async () => {
    if (!confirmTrip) return
    setDeleting(true)
    try {
      await deleteTrip(confirmTrip.trip_id)
      setTrips((prev) => prev.filter((t) => t.trip_id !== confirmTrip.trip_id))
      setConfirmTrip(null)
    } catch {
      setError(t('trips.error_delete'))
      setConfirmTrip(null)
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div>
      {confirmTrip && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" onClick={() => setConfirmTrip(null)} />
          <div className="relative bg-white rounded-2xl shadow-xl p-6 w-full max-w-sm text-center">
            <div className="flex items-center justify-center w-12 h-12 rounded-full bg-red-100 mx-auto mb-4 text-2xl">🗑️</div>
            <h3 className="font-semibold text-gray-900 text-lg mb-1">{t('trips.delete_confirm_title')}</h3>
            <p className="text-sm text-gray-500 mb-6">
              {t('trips.delete_confirm_msg')} <span className="font-medium text-gray-700">"{confirmTrip.title}"</span>?
            </p>
            <div className="flex gap-3">
              <button onClick={() => setConfirmTrip(null)} disabled={deleting} className="btn-secondary flex-1 justify-center">
                {t('common.cancel')}
              </button>
              <button onClick={handleDelete} disabled={deleting} className="btn-danger flex-1 justify-center">
                {deleting ? t('common.loading') : t('common.delete')}
              </button>
            </div>
          </div>
        </div>
      )}

      <PageHeader
        icon="🗂️"
        title={t('nav.my_trips')}
        subtitle={t('trips.subtitle')}
        action={
          <Link to="/trips/create" className="btn-primary">
            + {t('nav.plan_trip')}
          </Link>
        }
      />

      {error && <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>}

      {!loading && trips.length === 0 && !error && (
        <div className="card text-center py-16">
          <p className="text-4xl mb-3">✈️</p>
          <p className="font-medium text-gray-700">{t('trips.no_trips_title')}</p>
          <p className="text-sm text-gray-400 mt-1 mb-5">{t('trips.no_trips_msg')}</p>
          <Link to="/trips/create" className="btn-primary inline-flex">{t('trips.start_planning')}</Link>
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
        {trips.map((trip) => {
          const style = STATUS_STYLE[trip.status?.toLowerCase()] ?? STATUS_STYLE.planned
          return (
            <div key={trip.trip_id} className="card hover:shadow-md transition-shadow flex flex-col">
              <div className="flex items-start justify-between mb-3">
                <h3 className="font-semibold text-gray-900 text-base leading-tight flex-1 mr-2 truncate">{trip.title}</h3>
                <span className={`badge ${style.badge} flex-shrink-0`}>{style.label}</span>
              </div>
              <div className="flex items-center flex-wrap gap-x-2 gap-y-1 text-[13px] text-gray-600 mb-3 font-medium">
                <span className="text-base">📅</span>
                <span>{fmt(trip.start_date)}</span>
                <span className="text-gray-300">→</span>
                <span>{fmt(trip.end_date)}</span>
              </div>
              {trip.total_price > 0 && (
                <div className="flex items-center gap-1.5 text-sm text-gray-500 mb-4">
                  <span>💰</span> ${Number(trip.total_price).toLocaleString()} {t('trips.budget_label')}
                </div>
              )}
              <div className="mt-auto pt-3 border-t border-gray-100 flex gap-2">
                <Link to={`/trips/${trip.trip_id}/itinerary`} className="btn-primary flex-1 justify-center text-xs">
                  {t('trips.view_itinerary')}
                </Link>
                <button onClick={() => setConfirmTrip(trip)} className="p-2 rounded-lg text-gray-400 hover:text-red-600 hover:bg-red-50 transition-colors">
                  🗑️
                </button>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}