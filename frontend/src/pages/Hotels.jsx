import { useState, useEffect } from 'react'
import { getHotels, getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

export default function Hotels() {
  const [hotels, setHotels] = useState([])
  const [cities, setCities] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCityId, setSelectedCityId] = useState('')
  const [selectedStars, setSelectedStars] = useState('')
  const [minRating, setMinRating] = useState('')
  const [maxPrice, setMaxPrice] = useState('')

  useEffect(() => {
    setLoading(true)
    Promise.all([getHotels(), getCities()])
      .then(([hotelsData, citiesData]) => {
        setHotels(hotelsData || [])
        setCities(citiesData || [])
      })
      .catch(err => {
        console.error("Error loading hotels:", err)
        setError("Failed to load hotel data.")
      })
      .finally(() => setLoading(false))
  }, [])

  const filteredHotels = hotels.filter(hotel => {
    const matchesSearch = hotel.name.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesCity = selectedCityId ? hotel.city_id === parseInt(selectedCityId) : true
    const matchesStars = selectedStars ? hotel.stars === parseInt(selectedStars) : true
    const matchesRating = minRating ? hotel.rating >= parseFloat(minRating) : true
    const matchesPrice = maxPrice ? hotel.price_per_night <= parseFloat(maxPrice) : true
    return matchesSearch && matchesCity && matchesStars && matchesRating && matchesPrice
  })

  const columns = [
    { key: 'name', label: 'Hotel Name' },
    { 
      key: 'city_id', 
      label: 'City',
      render: (id) => cities.find(c => c.city_id === id)?.name || `ID: ${id}`
    },
    { 
      key: 'stars', 
      label: 'Stars', 
      render: (val) => (
        <span className="text-yellow-500 font-bold">
          {'★'.repeat(val)}{'☆'.repeat(5 - val)}
        </span>
      )
    },
    { 
      key: 'rating', 
      label: 'Rating', 
      render: (val) => (
        <span className={`px-2 py-1 rounded text-xs font-bold ${val >= 8 ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'}`}>
          {val.toFixed(1)} / 10
        </span>
      )
    },
    { 
      key: 'price_per_night', 
      label: 'Price per Night', 
      render: (val) => <span className="font-semibold text-brand-700">${val.toLocaleString()}</span>
    }
  ]

  if (error) return <div className="p-8 text-red-600 font-medium">{error}</div>

  return (
    <div>
      <PageHeader icon="🏨" title="Hotels" subtitle="Manage accommodations and pricing" />

      <div className="mb-6 grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4 items-end">
        <div>
          <label className="label">Search Hotel</label>
          <input 
            type="text" 
            className="input" 
            placeholder="Hotel name..." 
            value={searchTerm} 
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        <div>
          <label className="label">City</label>
          <select className="input" value={selectedCityId} onChange={(e) => setSelectedCityId(e.target.value)}>
            <option value="">All Cities</option>
            {cities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
          </select>
        </div>

        <div>
          <label className="label">Stars</label>
          <select className="input" value={selectedStars} onChange={(e) => setSelectedStars(e.target.value)}>
            <option value="">Stars</option>
            {[5, 4, 3, 2, 1].map(s => <option key={s} value={s}>{s} Stars</option>)}
          </select>
        </div>

        <div>
          <label className="label">Rating</label>
          <select className="input" value={minRating} onChange={(e) => setMinRating(e.target.value)}>
            <option value="">Rating</option>
            {[9, 8, 7, 6].map(r => <option key={r} value={r}>{r}+ Exceptional</option>)}
          </select>
        </div>

        <div>
          <label className="label">Max Price ($)</label>
          <input 
            type="number" 
            className="input" 
            placeholder="Max price..." 
            value={maxPrice} 
            onChange={(e) => setMaxPrice(e.target.value)}
          />
        </div>
      </div>

      <div className="card">
        <DataTable columns={columns} data={filteredHotels} loading={loading} />
      </div>
    </div>
  )
}