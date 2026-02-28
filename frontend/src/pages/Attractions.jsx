import { useState, useEffect } from 'react'
import { getAttractions } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

const columns = [
  { key: 'name',          label: 'Attraction' },
  { key: 'category',      label: 'Category', render: (v) => v ? <span className="badge bg-purple-100 text-purple-700 capitalize">{v}</span> : '—' },
  { key: 'rating',        label: 'Rating', render: (v) => v ? `${v}/10` : '—' },
  { key: 'entry_fee',     label: 'Entry Fee', render: (v) => v != null ? (v === 0 ? 'Free' : `$${Number(v).toLocaleString()}`) : '—' },
  { key: 'city_id',       label: 'City ID' },
]

export default function Attractions() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getAttractions()
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
      <PageHeader icon="🎡" title="Attractions" subtitle={loading ? 'Loading…' : `${data.length} attractions`} />
      {loading && <p className="text-gray-400 py-4">Loading…</p>}
      {error && <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>}
      {!loading && !error && <DataTable columns={columns} data={data} searchKeys={['name', 'category']} />}
    </div>
  )
}
