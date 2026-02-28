import { useState, useEffect } from 'react'
import { getUserPreferences, setPreferences } from '../api/users'
import { getCities } from '../api/resources'
import PageHeader from '../components/PageHeader'

export default function Preferences() {
  const [cities, setCities] = useState([])
  const [form, setForm] = useState({
    home_city_id: '',
    budget_min: '',
    budget_max: '',
    preferred_categories: '',
  })
  const [savedAt, setSavedAt] = useState(null)
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState('')

  // Load cities and existing preferences on mount
  useEffect(() => {
    getCities().then(setCities).catch(() => {})

    getUserPreferences()
      .then((prefs) => {
        if (prefs) {
          setForm({
            home_city_id: prefs.home_city_id ? String(prefs.home_city_id) : '',
            budget_min: prefs.budget_min ? String(prefs.budget_min) : '',
            budget_max: prefs.budget_max ? String(prefs.budget_max) : '',
            preferred_categories: prefs.preferred_categories ?? '',
          })
          setSavedAt(prefs.updated_at)
        }
      })
      .catch(() => {})
      .finally(() => setFetching(false))
  }, [])

  const set = (field) => (e) => setForm({ ...form, [field]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await setPreferences({
        home_city_id: parseInt(form.home_city_id),
        budget_min: parseFloat(form.budget_min),
        budget_max: parseFloat(form.budget_max),
        travel_style: '',
        preferred_categories: form.preferred_categories,
      })
      setSavedAt(new Date().toISOString())
    } catch (err) {
      setError(err.response?.data || 'Failed to save preferences.')
    } finally {
      setLoading(false)
    }
  }

  if (fetching) {
    return (
      <div>
        <PageHeader icon="⚙️" title="Preferences" />
        <div className="flex items-center gap-2 text-gray-400 py-8">
          <div className="w-4 h-4 border-2 border-gray-300 border-t-brand-600 rounded-full animate-spin" />
          Loading preferences…
        </div>
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        icon="⚙️"
        title="Preferences"
        subtitle="Set your travel preferences to get better trip recommendations"
      />

      <div className="card max-w-lg">
        {/* Last saved timestamp */}
        {savedAt && (
          <div className="mb-5 flex items-center gap-2 p-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-700">
            <span>✅</span>
            <span>
              Preferences last saved on{' '}
              {new Date(savedAt).toLocaleDateString('en-US', {
                year: 'numeric', month: 'long', day: 'numeric',
                hour: '2-digit', minute: '2-digit',
              })}
            </span>
          </div>
        )}

        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="label">Home City</label>
            <select className="input" value={form.home_city_id} onChange={set('home_city_id')} required>
              <option value="">Select your home city…</option>
              {cities.map((c) => (
                <option key={c.city_id} value={c.city_id}>
                  {c.name} {c.iata_code ? `(${c.iata_code})` : ''}
                </option>
              ))}
            </select>
            <p className="text-xs text-gray-400 mt-1">Used as the departure city when planning trips</p>
          </div>

          <div>
            <label className="label">Budget Range (USD)</label>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <input
                  type="number"
                  className="input"
                  value={form.budget_min}
                  onChange={set('budget_min')}
                  min="0"
                  placeholder="Min — e.g. 500"
                  required
                />
                <p className="text-xs text-gray-400 mt-1">Minimum budget</p>
              </div>
              <div>
                <input
                  type="number"
                  className="input"
                  value={form.budget_max}
                  onChange={set('budget_max')}
                  min="0"
                  placeholder="Max — e.g. 5000"
                  required
                />
                <p className="text-xs text-gray-400 mt-1">Maximum budget</p>
              </div>
            </div>
          </div>

          <div>
            <label className="label">Preferred Categories</label>
            <input
              type="text"
              className="input"
              value={form.preferred_categories}
              onChange={set('preferred_categories')}
              placeholder="e.g. museums, beaches, hiking, food"
            />
            <p className="text-xs text-gray-400 mt-1">Comma-separated list of interests</p>
          </div>

          <button type="submit" disabled={loading} className="btn-primary w-full justify-center">
            {loading ? 'Saving…' : 'Save Preferences'}
          </button>
        </form>
      </div>
    </div>
  )
}
