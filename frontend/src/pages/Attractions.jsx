import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getAttractions, getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

const CATEGORY_COLORS = {
  museum: 'bg-blue-100 text-blue-800 border-blue-200',
  viewpoint: 'bg-green-100 text-green-800 border-green-200',
  gallery: 'bg-purple-100 text-purple-800 border-purple-200',
  attraction: 'bg-orange-100 text-orange-800 border-orange-200',
  monument: 'bg-indigo-100 text-indigo-800 border-indigo-200',
  historic: 'bg-red-100 text-red-800 border-red-200',
}

export default function Attractions() {
  const { t } = useTranslation();
  const [attractions, setAttractions] = useState([])
  const [cities, setCities] = useState([])
  const [loading, setLoading] = useState(true)
  const [lastUpdated, setLastUpdated] = useState(null)
  
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCityId, setSelectedCityId] = useState('')
  const [selectedCategory, setSelectedCategory] = useState('')
  const [minRating, setMinRating] = useState('')

  useEffect(() => {
    setLoading(true)
    Promise.all([getAttractions(), getCities()])
      .then(([attractionsData, citiesData]) => {
        if (attractionsData && attractionsData.length > 0) {
          const latest = attractionsData.reduce((prev, current) => 
            (prev.updated_at > current.updated_at) ? prev : current
          );
          setLastUpdated(latest.updated_at);
        }

        const sortedAttractions = (attractionsData || []).sort((a, b) => 
          a.name.localeCompare(b.name)
        );
        const sortedCities = (citiesData || []).sort((a, b) => 
          a.name.localeCompare(b.name)
        );
        setAttractions(sortedAttractions);
        setCities(sortedCities);
      })
      .catch(err => console.error("Error loading data:", err))
      .finally(() => setLoading(false))
  }, [])

  const sortedCategories = Object.keys(CATEGORY_COLORS).sort((a, b) => a.localeCompare(b));

  const filteredData = attractions.filter(item => {
    const matchesSearch = item.name.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesCity = selectedCityId ? item.city_id === parseInt(selectedCityId) : true
    const matchesCategory = selectedCategory ? item.category === selectedCategory : true
    const matchesRating = minRating ? item.rating >= parseFloat(minRating) : true
    return matchesSearch && matchesCity && matchesCategory && matchesRating
  })

  const columns = [
    { 
      key: 'name', 
      label: t('attractions.table.name') },
    { 
      key: 'city_id', 
      label: t('attractions.table.cities'),
      render: (id) => cities.find(c => c.city_id === id)?.name || `ID: ${id}`
    },
    { 
      key: 'category', 
      label: t('attractions.table.category'),
      render: (cat) => (
      <span className={`px-2 py-1 rounded-full text-xs font-medium border ${CATEGORY_COLORS[cat] || 'bg-gray-100 text-gray-800 border-gray-200'}`}>
        {t(`attractions.categories.${cat}`)} 
      </span>
    )
    },
    { 
      key: 'rating', 
      label: t('attractions.table.rating'), 
      render: (val) => `⭐ ${val.toFixed(1)}` },
    { 
      key: 'website', 
      label: t('attractions.table.info'), 
      render: (url) => {
        if (!url) return <span className="text-gray-400 italic text-xs">{t('attractions.no_link')}</span>;
        const href = url.startsWith('http') ? url : `https://${url}`;
        return (
          <a 
            href={href} 
            target="_blank" 
            rel="noopener noreferrer" 
            className="inline-flex items-center px-3 py-1 bg-purple-50 text-purple-700 rounded-lg hover:bg-purple-100 transition-colors border border-purple-200 font-medium text-xs"
          >
            {t('attractions.view_website')} ↗
          </a>
        );
      }
    },
  ]

  return (
    <div>
      <PageHeader 
        icon="🏛️" 
        title={t('nav.attractions')} 
        subtitle={t('attractions.subtitle')}
      />

      {lastUpdated && (
          <div className="mb-4 md:mb-6 flex items-center gap-2 px-3 py-1.5 bg-white border border-slate-100 rounded-full shadow-sm w-fit">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            <span className="text-[10px] font-black uppercase tracking-widest text-slate-400">
              {t('common.last_updated')} {new Date(lastUpdated).toLocaleString()}
            </span>
          </div>
        )}

      <div className="mb-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 items-end bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <div>
          <label className="label">{t('attractions.filter.name')}</label>
          <input 
            type="text" className="input" placeholder={t('attractions.filter.placeholder')}
            value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        <div>
          <label className="label">{t('attractions.filter.cities')}</label>
          <select className="input" value={selectedCityId} onChange={(e) => setSelectedCityId(e.target.value)}>
            <option value="">{t('attractions.filter.all_cities')}</option>
            {cities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
          </select>
        </div>

        <div>
          <label className="label">{t('attractions.filter.category')}</label>
          <select className="input" value={selectedCategory} onChange={(e) => setSelectedCategory(e.target.value)}>
            <option value="">{t('attractions.filter.all_categories')}</option>
            {sortedCategories.map(cat => (
              <option key={cat} value={cat}>{t(`attractions.categories.${cat}`)}</option>
            ))}
          </select>
        </div>

        <div>
          <label className="label">{t('attractions.filter.rating')}</label>
          <select className="input" value={minRating} onChange={(e) => setMinRating(e.target.value)}>
            <option value="">{t('attractions.filter.all_ratings')}</option>
            {[4.5, 4.0, 3.5, 3.0].map(r => (
              <option key={r} value={r}>{r}+ {t('attractions.filter.stars')}</option>
            ))}
          </select>
        </div>
      </div>

      <div className="card shadow-md">
        <DataTable columns={columns} data={filteredData} loading={loading} />
      </div>
    </div>
  )
}