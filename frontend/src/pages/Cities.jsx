import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next' 
import { getCities, getCountries } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

export default function Cities() {
  const { t } = useTranslation(); 
  const [cities, setCities] = useState([])
  const [countries, setCountries] = useState([])
  const [selectedCountryId, setSelectedCountryId] = useState('')
  const [searchTerm, setSearchTerm] = useState('') 
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [lastUpdated, setLastUpdated] = useState(null)

  useEffect(() => {
    setLoading(true)
    Promise.all([getCities(), getCountries()])
      .then(([citiesData, countriesData]) => {
        if (citiesData && citiesData.length > 0) {
          const latest = citiesData.reduce((prev, current) => 
            (new Date(prev.updated_at) > new Date(current.updated_at)) ? prev : current
          );
          setLastUpdated(latest.updated_at);
        }

        const sortedCities = (citiesData || []).sort((a, b) => 
          a.name.localeCompare(b.name)
        );

        const sortedCountries = (countriesData || []).sort((a, b) => 
          a.name.localeCompare(b.name)
        );

        setCities(sortedCities);
        setCountries(sortedCountries);
      })
      .catch(err => {
        console.error("Error loading cities data:", err);
        setError(t('cities.error_load')); 
      })
      .finally(() => setLoading(false));
  }, [t]);

  const filteredCities = cities.filter(city => {
    const matchesCountry = selectedCountryId ? city.country_id === parseInt(selectedCountryId) : true
    const matchesSearch = city.name.toLowerCase().includes(searchTerm.toLowerCase())
    return matchesCountry && matchesSearch
  })

  const columns = [
    { key: 'name', label: t('cities.table.city_name') },
    { 
      key: 'country_id', 
      label: t('cities.table.country'),
      render: (id) => countries.find(c => c.country_id === id)?.name || `ID: ${id}`
    },
    { 
      key: 'description', 
      label: t('cities.table.description'),
      render: (text) => text || t('common.no_description') 
    },
  ]

  if (error) return <div className="p-8 text-red-600">{error}</div>

  return (
    <div>
      <PageHeader 
        icon="🏙️" 
        title={t('nav.cities')} 
        subtitle={t('cities.subtitle')}
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

      <div className="mb-6 flex flex-wrap gap-4 items-end bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <div className="w-full max-w-xs">
          <label className="label">{t('cities.filter.search_label')}</label>
          <input 
            type="text"
            className="input"
            placeholder={t('cities.filter.placeholder')}
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
      </div>

        <div className="w-full max-w-xs">
          <label className="label">{t('cities.filter.country')}</label>
          <select 
            className="input"
            value={selectedCountryId}
            onChange={(e) => setSelectedCountryId(e.target.value)}
          >
            <option value="">{t('cities.filter.all_countries')}</option>
            {countries.map(c => (
              <option key={c.country_id} value={c.country_id}>{c.name}</option>
            ))}
          </select>
        </div>
      </div>

      <div className="card shadow-md">
        <DataTable 
          columns={columns} 
          data={filteredCities} 
          loading={loading} 
        />
      </div>
    </div>
  )
}