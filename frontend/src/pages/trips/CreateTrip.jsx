import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { createTrip, generateTripOptions, selectTripOption } from '../../api/trips'
import { getCities } from '../../api/resources'
import PageHeader from '../../components/PageHeader'

const TIER_STYLE = {
  budget:   { border: 'border-green-300',  bg: 'bg-green-50',   badge: 'bg-green-100 text-green-800' },
  standard: { border: 'border-blue-300',   bg: 'bg-blue-50',    badge: 'bg-blue-100 text-blue-800' },
  premium:  { border: 'border-purple-300', bg: 'bg-purple-50',  badge: 'bg-purple-100 text-purple-800' },
}

const defaultStyle = { border: 'border-gray-200', bg: 'bg-white', badge: 'bg-gray-100 text-gray-800' }

export default function CreateTrip() {
  const navigate = useNavigate()
  const [step, setStep] = useState(1)
  const [tripId, setTripId] = useState(null)
  const [cities, setCities] = useState([])
  const [options, setOptions] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [form, setForm] = useState({
    name: '',
    destination_city_id: '',
    start_date: '',
    end_date: '',
    total_price: '',
  })

  useEffect(() => {
    getCities().then(setCities).catch(() => {})
  }, [])

  const set = (field) => (e) => setForm({ ...form, [field]: e.target.value })

  const handleCreate = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const { trip_id } = await createTrip({
        name: form.name,
        destination_city_id: parseInt(form.destination_city_id),
        start_date: form.start_date,
        end_date: form.end_date,
        total_price: parseFloat(form.total_price),
      })
      setTripId(trip_id)
      const opts = await generateTripOptions(trip_id)
      setOptions(Array.isArray(opts) ? opts : [])
      setStep(2)
    } catch (err) {
      setError(err.response?.data || 'Failed to create trip. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const handleSelect = async (opt) => {
    setError('')
    setLoading(true)
    try {
      await selectTripOption(tripId, {
        tier: opt.tier,
        hotel_id: opt.hotel?.hotel_id ?? 0,
        outbound_flight_id: opt.outbound_flight?.flight_id ?? 0,
        inbound_flight_id: opt.inbound_flight?.flight_id ?? 0,
      })
      navigate(`/trips/${tripId}/itinerary`)
    } catch (err) {
      setError(err.response?.data || 'Failed to confirm selection. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <PageHeader
        icon="✈️"
        title="Plan a Trip"
        subtitle={step === 1 ? 'Fill in your trip details' : 'Choose your travel package'}
      />

      {/* Step indicator */}
      <div className="flex items-center gap-2 mb-6 text-sm font-medium">
        {['Trip Details', 'Choose Package', 'Itinerary'].map((label, i) => {
          const n = i + 1
          const active = step === n
          const done = step > n
          return (
            <div key={n} className="flex items-center gap-2">
              <span
                className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold ${
                  active ? 'bg-brand-600 text-white' : done ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-500'
                }`}
              >
                {done ? '✓' : n}
              </span>
              <span className={active ? 'text-brand-600' : done ? 'text-green-600' : 'text-gray-400'}>{label}</span>
              {i < 2 && <span className="text-gray-300 mx-1">→</span>}
            </div>
          )
        })}
      </div>

      {error && (
        <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>
      )}

      {/* Step 1: Trip form */}
      {step === 1 && (
        <div className="card max-w-lg">
          <form onSubmit={handleCreate} className="space-y-4">
            <div>
              <label className="label">Trip Name</label>
              <input
                type="text"
                className="input"
                value={form.name}
                onChange={set('name')}
                placeholder="e.g. Summer Holiday 2025"
                required
              />
            </div>

            <div>
              <label className="label">Destination City</label>
              <select className="input" value={form.destination_city_id} onChange={set('destination_city_id')} required>
                <option value="">Select a city…</option>
                {cities.map((c) => (
                  <option key={c.city_id} value={c.city_id}>
                    {c.name} {c.iata_code ? `(${c.iata_code})` : ''}
                  </option>
                ))}
              </select>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="label">Start Date</label>
                <input type="date" className="input" value={form.start_date} onChange={set('start_date')} required />
              </div>
              <div>
                <label className="label">End Date</label>
                <input type="date" className="input" value={form.end_date} onChange={set('end_date')} required />
              </div>
            </div>

            <div>
              <label className="label">Total Budget (USD)</label>
              <input
                type="number"
                className="input"
                value={form.total_price}
                onChange={set('total_price')}
                min="0"
                step="0.01"
                placeholder="e.g. 3000"
                required
              />
            </div>

            <button type="submit" disabled={loading} className="btn-primary w-full justify-center">
              {loading ? 'Generating packages…' : 'Continue →'}
            </button>
          </form>
        </div>
      )}

      {/* Step 2: Choose option */}
      {step === 2 && (
        <div className="space-y-4 max-w-3xl">
          {options.length === 0 ? (
            <div className="card text-center py-12 text-gray-400">
              No package options available for this trip configuration.
            </div>
          ) : (
            options.map((opt) => {
              const style = TIER_STYLE[opt.tier] ?? defaultStyle
              return (
                <div
                  key={opt.tier}
                  className={`rounded-xl border-2 p-6 ${style.border} ${style.bg}`}
                >
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1">
                      <span className={`badge ${style.badge} uppercase tracking-wider mb-3 inline-block`}>
                        {opt.tier}
                      </span>

                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                        {opt.outbound_flight && (
                          <div className="space-y-0.5">
                            <p className="font-semibold text-gray-900">Outbound Flight</p>
                            <p className="text-gray-700">{opt.outbound_flight.airline}</p>
                            <p className="text-gray-500">{opt.outbound_flight.duration_minutes} min · ${opt.outbound_flight.price?.toLocaleString()}</p>
                          </div>
                        )}
                        {opt.inbound_flight && (
                          <div className="space-y-0.5">
                            <p className="font-semibold text-gray-900">Return Flight</p>
                            <p className="text-gray-700">{opt.inbound_flight.airline}</p>
                            <p className="text-gray-500">{opt.inbound_flight.duration_minutes} min · ${opt.inbound_flight.price?.toLocaleString()}</p>
                          </div>
                        )}
                        {opt.hotel && (
                          <div className="space-y-0.5">
                            <p className="font-semibold text-gray-900">Hotel</p>
                            <p className="text-gray-700">{opt.hotel.name} {'⭐'.repeat(opt.hotel.stars ?? 0)}</p>
                            <p className="text-gray-500">${opt.hotel.price_per_night}/night · Rating: {opt.hotel.rating}</p>
                          </div>
                        )}
                        <div className="space-y-0.5">
                          <p className="font-semibold text-gray-900">Budget Breakdown</p>
                          <p className="text-gray-700">Logistics: ${opt.logistics_budget?.toLocaleString()}</p>
                          <p className="text-gray-700">Activities: ${opt.activites_budget?.toLocaleString()}</p>
                          {opt.more_money > 0 && (
                            <p className="text-green-600 font-medium">+${opt.more_money?.toLocaleString()} remaining</p>
                          )}
                        </div>
                      </div>
                    </div>

                    <div className="flex-shrink-0 text-right">
                      <p className="text-2xl font-bold text-gray-900">
                        ${opt.total_price_of_money?.toLocaleString()}
                      </p>
                      <p className="text-xs text-gray-400 mb-4">total cost</p>
                      <button
                        onClick={() => handleSelect(opt)}
                        disabled={loading}
                        className="btn-primary"
                      >
                        {loading ? '…' : 'Select'}
                      </button>
                    </div>
                  </div>
                </div>
              )
            })
          )}

          <button onClick={() => setStep(1)} className="btn-secondary">
            ← Back to Details
          </button>
        </div>
      )}
    </div>
  )
}
