import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getFlights, getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

export default function Flights() {
  const { t } = useTranslation();
  const [flights, setFlights] = useState([])
  const [cities, setCities] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [fromCityId, setFromCityId] = useState('')
  const [toCityId, setToCityId] = useState('')
  const [maxPrice, setMaxPrice] = useState('')
  const [lastUpdated, setLastUpdated] = useState(null)

  const formatDuration = (minutes) => {
    if (!minutes) return '—'
    const h = Math.floor(minutes / 60)
    const m = minutes % 60
    return `${h}${t('flights.units.hour')} ${m}${t('flights.units.minute')}`
  }

  useEffect(() => {
    setLoading(true)
    Promise.all([getFlights(), getCities()])
      .then(([flightsData, citiesData]) => {
        if (flightsData && flightsData.length > 0) {
          const latest = flightsData.reduce((prev, current) => 
            (new Date(prev.updated_at) > new Date(current.updated_at)) ? prev : current
          );
          setLastUpdated(latest.updated_at);
        }

        const sortedCities = (citiesData || []).sort((a, b) => 
          a.name.localeCompare(b.name)
        );
        setFlights(flightsData || [])
        setCities(sortedCities)
      })
      .catch(err => {
        console.error("Error loading flights:", err)
        setError(t('flights.error_load'))
      })
      .finally(() => setLoading(false))
  }, [t])

  const filteredFlights = flights.filter(flight => {
    const matchesFrom = fromCityId ? flight.from_city_id === parseInt(fromCityId) : true
    const matchesTo = toCityId ? flight.to_city_id === parseInt(toCityId) : true
    const matchesPrice = maxPrice ? flight.price <= parseFloat(maxPrice) : true
    return matchesFrom && matchesTo && matchesPrice
  })

  const getCityName = (id) => cities.find(c => c.city_id === id)?.name || `ID: ${id}`

  const columns = [
    { 
      key: 'airline', 
      label: t('flights.table.airline'),
      render: (val) => (
        <div className="flex flex-col">
          <span className="font-bold text-gray-800">{val}</span>
        </div>
      )
    },
    { 
      key: 'from_city_id', 
      label: t('flights.table.origin'),
      render: (id) => <span className="font-medium text-blue-600">{getCityName(id)}</span>
    },
    { 
      key: 'to_city_id', 
      label: t('flights.table.destination'),
      render: (id) => <span className="font-medium text-green-600">{getCityName(id)}</span>
    },
    { 
      key: 'duration_minutes', 
      label: t('flights.table.duration'), 
      render: (v) => formatDuration(v) 
    },
    { 
      key: 'price', 
      label: t('flights.table.price'), 
      render: (val) => <span className="font-bold text-green-700">${val.toLocaleString()}</span>
    },
    { 
      key: 'website', 
      label: t('flights.table.source'), 
      render: (url) => {
        if (!url) return '—';
        const href = url.startsWith('http') ? url : `https://${url}`;
        return (
          <a 
            href={href} 
            target="_blank" 
            rel="noopener noreferrer" 
            className="inline-flex items-center px-3 py-1 bg-gray-50 text-gray-600 rounded-lg hover:bg-gray-100 transition-colors border border-gray-200 font-medium text-xs"
          >
            <span>{t('flights.visit_source')} {url}</span>
            <span className="ml-1 text-[10px]">↗</span>
          </a>
        );
      }
    }
  ]

  if (error) return <div className="p-8 text-red-600 font-medium">{error}</div>

  return (
    <div>
      <PageHeader 
      icon="✈️" 
      title={t('nav.flights')}
      subtitle={t('flights.subtitle')} 
      />

      {lastUpdated && (
          <div className="mb-4 md:mb-6 flex items-center gap-2 px-3 py-1.5 bg-white border border-slate-100 rounded-full shadow-sm w-fit transition-all hover:border-slate-200">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            <span className="text-[10px] font-black uppercase tracking-widest text-slate-400">
              {t('common.last_updated')} {new Date(lastUpdated).toLocaleString()}
            </span>
          </div>
        )}

      <div className="mb-6 grid grid-cols-1 md:grid-cols-3 gap-4 items-end bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <div>
          <label className="label">{t('flights.filter.from_label')}</label>
          <select className="input" value={fromCityId} onChange={(e) => setFromCityId(e.target.value)}>
            <option value="">{t('flights.filter.all_origins')}</option>
            {cities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
          </select>
        </div>

        <div>
          <label className="label">{t('flights.filter.to_label')}</label>
          <select className="input" value={toCityId} onChange={(e) => setToCityId(e.target.value)}>
            <option value="">{t('flights.filter.all_destinations')}</option>
            {cities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
          </select>
        </div>

        <div>
          <label className="label">{t('flights.filter.price_label')}</label>
          <input 
            type="number" className="input" placeholder={t('flights.filter.price_placeholder')}
            value={maxPrice} onChange={(e) => setMaxPrice(e.target.value)}
          />
        </div>
      </div>

      <div className="card shadow-md overflow-hidden">
        <DataTable columns={columns} data={filteredFlights} loading={loading} />
      </div>
    </div>
  )
}