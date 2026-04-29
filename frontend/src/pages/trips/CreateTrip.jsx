import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { createTrip, generateTripOptions, selectTripOption } from '../../api/trips'
import { getCities, getCountries } from '../../api/resources'
import PageHeader from '../../components/PageHeader'

export default function CreateTrip() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { id: urlTripId } = useParams()
  const [step, setStep] = useState(urlTripId ? 2 : 1)
  const [tripId, setTripId] = useState(urlTripId ? parseInt(urlTripId) : null)
  const [countries, setCountries] = useState([])
  const [cities, setCities] = useState([])
  const [filteredCities, setFilteredCities] = useState([])
  const [selectedCountry, setSelectedCountry] = useState('')
  const [options, setOptions] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [timeMode, setTimeMode] = useState('duration')

  const TIER_STYLE = {
    budget:   { border: 'border-green-300',  bg: 'bg-green-50',     badge: 'bg-green-100 text-green-800', icon: '🍃', label: t('trips.tiers.budget') },
    standard: { border: 'border-blue-300',   bg: 'bg-blue-50',      badge: 'bg-blue-100 text-blue-800',  icon: '🌟', label: t('trips.tiers.standard') },
    premium:  { border: 'border-purple-300', bg: 'bg-purple-50',     badge: 'bg-purple-100 text-purple-800', icon: '💎', label: t('trips.tiers.premium') },
  }

  const [form, setForm] = useState({
    name: '',
    destination_city_id: '',
    start_date: '',
    end_date: '',
    duration: '',
    total_price: '',
  })

  useEffect(() => {
    if (urlTripId) {
      setTripId(parseInt(urlTripId));
      setLoading(true);
      generateTripOptions(urlTripId)
      .then(opts => {
        setOptions(Array.isArray(opts) ? opts : []);
        setStep(2); 
      })
      .catch(() => setError(t('trips.error_load_options')))
      .finally(() => setLoading(false));
    }
  }, [urlTripId, t]);

  useEffect(() => {

  Promise.all([getCountries(), getCities()])
    .then(([countriesData, citiesData]) => {
      const sortedCountries = (countriesData || []).sort((a, b) => 
        a.name.localeCompare(b.name)
      );
      
      const sortedCities = (citiesData || []).sort((a, b) => 
        a.name.localeCompare(b.name)
      );

      setCountries(sortedCountries);
      setCities(sortedCities);
    })
    .catch(() => setError(t('countries.error_load')));
}, [t]);

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
      setError(t('trips.error_create'))
    } finally {
      setLoading(false)
    }
  }

  const handleSelect = async (opt) => {
  setLoading(true);
  setError('');
  try {
    await selectTripOption(tripId, {
      tier: opt.tier,
      hotel_id: opt.hotel?.hotel_id || 0,
      outbound_flight_id: opt.outbound_flight?.flight_id || 0,
      inbound_flight_id: opt.inbound_flight?.flight_id || 0,
    });
    navigate(`/trips/${tripId}/itinerary`);
  } catch (err) {
        const backendMessage = err.response?.data?.error;
        
        if (backendMessage) {
            setError(backendMessage); 
        } else {
            setError(t('trips.error_select'));
        }
        
        window.scrollTo({ top: 0, behavior: 'smooth' });
    } finally {
        setLoading(false)
    }
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-6">
      <PageHeader
        icon="✈️"
        title={t('nav.plan_trip')}
        subtitle={step === 1 ? t('trips.create_subtitle_1') : t('trips.create_subtitle_2')}
      />

      <div className="flex items-center justify-center gap-6 mb-10 mt-4">
        {[
          { n: 1, label: t('trips.steps.details') },
          { n: 2, label: t('trips.steps.package') },
          { n: 3, label: t('trips.steps.itinerary') }
        ].map((s, i) => (
          <div key={s.n} className="flex items-center gap-3">
            <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm transition-all shadow-sm ${
              step === s.n ? 'bg-blue-600 text-white ring-4 ring-blue-100' : 
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

      {error && (
        <div className="mb-6 p-4 bg-red-50 border-l-4 border-red-500 text-red-700 shadow-lg rounded-r-xl flex items-center gap-3">
          <div className="bg-red-500 text-white w-6 h-6 rounded-full flex items-center justify-center font-bold">!</div>
          <div className="flex-1">
            <p className="font-bold">{error}</p>
            {error.includes("նախընտրած") && (
              <button 
                onClick={() => navigate('/preferences')} 
                className="text-xs underline hover:text-red-900 mt-1 block"
              >
                Փոխել նախասիրությունները
              </button>
            )}
          </div>
        </div>
      )}

      {step === 1 ? (
        <div className="card shadow-xl p-8 border-t-4 border-blue-600">
          <form onSubmit={handleCreate} className="space-y-6">
            <div>
              <label className="label">{t('trips.form.name_label')}</label>
              <input type="text" className="input" value={form.name} onChange={handleChange('name')} placeholder={t('trips.form.name_placeholder')} required />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="label">{t('cities.table.country')}</label>
                <select className="input" value={selectedCountry} onChange={handleCountryChange} required>
                  <option value="">{t('trips.form.select_country')}</option>
                  {countries.map(c => <option key={c.country_id} value={c.country_id}>{c.name}</option>)}
                </select>
              </div>
              <div>
                <label className="label">{t('trips.form.cities')}</label>
                <select className="input" value={form.destination_city_id} onChange={handleChange('destination_city_id')} disabled={!selectedCountry} required>
                  <option value="">{t('trips.form.select_city')}</option>
                  {filteredCities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
                </select>
              </div>
            </div>

            <div className="space-y-4 bg-gray-50 p-4 rounded-xl border border-gray-100">
              <div className="flex gap-2 p-1 bg-gray-200/50 rounded-lg w-fit">
                <button type="button" onClick={() => setTimeMode('duration')} className={`px-4 py-1.5 rounded-md text-xs font-bold transition-all ${timeMode === 'duration' ? 'bg-white shadow text-blue-600' : 'text-gray-500'}`}>{t('trips.form.duration_mode')}</button>
                <button type="button" onClick={() => setTimeMode('dates')} className={`px-4 py-1.5 rounded-md text-xs font-bold transition-all ${timeMode === 'dates' ? 'bg-white shadow text-blue-600' : 'text-gray-500'}`}>{t('trips.form.dates_mode')}</button>
              </div>

              {timeMode === 'duration' ? (
                <div className="flex items-center gap-3">
                  <input type="number" className="input w-32" value={form.duration} onChange={handleChange('duration')} min="1" placeholder={t('flights.units.day')} required />
                  <span className="text-gray-500 font-medium">{t('trips.form.days_in_city')}</span>
                </div>
              ) : (
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-[10px] font-bold text-gray-400 uppercase">{t('flights.filter.from_date')}</label>
                    <input type="date" 
                      placeholder="օր / ամիս / տարի" 
                      className="input" value={form.start_date} onChange={handleChange('start_date')} required />
                  </div>
                  <div>
                    <label className="text-[10px] font-bold text-gray-400 uppercase">{t('flights.filter.to_date')}</label>
                    <input type="date" 
                      placeholder="օր / ամիս / տարի" 
                      className="input" value={form.end_date} onChange={handleChange('end_date')} required />
                  </div>
                </div>
              )}
            </div>

            <div>
              <label className="label">{t('trips.form.budget_label')}</label>
              <input 
                type="number" 
                className="input" 
                value={form.total_price} 
                onChange={handleChange('total_price')} 
                placeholder={t('trips.form.budget_placeholder')}
                required 
              />
            </div>

            <button type="submit" disabled={loading} className="btn-primary w-full py-4 text-lg justify-center shadow-lg">
              {loading ? t('trips.form.generating') : t('trips.form.generate_btn')}
            </button>
          </form>
        </div>
      ) : (
        <div className="space-y-6">
          {options.map((opt) => {
            const style = TIER_STYLE[opt.tier] || {}
            return (
              <div key={opt.tier} className={`relative overflow-hidden rounded-2xl border-2 transition-all hover:scale-[1.01] ${style.border} ${style.bg} p-6 shadow-sm`}>
                <div className="flex flex-col md:flex-row justify-between gap-6">
                  <div className="flex-1 space-y-4">
                    <div className="flex items-center gap-2">
                      <span className="text-2xl">{style.icon}</span>
                      <span className={`px-3 py-1 rounded-full text-xs font-bold uppercase tracking-widest ${style.badge}`}>{style.label}</span>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                      <div className="space-y-1">
                        <p className="text-[10px] font-black text-gray-400 uppercase tracking-tighter">{t('trips.package.accommodation')}</p>
                        <p className="font-bold text-gray-900 leading-tight">{opt.hotel?.name || t('trips.package.local_stay')}</p>
                        <p className="text-xs text-gray-600">{opt.hotel?.stars} ⭐ · ${opt.hotel?.price_per_night}/{t('flights.units.night')}</p>
                      </div>
                      <div className="space-y-1">
                        <p className="text-[10px] font-black text-gray-400 uppercase tracking-tighter">{t('trips.package.transport')}</p>
                        <p className="font-bold text-gray-900 leading-tight">{opt.outbound_flight?.airline || t('trips.package.standard_travel')}</p>
                        <p className="text-xs text-gray-600">{t('trips.package.round_trip')}: ${((opt.outbound_flight?.price || 0) + (opt.inbound_flight?.price || 0)).toLocaleString()}</p>
                      </div>
                    </div>
                    
                    <div className="pt-4 border-t border-gray-200/50 flex flex-wrap gap-x-6 gap-y-2 text-xs text-gray-500">
                      <span>🎭 {t('trips.package.activities')}: <b className="text-gray-900">${opt.activites_budget}</b></span>
                      <span>🚗 {t('trips.package.logistics')}: <b className="text-gray-900">${opt.logistics_budget}</b></span>
                      {opt.more_money > 0 && <span className="text-green-600 font-bold">💰 {t('trips.package.saved')}: ${opt.more_money}</span>}
                    </div>
                  </div>

                  <div className="md:w-52 flex flex-col items-center justify-center bg-white rounded-xl p-5 border border-gray-100 shadow-inner">
                    <p className="text-xs text-gray-400 uppercase font-bold mb-1">{t('trips.package.est_total')}</p>
                    <p className="text-3xl font-black text-gray-900">${opt.total_price_of_money?.toLocaleString()}</p>
                    <button onClick={() => handleSelect(opt)} disabled={loading} className="btn-primary w-full mt-4 justify-center py-3">{t('trips.package.select_btn')}</button>
                  </div>
                </div>
              </div>
            )
          })}
          <button onClick={() => setStep(1)} className="flex items-center gap-2 text-gray-500 font-bold hover:text-blue-600 transition-colors">
            ← {t('trips.package.back_btn')}
          </button>
        </div>
      )}
    </div>
  )
}