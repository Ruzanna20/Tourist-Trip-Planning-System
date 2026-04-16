import { useState, useEffect } from 'react'
import { getUserPreferences, setPreferences } from '../api/users'
import { getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'

const CATEGORIES = [
  { id: 'museum', label: 'Museum' },
  { id: 'viewpoint', label: 'Viewpoint' },
  { id: 'gallery', label: 'Gallery' },
  { id: 'attraction', label: 'Attraction' },
  { id: 'monument', label: 'Monument' },
  { id: 'historic', label: 'Historic' },
];

export default function Preferences() {
  const [cities, setCities] = useState([])
  const [form, setForm] = useState({
    home_city_id: '',
    preferred_categories: [], 
  })
  const [savedAt, setSavedAt] = useState(null)
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState('')

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

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await setPreferences({
        home_city_id: parseInt(form.home_city_id),
        preferred_categories: form.preferred_categories.join(','),
      })
      setSavedAt(new Date().toISOString())
    } catch (err) {
      setError(err.response?.data || 'Failed to save preferences.')
    } finally {
      setLoading(false)
    }
  }

  if (fetching) return <div className="p-8 text-gray-500 text-center">Loading preferences...</div>

  return (
    <div className="max-w-2xl mx-auto">
      <PageHeader
        icon="⚙️"
        title="Preferences"
        subtitle="Set your travel preferences to get better trip recommendations"
      />

      <div className="card mt-6">
        {savedAt && (
          <div className="mb-5 p-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-700">
            ✅ Last saved on {new Date(savedAt).toLocaleString()}
          </div>
        )}

        {error && <div className="mb-4 p-3 bg-red-50 text-red-700 rounded-lg">{error}</div>}

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Home City Section */}
          <div>
            <label className="label">Home City</label>
            <select 
              className="input" 
              value={form.home_city_id} 
              onChange={(e) => setForm({...form, home_city_id: e.target.value})} 
              required
            >
              <option value="">Select your home city…</option>
              {cities.map((c) => (
                <option key={c.city_id} value={c.city_id}>{c.name}</option>
              ))}
            </select>
          </div>

          {/* Categories Section */}
          <div>
            <label className="label mb-3">Preferred Categories</label>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
              {CATEGORIES.map((cat) => (
                <button
                  key={cat.id}
                  type="button"
                  onClick={() => handleCategoryChange(cat.id)}
                  className={`p-3 rounded-lg border text-sm transition-all ${
                    form.preferred_categories.includes(cat.id)
                      ? 'border-brand-600 bg-brand-50 text-brand-700 shadow-sm'
                      : 'border-gray-200 bg-white text-gray-600 hover:border-gray-300'
                  }`}
                >
                  {cat.label}
                </button>
              ))}
            </div>
            <p className="text-xs text-gray-400 mt-2">Select the types of places you enjoy visiting</p>
          </div>

          <button type="submit" disabled={loading} className="btn-primary w-full py-3">
            {loading ? 'Saving...' : 'Save Preferences'}
          </button>
        </form>
      </div>
    </div>
  )
}