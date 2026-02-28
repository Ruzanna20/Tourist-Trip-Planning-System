import { useState, useEffect } from 'react'
import { getHotels } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

const columns = [
  { key: 'name',           label: 'Hotel' },
  { key: 'stars',          label: 'Stars', render: (v) => '⭐'.repeat(v ?? 0) || '—' },
  { key: 'rating',         label: 'Rating', render: (v) => v ? `${v}/10` : '—' },
  { key: 'price_per_night',label: 'Price/Night', render: (v) => v != null ? `$${Number(v).toLocaleString()}` : '—' },
  { key: 'address',        label: 'Address', render: (v) => <span className="text-gray-500 text-xs">{v || '—'}</span> },
]

export default function Hotels() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getHotels()
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
      <PageHeader icon="🏨" title="Hotels" subtitle={loading ? 'Loading…' : `${data.length} hotels`} />
      {loading && <p className="text-gray-400 py-4">Loading…</p>}
      {error && <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>}
      {!loading && !error && <DataTable columns={columns} data={data} searchKeys={['name', 'address']} />}
    </div>
  )
}
