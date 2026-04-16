import { useState, useEffect } from 'react'
import { getAttractions, getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

const CATEGORY_COLORS = {
  museum: 'bg-blue-100 text-blue-800 border-blue-200',
  viewpoint: 'bg-green-100 text-green-800 border-green-200',
  gallery: 'bg-purple-100 text-purple-800 border-purple-200',
  attraction: 'bg-orange-100 text-orange-800 border-orange-200',
  monument: 'bg-indigo-100 text-indigo-800 border-indigo-200',
  historic: 'bg-red-100 text-red-800 border-red-200',
}

export default function Attractions() {
  const [attractions, setAttractions] = useState([])
  const [cities, setCities] = useState([])
  const [loading, setLoading] = useState(true)
  
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCityId, setSelectedCityId] = useState('')
  const [selectedCategory, setSelectedCategory] = useState('')
  const [minRating, setMinRating] = useState('')

  useEffect(() => {
    setLoading(true)
    Promise.all([getAttractions(), getCities()])
      .then(([attractionsData, citiesData]) => {
        setAttractions(attractionsData || [])
        setCities(citiesData || [])
      })
      .catch(err => console.error("Error loading data:", err))
      .finally(() => setLoading(false))
  }, [])

  const filteredData = attractions.filter(item => {
    const matchesSearch = item.name.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesCity = selectedCityId ? item.city_id === parseInt(selectedCityId) : true
    const matchesCategory = selectedCategory ? item.category === selectedCategory : true
    const matchesRating = minRating ? item.rating >= parseFloat(minRating) : true
    return matchesSearch && matchesCity && matchesCategory && matchesRating
  })

  const columns = [
    { key: 'name', label: 'Attraction Name' },
    { 
      key: 'city_id', 
      label: 'City',
      render: (id) => cities.find(c => c.city_id === id)?.name || `ID: ${id}`
    },
    { 
      key: 'category', 
      label: 'Category',
      render: (cat) => (
        <span className={`px-2 py-1 rounded-full text-xs font-medium border ${CATEGORY_COLORS[cat] || 'bg-gray-100 text-gray-800 border-gray-200'}`}>
          {cat}
        </span>
      )
    },
    { key: 'rating', label: 'Rating', render: (val) => `⭐ ${val.toFixed(1)}` },
  ]

  return (
    <div>
      <PageHeader icon="🏛️" title="Attractions" subtitle="Explore tourist spots with custom filters" />

      <div className="mb-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 items-end">
        <div>
          <label className="label">Search Name</label>
          <input 
            type="text" className="input" placeholder="Attraction name..."
            value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)}
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
          <label className="label">Category</label>
          <select className="input" value={selectedCategory} onChange={(e) => setSelectedCategory(e.target.value)}>
            <option value="">All Categories</option>
            {Object.keys(CATEGORY_COLORS).map(cat => <option key={cat} value={cat}>{cat}</option>)}
          </select>
        </div>

        <div>
          <label className="label">Rating</label>
          <select className="input" value={minRating} onChange={(e) => setMinRating(e.target.value)}>
            <option value="">Rating</option>
            {[4.5, 4.0, 3.5, 3.0].map(r => <option key={r} value={r}>{r}+ Stars</option>)}
          </select>
        </div>
      </div>

      <div className="card">
        <DataTable columns={columns} data={filteredData} loading={loading} />
      </div>
    </div>
  )
}