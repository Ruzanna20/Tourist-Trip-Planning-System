import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { getTripItinerary, getEnrichedActivities } from '../../api/trips'
import PageHeader from '../../components/PageHeader'

const ACTIVITY_META = {
  hotel:      { icon: '🏨', color: 'bg-amber-50 border-amber-200',   label: 'Hotel Check-in' },
  attraction: { icon: '🎡', color: 'bg-purple-50 border-purple-200', label: 'Attraction' },
  restaurant: { icon: '🍽️', color: 'bg-orange-50 border-orange-200', label: 'Restaurant' },
  flight:     { icon: '🛫', color: 'bg-sky-50 border-sky-200',       label: 'Flight' },
}

function fmt(iso) {
  if (!iso) return '—'
  return new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function fmtDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('en-US', {
    weekday: 'long', year: 'numeric', month: 'long', day: 'numeric',
  })
}

function ActivityCard({ act }) {
  const meta = ACTIVITY_META[act.activity_type] ?? { icon: '📌', color: 'bg-gray-50 border-gray-200', label: act.activity_type }

  const renderDetail = () => {
    switch (act.activity_type) {
      case 'hotel':
        return (
          <>
            {act.entity_detail && <p className="text-xs text-gray-500">📍 {act.entity_detail}</p>}
            {act.entity_extra && (
              <p className="text-xs text-gray-500 line-clamp-2">{act.entity_extra}</p>
            )}
          </>
        )
      case 'attraction':
        return (
          <>
            {act.entity_detail && (
              <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs bg-purple-100 text-purple-700 capitalize">
                {act.entity_detail}
              </span>
            )}
            {act.entity_extra && Number(act.entity_extra) > 0 && (
              <p className="text-xs text-gray-500 mt-1">Entry fee: ${Number(act.entity_extra).toFixed(2)}</p>
            )}
            {act.entity_extra === '0' && (
              <p className="text-xs text-green-600 mt-1">Free entry</p>
            )}
          </>
        )
      case 'restaurant':
        return (
          <>
            {act.entity_detail && (
              <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs bg-orange-100 text-orange-700 capitalize">
                {act.entity_detail}
              </span>
            )}
            {act.entity_extra && (
              <p className="text-xs text-gray-500 mt-1">Price: {act.entity_extra}</p>
            )}
          </>
        )
      case 'flight':
        return (
          <>
            {act.entity_detail && <p className="text-xs text-gray-500">⏱ {act.entity_detail}</p>}
            {act.entity_extra && <p className="text-xs text-gray-500">💰 ${Number(act.entity_extra).toLocaleString()}</p>}
          </>
        )
      default:
        return null
    }
  }

  return (
    <div className={`flex gap-3 p-3 rounded-lg border ${meta.color}`}>
      <span className="text-xl flex-shrink-0 mt-0.5">{meta.icon}</span>
      <div className="flex-1 min-w-0 space-y-1">
        <div className="flex items-baseline justify-between gap-2">
          <p className="text-sm font-semibold text-gray-900 truncate">
            {act.entity_name || meta.label}
          </p>
          <p className="text-xs text-gray-400 flex-shrink-0">
            {fmt(act.start_time)} – {fmt(act.end_time)}
          </p>
        </div>
        {renderDetail()}
        {act.entity_rating > 0 && (
          <p className="text-xs text-yellow-600">★ {act.entity_rating.toFixed(1)}</p>
        )}
        {act.notes && (
          <p className="text-xs text-gray-500 italic border-t border-gray-200 pt-1 mt-1">{act.notes}</p>
        )}
      </div>
    </div>
  )
}

export default function Itinerary() {
  const { id } = useParams()
  const [days, setDays] = useState([])
  const [activities, setActivities] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getTripItinerary(id)
      .then((data) => {
        const list = Array.isArray(data) ? data : []
        setDays(list)
        list.forEach((day) => {
          getEnrichedActivities(day.Itinerary_id)
            .then((acts) =>
              setActivities((prev) => ({ ...prev, [day.Itinerary_id]: Array.isArray(acts) ? acts : [] })),
            )
            .catch(() =>
              setActivities((prev) => ({ ...prev, [day.Itinerary_id]: [] })),
            )
        })
      })
      .catch(() => setError('Failed to load itinerary. The plan may still be processing — try refreshing.'))
      .finally(() => setLoading(false))
  }, [id])

  return (
    <div>
      <PageHeader
        icon="📋"
        title="Trip Itinerary"
        subtitle={`Trip #${id}`}
        action={
          <Link to="/trips" className="btn-secondary">
            ← My Trips
          </Link>
        }
      />

      {loading && (
        <div className="flex items-center gap-2 text-gray-400 py-8">
          <div className="w-4 h-4 border-2 border-gray-300 border-t-brand-600 rounded-full animate-spin" />
          Loading itinerary…
        </div>
      )}

      {error && (
        <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>
      )}

      {!loading && !error && days.length === 0 && (
        <div className="card text-center py-16">
          <p className="text-4xl mb-3">⏳</p>
          <p className="font-medium text-gray-700">Itinerary is being generated</p>
          <p className="text-sm text-gray-400 mt-1">Your trip plan is still processing. Please refresh in a moment.</p>
        </div>
      )}

      <div className="space-y-5">
        {days.map((day) => {
          const dayActivities = activities[day.Itinerary_id]
          return (
            <div key={day.Itinerary_id} className="card">
              {/* Day header */}
              <div className="flex items-center gap-3 mb-4 pb-3 border-b border-gray-100">
                <div className="w-10 h-10 rounded-full bg-brand-600 text-white flex items-center justify-center font-bold text-sm flex-shrink-0">
                  {day.day_number}
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900">Day {day.day_number}</h3>
                  {day.date && <p className="text-xs text-gray-500">{fmtDate(day.date)}</p>}
                </div>
              </div>

              {day.notes && (
                <p className="text-sm text-gray-600 mb-3 italic">{day.notes}</p>
              )}

              {dayActivities === undefined ? (
                <div className="flex items-center gap-2 text-xs text-gray-400">
                  <div className="w-3 h-3 border border-gray-300 border-t-brand-400 rounded-full animate-spin" />
                  Loading activities…
                </div>
              ) : dayActivities.length === 0 ? (
                <p className="text-xs text-gray-400">No activities scheduled for this day.</p>
              ) : (
                <div className="space-y-2">
                  {dayActivities.map((act) => (
                    <ActivityCard key={act.activity_id} act={act} />
                  ))}
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
