import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { createTrip, generateTripOptions, selectTripOption } from '../../api/trips'
import { getCities, getCountries } from '../../api/resources'
import PageHeader from '../../components/PageHeader'

const TIER_STYLE = {
  budget:   { border: 'border-green-300',  bg: 'bg-green-50',    badge: 'bg-green-100 text-green-800', icon: '🍃' },
  standard: { border: 'border-blue-300',   bg: 'bg-blue-50',     badge: 'bg-blue-100 text-blue-800',  icon: '🌟' },
  premium:  { border: 'border-purple-300', bg: 'bg-purple-50',   badge: 'bg-purple-100 text-purple-800', icon: '💎' },
}

export default function CreateTrip() {
  const navigate = useNavigate()
  const [step, setStep] = useState(1)
  const [tripId, setTripId] = useState(null)
  const [countries, setCountries] = useState([])
  const [cities, setCities] = useState([])
  const [filteredCities, setFilteredCities] = useState([])
  const [selectedCountry, setSelectedCountry] = useState('')
  const [options, setOptions] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [timeMode, setTimeMode] = useState('duration')

  const [form, setForm] = useState({
    name: '',
    destination_city_id: '',
    start_date: '',
    end_date: '',
    duration: '',
    total_price: '',
  })

  useEffect(() => {
    Promise.all([getCountries(), getCities()])
      .then(([countriesData, citiesData]) => {
        setCountries(countriesData || [])
        setCities(citiesData || [])
      })
      .catch(() => setError('Failed to load locations.'))
  }, [])

  const handleChange = (field) => (e) => setForm({ ...form, [field]: e.target.value })

  const handleCountryChange = (e) => {
    const countryId = e.target.value
    setSelectedCountry(countryId)
    setFilteredCities(cities.filter(c => c.country_id === parseInt(countryId)))
    setForm({ ...form, destination_city_id: '' })
  }

  const handleCreate = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    let startDate = form.start_date
    let endDate = form.end_date
    let duration = parseInt(form.duration)

    if (timeMode === 'duration') {
      const start = new Date() 
      const end = new Date()
      end.setDate(start.getDate() + duration)
      
      startDate = start.toISOString().split('T')[0]
      endDate = end.toISOString().split('T')[0]
    } else {
      const s = new Date(startDate)
      const e = new Date(endDate)
      duration = Math.ceil((e - s) / (1000 * 60 * 60 * 24))
    }

    try {
      const payload = {
        name: form.name,
        destination_city_id: parseInt(form.destination_city_id),
        total_price: parseFloat(form.total_price),
        duration: duration,
        start_date: startDate,
        end_date: endDate,
      }

      const { trip_id } = await createTrip(payload)
      setTripId(trip_id)
      const opts = await generateTripOptions(trip_id)
      setOptions(Array.isArray(opts) ? opts : [])
      setStep(2)
    } catch (err) {
      setError(err.response?.data || 'Failed to start planning. Check dates.')
    } finally {
      setLoading(false)
    }
  }

  const handleSelect = async (opt) => {
    setLoading(true)
    try {
      await selectTripOption(tripId, {
        tier: opt.tier,
        hotel_id: opt.hotel?.hotel_id || 0,
        outbound_flight_id: opt.outbound_flight?.flight_id || 0,
        inbound_flight_id: opt.inbound_flight?.flight_id || 0,
      })
      navigate(`/trips/${tripId}/itinerary`)
    } catch (err) {
      setError('Selection failed.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-6">
      <PageHeader
        icon="✈️"
        title="Plan a Trip"
        subtitle={step === 1 ? 'Fill in your trip details' : 'Choose your travel package'}
      />

      <div className="flex items-center justify-center gap-6 mb-10 mt-4">
        {[
          { n: 1, label: 'Trip Details' },
          { n: 2, label: 'Choose Package' },
          { n: 3, label: 'Itinerary' }
        ].map((s, i) => (
          <div key={s.n} className="flex items-center gap-3">
            <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm transition-all shadow-sm ${
              step === s.n ? 'bg-brand-600 text-white ring-4 ring-brand-100' : 
              step > s.n ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-500'
            }`}>
              {step > s.n ? '✓' : s.n}
            </div>
            <span className={`font-medium ${step >= s.n ? 'text-gray-900' : 'text-gray-400'}`}>
              {s.label}
            </span>
            {i < 2 && <span className="text-gray-300 font-light">→</span>}
          </div>
        ))}
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border-l-4 border-red-500 text-red-700 text-sm shadow-sm">{error}</div>}

      {step === 1 ? (
        <div className="card shadow-xl p-8 border-t-4 border-brand-600">
          <form onSubmit={handleCreate} className="space-y-6">
            <div>
              <label className="label">Trip Name</label>
              <input type="text" className="input" value={form.name} onChange={handleChange('name')} placeholder="My Adventure" required />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="label">Country</label>
                <select className="input" value={selectedCountry} onChange={handleCountryChange} required>
                  <option value="">Select country...</option>
                  {countries.map(c => <option key={c.country_id} value={c.country_id}>{c.name}</option>)}
                </select>
              </div>
              <div>
                <label className="label">City</label>
                <select className="input" value={form.destination_city_id} onChange={handleChange('destination_city_id')} disabled={!selectedCountry} required>
                  <option value="">Select city...</option>
                  {filteredCities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
                </select>
              </div>
            </div>

            <div className="space-y-4 bg-gray-50 p-4 rounded-xl border border-gray-100">
              <div className="flex gap-2 p-1 bg-gray-200/50 rounded-lg w-fit">
                <button type="button" onClick={() => setTimeMode('duration')} className={`px-4 py-1.5 rounded-md text-xs font-bold transition-all ${timeMode === 'duration' ? 'bg-white shadow text-brand-600' : 'text-gray-500'}`}>DURATION</button>
                <button type="button" onClick={() => setTimeMode('dates')} className={`px-4 py-1.5 rounded-md text-xs font-bold transition-all ${timeMode === 'dates' ? 'bg-white shadow text-brand-600' : 'text-gray-500'}`}>SPECIFIC DATES</button>
              </div>

              {timeMode === 'duration' ? (
                <div className="flex items-center gap-3">
                  <input type="number" className="input w-32" value={form.duration} onChange={handleChange('duration')} min="1" placeholder="Days" required />
                  <span className="text-gray-500 font-medium">days in the city</span>
                </div>
              ) : (
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-[10px] font-bold text-gray-400 uppercase">From</label>
                    <input type="date" className="input" value={form.start_date} onChange={handleChange('start_date')} required />
                  </div>
                  <div>
                    <label className="text-[10px] font-bold text-gray-400 uppercase">To</label>
                    <input type="date" className="input" value={form.end_date} onChange={handleChange('end_date')} required />
                  </div>
                </div>
              )}
            </div>

            <div>
              <label className="label">Budget (USD)</label>
              <input type="number" className="input" value={form.total_price} onChange={handleChange('total_price')} placeholder="e.g. 2000" required />
            </div>

            <button type="submit" disabled={loading} className="btn-primary w-full py-4 text-lg justify-center shadow-lg shadow-brand-200">
              {loading ? 'AI is generating options...' : 'Generate Plan →'}
            </button>
          </form>
        </div>
      ) : (
        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
          {options.map((opt) => {
            const style = TIER_STYLE[opt.tier] || {}
            return (
              <div key={opt.tier} className={`relative overflow-hidden rounded-2xl border-2 transition-all hover:scale-[1.01] ${style.border} ${style.bg} p-6 shadow-sm`}>
                <div className="flex flex-col md:flex-row justify-between gap-6">
                  <div className="flex-1 space-y-4">
                    <div className="flex items-center gap-2">
                      <span className="text-2xl">{style.icon}</span>
                      <span className={`px-3 py-1 rounded-full text-xs font-bold uppercase tracking-widest ${style.badge}`}>{opt.tier}</span>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                      <div className="space-y-1">
                        <p className="text-[10px] font-black text-gray-400 uppercase tracking-tighter">Accommodation</p>
                        <p className="font-bold text-gray-900 leading-tight">{opt.hotel?.name || 'Local Guesthouse'}</p>
                        <p className="text-xs text-gray-600">{opt.hotel?.stars} ⭐ · ${opt.hotel?.price_per_night}/night</p>
                      </div>
                      <div className="space-y-1">
                        <p className="text-[10px] font-black text-gray-400 uppercase tracking-tighter">Transport</p>
                        <p className="font-bold text-gray-900 leading-tight">{opt.outbound_flight?.airline || 'Standard Travel'}</p>
                        <p className="text-xs text-gray-600">Round trip total: ${((opt.outbound_flight?.price || 0) + (opt.inbound_flight?.price || 0)).toLocaleString()}</p>
                      </div>
                    </div>
                    
                    <div className="pt-4 border-t border-gray-200/50 flex flex-wrap gap-x-6 gap-y-2 text-xs text-gray-500">
                      <span>🎭 Activities Budget: <b className="text-gray-900">${opt.activites_budget}</b></span>
                      <span>🚗 Logistics: <b className="text-gray-900">${opt.logistics_budget}</b></span>
                      {opt.more_money > 0 && <span className="text-green-600 font-bold">💰 Saved: ${opt.more_money}</span>}
                    </div>
                  </div>

                  <div className="md:w-52 flex flex-col items-center justify-center bg-white rounded-xl p-5 border border-gray-100 shadow-inner">
                    <p className="text-xs text-gray-400 uppercase font-bold mb-1">Est. Total</p>
                    <p className="text-3xl font-black text-gray-900">${opt.total_price_of_money?.toLocaleString()}</p>
                    <button onClick={() => handleSelect(opt)} disabled={loading} className="btn-primary w-full mt-4 justify-center py-3">Select Plan</button>
                  </div>
                </div>
              </div>
            )
          })}
          <button onClick={() => setStep(1)} className="flex items-center gap-2 text-gray-500 font-bold hover:text-brand-600 transition-colors">
            ← Change Details
          </button>
        </div>
      )}
    </div>
  )
}