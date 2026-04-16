import { useState, useEffect } from 'react'
import { getFlights, getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

function formatDuration(minutes) {
  if (!minutes) return '—'
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return `${h}h ${m}m`
}

export default function Flights() {
  const [flights, setFlights] = useState([])
  const [cities, setCities] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [fromCityId, setFromCityId] = useState('')
  const [toCityId, setToCityId] = useState('')
  const [maxPrice, setMaxPrice] = useState('')

  useEffect(() => {
    setLoading(true)
    Promise.all([getFlights(), getCities()])
      .then(([flightsData, citiesData]) => {
        setFlights(flightsData || [])
        setCities(citiesData || [])
      })
      .catch(err => {
        console.error("Error loading flights:", err)
        setError("Failed to load flight data.")
      })
      .finally(() => setLoading(false))
  }, [])

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
      label: 'Airline & Info',
      render: (val) => (
        <div className="flex flex-col">
          <span className="font-bold text-gray-800">{val}</span>
        </div>
      )
    },
    { 
      key: 'from_city_id', 
      label: 'Origin city',
      render: (id) => <span className="font-medium text-blue-600">{getCityName(id)}</span>
    },
    { 
      key: 'to_city_id', 
      label: 'Destination city',
      render: (id) => <span className="font-medium text-green-600">{getCityName(id)}</span>
    },
    { 
      key: 'duration_minutes', 
      label: 'Duration', 
      render: (v) => formatDuration(v) 
    },
    { 
      key: 'price', 
      label: 'Price', 
      render: (val) => <span className="font-bold text-green-700">${val.toLocaleString()}</span>
    }
  ]

  if (error) return <div className="p-8 text-red-600">{error}</div>

  return (
    <div>
      <PageHeader icon="✈️" title="Flights" subtitle="Manage routes and airfares" />

      <div className="mb-6 grid grid-cols-1 md:grid-cols-3 gap-4 items-end bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <div>
          <label className="label">From City</label>
          <select className="input" value={fromCityId} onChange={(e) => setFromCityId(e.target.value)}>
            <option value="">All Origins</option>
            {cities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
          </select>
        </div>

        <div>
          <label className="label">To City</label>
          <select className="input" value={toCityId} onChange={(e) => setToCityId(e.target.value)}>
            <option value="">All Destinations</option>
            {cities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
          </select>
        </div>

        <div>
          <label className="label">Max Price ($)</label>
          <input 
            type="number" 
            className="input" 
            placeholder="e.g. 500" 
            value={maxPrice} 
            onChange={(e) => setMaxPrice(e.target.value)}
          />
        </div>
      </div>

      <div className="card">
        <DataTable columns={columns} data={filteredFlights} loading={loading} />
      </div>
    </div>
  )
}