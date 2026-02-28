import { useState, useEffect } from 'react'
import { getCountries } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

const columns = [
  { key: 'name',       label: 'Country' },
  { key: 'code',       label: 'Code', render: (v) => <span className="badge bg-gray-100 text-gray-700 font-mono">{v}</span> },
]

export default function Countries() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getCountries()
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
      <PageHeader icon="🌍" title="Countries" subtitle={loading ? 'Loading…' : `${data.length} countries`} />
      {loading && <p className="text-gray-400 py-4">Loading…</p>}
      {error && <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>}
      {!loading && !error && <DataTable columns={columns} data={data} searchKeys={['name', 'code']} />}
    </div>
  )
}
