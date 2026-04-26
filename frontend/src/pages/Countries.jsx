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
  const [lastUpdated, setLastUpdated] = useState(null)

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
        if (rawData && rawData.length > 0) {
          const latest = rawData.reduce((prev, current) => 
            (new Date(prev.updated_at) > new Date(current.updated_at)) ? prev : current
          );
          setLastUpdated(latest.updated_at);
        }

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
      {lastUpdated && (
          <div className="mb-4 md:mb-6 flex items-center gap-2 px-3 py-1.5 bg-white border border-slate-100 rounded-full shadow-sm w-fit transition-all hover:border-slate-200">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            <span className="text-[10px] font-black uppercase tracking-widest text-slate-400">
              {t('common.last_updated')} {new Date(lastUpdated).toLocaleString()}
            </span>
          </div>
        )}

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