import { useState, useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { getUserPreferences, setPreferences } from '../api/users'
import { getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'
import Swal from 'sweetalert2'

export default function Preferences() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const planRef = useRef(null);
  
  const CATEGORIES = [
    { id: 'museum', label: t('attractions.categories.museum') },
    { id: 'viewpoint', label: t('attractions.categories.viewpoint') },
    { id: 'gallery', label: t('attractions.categories.gallery') },
    { id: 'attraction', label: t('attractions.categories.attraction') },
    { id: 'monument', label: t('attractions.categories.monument') },
    { id: 'historic', label: t('attractions.categories.historic') },
  ];

  const [cities, setCities] = useState([])
  const [form, setForm] = useState({
    home_city_id: '',
    preferred_categories: [], 
  })
  const [savedAt, setSavedAt] = useState(null)
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState('')
  const [showPlanOption, setShowPlanOption] = useState(false)

  useEffect(() => {
    getCities().then(setCities).catch(() => {})

    getUserPreferences()
      .then((prefs) => {
        if (prefs) {
          setForm({
            home_city_id: prefs.home_city_id ? String(prefs.home_city_id) : '',
            preferred_categories: prefs.preferred_categories ? prefs.preferred_categories.split(',') : [],
          })
          setSavedAt(prefs.updated_at)
        }
      })
      .catch(() => {})
      .finally(() => setFetching(false))
  }, [])

  const handleCategoryChange = (categoryId) => {
    setForm(prev => {
      const current = prev.preferred_categories;
      const updated = current.includes(categoryId)
        ? current.filter(id => id !== categoryId)
        : [...current, categoryId];
      return { ...prev, preferred_categories: updated };
    });
  };

  const formatDateTime = (isoString) => {
    if (!isoString) return '';
    
    const date = new Date(isoString);
    if (isNaN(date.getTime())) return isoString;

    const day = String(date.getUTCDate()).padStart(2, '0');
    const month = String(date.getUTCMonth() + 1).padStart(2, '0');
    const year = date.getUTCFullYear();
    
    const hours = String(date.getUTCHours()).padStart(2, '0');
    const minutes = String(date.getUTCMinutes()).padStart(2, '0');
    const seconds = String(date.getUTCSeconds()).padStart(2, '0');

    return `${day}.${month}.${year}, ${hours}:${minutes}:${seconds}`;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    setShowPlanOption(false); 

    try {
      await setPreferences({
        home_city_id: parseInt(form.home_city_id),
        preferred_categories: form.preferred_categories.join(','),
      });

      const now = new Date().toISOString();
      setSavedAt(now);
      setShowPlanOption(true);

      setTimeout(() => {
        planRef.current?.scrollIntoView({ 
          behavior: 'smooth', 
          block: 'center' 
        });
      }, 200);

      Swal.fire({
        icon: 'success',
        title: t('preferences.saved_success'),
        showConfirmButton: false,
        timer: 1500
      });

    } catch (err) {      
      const errorMessage = err.response?.data?.error || err.response?.data || t('preferences.error_save');
      setError(errorMessage);
      setShowPlanOption(false); 

      Swal.fire({
        icon: 'error',
        title: 'Սխալ',
        text: errorMessage,
      });
    } finally {
      setLoading(false);
    }
  };

  if (fetching) return <div className="p-8 text-gray-500 text-center">{t('common.loading')}</div>

  return (
    <div className="max-w-2xl mx-auto pb-12">
      <PageHeader
        icon="⚙️"
        title={t('nav.preferences')}
        subtitle={t('preferences.subtitle')}
      />

      <div className="card mt-6 shadow-md border-t-4 border-blue-500">
        {savedAt && !error && (
          <div className="mb-5 p-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-700">
            ✅ {t('preferences.last_saved')} {formatDateTime(savedAt)}
          </div>
        )}

        {error && <div className="mb-4 p-3 bg-red-50 text-red-700 rounded-lg border-l-4 border-red-500 font-medium">{error}</div>}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="label">{t('preferences.home_city_label')}</label>
            <select 
              className="input" 
              value={form.home_city_id} 
              onChange={(e) => setForm({...form, home_city_id: e.target.value})} 
              required
            >
              <option value="">{t('preferences.select_city_placeholder')}</option>
              {cities.map((c) => (
                <option key={c.city_id} value={c.city_id}>{c.name}</option>
              ))}
            </select>
          </div>

          <div>
            <label className="label mb-3">{t('preferences.categories_label')}</label>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                {CATEGORIES.map((cat) => (
                  <button
                    key={cat.id}
                    type="button"
                    onClick={() => handleCategoryChange(cat.id)}
                    className={`flex items-center justify-center text-center px-4 py-3 min-h-[64px] rounded-xl border text-sm transition-all duration-200 ${
                      form.preferred_categories.includes(cat.id)
                        ? 'border-blue-600 bg-blue-50 text-blue-700 shadow-sm font-semibold ring-1 ring-blue-600'
                        : 'border-gray-200 bg-white text-gray-600 hover:border-gray-300 hover:bg-gray-50'
                    }`}
                  >
                    {cat.label}
                  </button>
                ))}
              </div>
              <p className="text-xs text-gray-400 mt-3">{t('preferences.categories_help')}</p>
            </div>

          <div className="flex justify-center w-full pt-8 pb-4">
            <button 
              type="submit" 
              disabled={loading} 
              className="bg-blue-600 text-white px-8 py-3 rounded-xl font-bold hover:bg-blue-700 transition-all shadow-md active:scale-95"
            >
              {loading ? t('common.loading') : t('preferences.save_btn')}
            </button>
          </div>
        </form>

        {showPlanOption && (
          <div ref={planRef} className="mt-8 pt-8 border-t border-gray-100 animate-fade-in">
            <div className="bg-blue-50 p-6 rounded-2xl border border-blue-100 text-center space-y-4 shadow-lg">
              <p className="text-blue-600 italic">Ձեր նախընտրությունները պահպանված են: Հիմա կարող եք պլանավորել ձեր ուղևորությունը:</p>
              <button 
                onClick={() => {
                  window.scrollTo(0, 0);
                  navigate('/trips/create')
                }}
                className="bg-blue-600 text-white px-8 py-3 rounded-xl font-bold hover:bg-blue-700 transition-all shadow-md active:scale-95"
              >
              {t('nav.plan_trip')}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}