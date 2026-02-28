import { useState, useEffect } from 'react'
import { getRestaurants } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

const columns = [
  { key: 'name',          label: 'Restaurant' },
  { key: 'cuisine',       label: 'Cuisine', render: (v) => v ? <span className="badge bg-orange-100 text-orange-700 capitalize">{v}</span> : '—' },
  { key: 'rating',        label: 'Rating', render: (v) => v ? `${v}/10` : '—' },
  { key: 'price_range',   label: 'Price Range', render: (v) => v || '—' },
  { key: 'city_id',       label: 'City ID' },
]

export default function Restaurants() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getRestaurants()
      .then((rawData) => {
        const sortedData = [...rawData].sort((a, b) => 
          a.name.localeCompare(b.name)
        );
        setData(sortedData);
      })
      .catch(() => setError('Failed to load attractions.'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div>
      <PageHeader icon="🍽️" title="Restaurants" subtitle={loading ? 'Loading…' : `${data.length} restaurants`} />
      {loading && <p className="text-gray-400 py-4">Loading…</p>}
      {error && <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>}
      {!loading && !error && <DataTable columns={columns} data={data} searchKeys={['name', 'cuisine']} />}
    </div>
  )
}
