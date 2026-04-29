import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getHotels, getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

export default function Hotels() {
  const { t } = useTranslation();
  const [hotels, setHotels] = useState([])
  const [cities, setCities] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCityId, setSelectedCityId] = useState('')
  const [selectedStars, setSelectedStars] = useState('')
  const [minRating, setMinRating] = useState('')
  const [maxPrice, setMaxPrice] = useState('')
  const [lastUpdated, setLastUpdated] = useState(null)

  useEffect(() => {
    setLoading(true)
    Promise.all([getHotels(), getCities()])
      .then(([hotelsData, citiesData]) => {
        if (hotelsData && hotelsData.length > 0) {
          const latest = hotelsData.reduce((prev, current) => 
            (new Date(prev.updated_at) > new Date(current.updated_at)) ? prev : current
          );
          setLastUpdated(latest.updated_at);
        }

        const sortedHotels = (hotelsData || []).sort((a, b) => a.name.localeCompare(b.name));
        const sortedCities = (citiesData || []).sort((a, b) => a.name.localeCompare(b.name));
        setHotels(sortedHotels);
        setCities(sortedCities);
      })
      .catch(err => {
        console.error("Error loading hotels:", err);
        setError(t('hotels.error_load'));
      })
      .finally(() => setLoading(false));
  }, [t]);

  const filteredHotels = hotels.filter(hotel => {
    const matchesSearch = hotel.name.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesCity = selectedCityId ? hotel.city_id === parseInt(selectedCityId) : true
    const matchesStars = selectedStars ? hotel.stars === parseInt(selectedStars) : true
    const matchesRating = minRating ? hotel.rating >= parseFloat(minRating) : true
    const matchesPrice = maxPrice ? hotel.price_per_night <= parseFloat(maxPrice) : true
    return matchesSearch && matchesCity && matchesStars && matchesRating && matchesPrice
  })

  const columns = [
    { key: 'name', label: t('hotels.table.name') },
    { 
      key: 'city_id', 
      label: t('hotels.table.cities'),
      render: (id) => cities.find(c => c.city_id === id)?.name || `ID: ${id}`
    },
    { 
      key: 'address', 
      label: t('hotels.table.address'), 
      render: (val) => <span className="text-gray-600 text-sm italic">{val || '—'}</span>
    },
    { 
      key: 'description', 
      label: t('hotels.table.description'), 
      render: (val) => (
        <div className="max-w-xs truncate text-gray-500 text-xs" title={val}>
          {val || t('hotels.no_description')}
        </div>
      )
    },
    { 
      key: 'stars', 
      label: t('hotels.table.stars'), 
      render: (val) => (
        <span className="text-yellow-500 font-bold whitespace-nowrap">
          {'★'.repeat(val)}{'☆'.repeat(5 - val)}
        </span>
      )
    },
    { 
      key: 'rating', 
      label: t('hotels.table.rating'), 
      render: (val) => (
        <span className={`px-2 py-1 rounded text-xs font-bold ${val >= 8 ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'}`}>
          {val.toFixed(1)} / 10
        </span>
      )
    },
    { 
      key: 'price_per_night', 
      label: t('hotels.table.price'), 
      render: (val) => <span className="font-semibold text-brand-700">${val.toLocaleString()}</span>
    },
    { 
      key: 'website', 
      label: t('hotels.table.info'), 
      render: (url) => {
        if (!url) return <span className="text-gray-400 italic text-xs">{t('hotels.no_info')}</span>;
        const href = url.startsWith('http') ? url : `https://${url}`;
        return (
          <a 
            href={href} 
            target="_blank" 
            rel="noopener noreferrer" 
            className="inline-flex items-center px-3 py-1 bg-amber-50 text-amber-700 rounded-lg hover:bg-amber-100 transition-colors border border-amber-200 font-medium text-xs"
          >
            {t('hotels.view_website')} ↗
          </a>
        );
      }
    },
  ]

  if (error) return <div className="p-8 text-red-600 font-medium">{error}</div>

  return (
    <div>
      <PageHeader 
        icon="🏨" 
        title={t('nav.hotels')} 
        subtitle={t('hotels.subtitle')}
      />

      {lastUpdated && (
          <div className="mb-4 md:mb-6 flex items-center gap-2 px-3 py-1.5 bg-white border border-slate-100 rounded-full shadow-sm w-fit transition-all hover:border-slate-200">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            <span className="text-[10px] font-black text-slate-400">
              {t('common.last_updated')} {new Date(lastUpdated).toLocaleString()}
            </span>
          </div>
        )}

      <div className="mb-6 grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4 items-end bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <div>
          <label className="label">{t('hotels.filter.search_label')}</label>
          <input 
            type="text" className="input" placeholder={t('hotels.filter.placeholder')}
            value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        <div>
          <label className="label">{t('hotels.filter.cities')}</label>
          <select className="input" value={selectedCityId} onChange={(e) => setSelectedCityId(e.target.value)}>
            <option value="">{t('hotels.filter.all_cities')}</option>
            {cities.map(c => <option key={c.city_id} value={c.city_id}>{c.name}</option>)}
          </select>
        </div>

        <div>
          <label className="label">{t('hotels.filter.stars')}</label>
          <select className="input" value={selectedStars} onChange={(e) => setSelectedStars(e.target.value)}>
            <option value="">{t('hotels.filter.all_stars')}</option>
            {[5, 4, 3, 2, 1].map(s => <option key={s} value={s}>{s} {t('hotels.filter.stars_label')}</option>)}
          </select>
        </div>

        <div>
          <label className="label">{t('hotels.filter.rating')}</label>
          <select className="input" value={minRating} onChange={(e) => setMinRating(e.target.value)}>
            <option value="">{t('hotels.filter.all_ratings')}</option>
            {[9, 8, 7, 6].map(r => <option key={r} value={r}>{r}+ {t('hotels.filter.rating_desc')}</option>)}
          </select>
        </div>

        <div>
          <label className="label">{t('hotels.filter.max_price')}</label>
          <input 
            type="number" className="input" placeholder={t('hotels.filter.price_placeholder')}
            value={maxPrice} onChange={(e) => setMaxPrice(e.target.value)}
          />
        </div>
      </div>

      <div className="card shadow-lg">
        <DataTable columns={columns} data={filteredHotels} loading={loading} />
      </div>
    </div>
  )
}