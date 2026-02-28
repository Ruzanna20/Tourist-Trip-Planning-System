import { useState, useEffect } from 'react'
import { getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

const columns = [
  { key: 'name',        label: 'City' },
  { key: 'country_id',  label: 'Country ID' },
  { key: 'description', label: 'Description', render: (v) => v ? <span className="text-gray-600 text-xs">{v.slice(0, 80)}{v.length > 80 ? '…' : ''}</span> : '—' },
]

export default function Cities() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getCities()
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
      <PageHeader icon="🏙️" title="Cities" subtitle={loading ? 'Loading…' : `${data.length} cities`} />
      {loading && <p className="text-gray-400 py-4">Loading…</p>}
      {error && <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>}
      {!loading && !error && <DataTable columns={columns} data={data} searchKeys={['name', 'iata_code']} />}
    </div>
  )
}
