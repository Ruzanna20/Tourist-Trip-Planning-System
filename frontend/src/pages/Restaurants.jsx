import { useState, useEffect } from 'react'
import { getRestaurants, getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

export default function Restaurants() {
  const [data, setData] = useState([])
  const [cities, setCities] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCityId, setSelectedCityId] = useState('')
  const [selectedCuisine, setSelectedCuisine] = useState('')
  const [minRating, setMinRating] = useState('')

  useEffect(() => {
    setLoading(true)
    Promise.all([getRestaurants(), getCities()])
      .then(([restData, citiesData]) => {
        const sortedData = [...restData].sort((a, b) => a.name.localeCompare(b.name));
        setData(sortedData);
        setCities(citiesData || []);
      })
      .catch(() => setError('Failed to load restaurant data.'))
      .finally(() => setLoading(false))
  }, [])

  const filteredData = data.filter(item => {
    const matchesSearch = item.name?.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesCity = selectedCityId ? item.city_id === parseInt(selectedCityId) : true
    const matchesCuisine = selectedCuisine ? item.cuisine === selectedCuisine : true
    const matchesRating = minRating ? item.rating >= parseFloat(minRating) : true
    return matchesSearch && matchesCity && matchesCuisine && matchesRating
  })

  const cuisines = [...new Set(data.map(r => r.cuisine).filter(Boolean))]

  const columns = [
    { key: 'name', label: 'Restaurant' },
    { 
      key: 'city_id', 
      label: 'City',
      render: (id) => cities.find(c => c.city_id === id)?.name || `ID: ${id}`
    },
    { 
      key: 'cuisine', 
      label: 'Cuisine', 
      render: (v) => v ? <span className="badge bg-orange-100 text-orange-700 capitalize">{v}</span> : '—' 
    },
    { 
      key: 'rating', 
      label: 'Rating', 
      render: (val) => (
        <span className={`px-2 py-1 rounded text-xs font-bold ${val >= 8 ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'}`}>
          {val.toFixed(1)} / 5
        </span>
      )
    },
    { 
      key: 'price_range', 
      label: 'Price Range', 
      render: (v) => v || '—' 
    },
  ]

  return (
    <div>
      <PageHeader icon="🍽️" title="Restaurants" subtitle={loading ? 'Loading…' : `${filteredData.length} restaurants`} />
      
      <div className="mb-6 grid grid-cols-1 md:grid-cols-4 gap-4 items-end bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <div>
          <label className="label">Search Name</label>
          <input type="text" className="input" placeholder="Name..." value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
        </div>
        <div>
          <label className="label">City</label>
          <select className="input" value={selectedCityId} onChange={(e) => setSelectedCityId(e.target.value)}>
            <option value="">All Cities</option>
            {cities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
          </select>
        </div>
        <div>
          <label className="label">Cuisine</label>
          <select className="input" value={selectedCuisine} onChange={(e) => setSelectedCuisine(e.target.value)}>
            <option value="">All Cuisines</option>
            {cuisines.map(c => <option key={c} value={c}>{c}</option>)}
          </select>
        </div>
        <div>
          <label className="label">Rating</label>
          <select className="input" value={minRating} onChange={(e) => setMinRating(e.target.value)}>
            <option value="">Rating</option>
            {[4.5, 4.0, 3.5, 3.0].map(r => <option key={r} value={r}>{r}+</option>)}
          </select>
        </div>
      </div>

      <div className="card">
        <DataTable columns={columns} data={filteredData} loading={loading} />
      </div>
    </div>
  )
}