import { useState, useEffect } from 'react'
import { getFlights } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

function formatDuration(minutes) {
  if (!minutes) return '—'
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return `${h}h ${m}m`
}

const columns = [
  { key: 'airline',          label: 'Airline' },
  { key: 'from_city_id',     label: 'From City' },
  { key: 'to_city_id',       label: 'To City' },
  { key: 'duration_minutes', label: 'Duration', render: (v) => formatDuration(v) },
  { key: 'price',            label: 'Price', render: (v) => v != null ? <span className="font-medium text-green-700">${Number(v).toLocaleString()}</span> : '—' },
]

export default function Flights() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getFlights()
      .then(setData)
      .catch(() => setError('Failed to load flights.'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div>
      <PageHeader icon="🛫" title="Flights" subtitle={loading ? 'Loading…' : `${data.length} flights`} />
      {loading && <p className="text-gray-400 py-4">Loading…</p>}
      {error && <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>}
      {!loading && !error && <DataTable columns={columns} data={data} searchKeys={['airline']} />}
    </div>
  )
}
