import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getCountries } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

export default function Countries() {
  const { t } = useTranslation();
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [searchTerm, setSearchTerm] = useState('')

  const columns = [
    { key: 'name', label: t('countries.table.country') },
    { 
      key: 'code', 
      label: t('countries.table.code'), 
      render: (v) => <span className="badge bg-gray-100 text-gray-700 font-mono">{v}</span> 
    },
  ]

  useEffect(() => {
    getCountries()
      .then((rawData) => {
        const sortedData = [...rawData].sort((a, b) => a.name.localeCompare(b.name));
        setData(sortedData);
      })
      .catch(() => setError(t('countries.error_load')))
      .finally(() => setLoading(false))
  }, [t])

  const filteredData = data.filter(country => 
    country.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    country.code.toLowerCase().includes(searchTerm.toLowerCase())
  )

  return (
    <div>
      <PageHeader 
        icon="🌍" 
        title={t('nav.countries')} 
        subtitle={t('countries.subtitle')}
        />

      <div className="mb-6 max-w-md">
        <label className="label">{t('countries.filter.search_label')}</label>
        <input 
          type="text"
          className="input"
          placeholder={t('countries.filter.placeholder')}
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>
      
      {error && <div className="p-3 bg-red-50 text-red-700 rounded-lg mb-4">{error}</div>}
      
      <div className="card shadow-md">
        <DataTable 
          columns={columns} 
          data={filteredData} 
          loading={loading}
        />
      </div>
    </div>
  )
}