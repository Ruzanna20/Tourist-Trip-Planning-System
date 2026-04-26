import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getRestaurants, getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

export default function Restaurants() {
  const { t } = useTranslation();
  const [data, setData] = useState([])
  const [cities, setCities] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCityId, setSelectedCityId] = useState('')
  const [selectedCuisine, setSelectedCuisine] = useState('')
  const [minRating, setMinRating] = useState('')

  useEffect(() => {
    setLoading(true)
    Promise.all([getRestaurants(), getCities()])
      .then(([restData, citiesData]) => {
        const sortedData = [...restData].sort((a, b) => a.name.localeCompare(b.name));
        setData(sortedData);

        const sortedCities = (citiesData || []).sort((a, b) => 
          a.name.localeCompare(b.name)
        );
        setCities(sortedCities);
      })
      .catch(() => setError(t('restaurants.error_load')))
      .finally(() => setLoading(false))
  }, [t])

  const cuisines = [...new Set(data.map(r => r.cuisine).filter(Boolean))].sort((a, b) => 
    a.localeCompare(b)
  );

  const filteredData = data.filter(item => {
    const matchesSearch = item.name?.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesCity = selectedCityId ? item.city_id === parseInt(selectedCityId) : true
    const matchesCuisine = selectedCuisine ? item.cuisine === selectedCuisine : true
    const matchesRating = minRating ? item.rating >= parseFloat(minRating) : true
    return matchesSearch && matchesCity && matchesCuisine && matchesRating
  })

  const columns = [
    { 
      key: 'name', 
      label: t('restaurants.table.restaurant') },
    { 
      key: 'city_id', 
      label: t('restaurants.table.cities'),
      render: (id) => cities.find(c => c.city_id === id)?.name || `ID: ${id}`
    },
    { 
      key: 'cuisine', 
      label: t('restaurants.table.cuisine'), 
      render: (v) => v ? <span className="badge bg-orange-100 text-orange-700 capitalize">{v}</span> : '—' 
    },
    { 
      key: 'rating', 
      label: t('restaurants.table.rating'), 
      render: (val) => (
        <span className={`px-2 py-1 rounded text-xs font-bold ${val >= 4.5 ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'}`}>
          {val.toFixed(1)} / 5
        </span>
      )
    },
    { 
      key: 'price_range', 
      label: t('restaurants.table.price'), 
      render: (v) => <span className="text-gray-600 font-medium">{v || '—'}</span>
    },
    { 
      key: 'website', 
      label: t('restaurants.table.info'), 
      render: (url) => {
        if (!url) return <span className="text-gray-400 italic text-xs">{t('restaurants.no_link')}</span>;
        const href = url.startsWith('http') ? url : `https://${url}`;
        return (
          <a 
            href={href} 
            target="_blank" 
            rel="noopener noreferrer" 
            className="inline-flex items-center px-3 py-1 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 transition-colors border border-blue-200 font-medium text-xs"
          >
            {t('restaurants.visit_website')} ↗
          </a>
        );
      }
    },
  ]

  if (error) return <div className="p-8 text-red-600">{error}</div>

  return (
    <div>
      <PageHeader 
        icon="🍽️" 
        title={t('nav.restaurants')} 
        subtitle={t('restaurants.subtitle')}
      />
      
      <div className="mb-6 grid grid-cols-1 md:grid-cols-4 gap-4 items-end bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <div>
          <label className="label">{t('restaurants.filter.search_label')}</label>
          <input type="text" className="input" placeholder={t('restaurants.filter.placeholder')} value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} />
        </div>
        <div>
          <label className="label">{t('restaurants.filter.cities')}</label>
          <select className="input" value={selectedCityId} onChange={(e) => setSelectedCityId(e.target.value)}>
            <option value="">{t('restaurants.filter.all_cities')}</option>
            {cities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
          </select>
        </div>
        <div>
          <label className="label">{t('restaurants.filter.cuisine')}</label>
          <select className="input" value={selectedCuisine} onChange={(e) => setSelectedCuisine(e.target.value)}>
            <option value="">{t('restaurants.filter.all_cuisines')}</option>
            {cuisines.map(c => <option key={c} value={c}>{c}</option>)}
          </select>
        </div>
        <div>
          <label className="label">{t('restaurants.filter.rating')}</label>
          <select className="input" value={minRating} onChange={(e) => setMinRating(e.target.value)}>
            <option value="">{t('restaurants.filter.all_ratings')}</option>
            {[4.8, 4.5, 4.0, 3.5].map(r => <option key={r} value={r}>{r}+ {t('restaurants.filter.stars')}</option>)}
          </select>
        </div>
      </div>

      <div className="card shadow-md">
        <DataTable columns={columns} data={filteredData} loading={loading} />
      </div>
    </div>
  )
}