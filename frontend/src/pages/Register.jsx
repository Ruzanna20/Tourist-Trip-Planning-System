import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { registerUser } from '../api/users'

export default function Register() {
  const navigate = useNavigate()
  const { t } = useTranslation()
  const [form, setForm] = useState({ first_name: '', last_name: '', email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const set = (field) => (e) => setForm({ ...form, [field]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await registerUser(form)
      navigate('/login')
    } catch (err) {
      setError(err.response?.data || t('auth.register.error_failed'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 font-sans text-left p-4">
      <div className="card w-full max-w-sm p-8 bg-white rounded-[32px] shadow-xl">
        <div className="text-center mb-6">
          <span className="text-4xl">🗺️</span>
          <h1 className="text-2xl font-black text-gray-900 mt-2">
            {t('auth.register.title')}
          </h1>
          <p className="text-sm text-gray-500 font-medium">
            {t('auth.register.subtitle')}
          </p>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-2xl text-sm text-red-700 font-bold italic text-left">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-[10px] font-black uppercase text-gray-400 tracking-widest mb-1 block ml-1">
                {t('auth.register.first_name_label')}
              </label>
              <input 
                type="text" 
                className="w-full p-4 bg-gray-50 border border-gray-100 rounded-2xl text-sm outline-none focus:ring-2 focus:ring-blue-500 transition-all" 
                value={form.first_name} 
                onChange={set('first_name')} 
                required 
              />
            </div>
            <div>
              <label className="text-[10px] font-black capitalize text-gray-400 tracking-widest mb-1 block ml-1">
                {t('auth.register.last_name_label')}
              </label>
              <input 
                type="text" 
                className="w-full p-4 bg-gray-50 border border-gray-100 rounded-2xl text-sm outline-none focus:ring-2 focus:ring-blue-500 transition-all" 
                value={form.last_name} 
                onChange={set('last_name')} 
                required 
              />
            </div>
          </div>
          <div>
            <label className="text-[10px] font-black capitalize text-gray-400 tracking-widest mb-1 block ml-1">
              {t('auth.register.email_label')}
            </label>
            <input 
              type="email" 
              className="w-full p-4 bg-gray-50 border border-gray-100 rounded-2xl text-sm outline-none focus:ring-2 focus:ring-blue-500 transition-all" 
              value={form.email} 
              onChange={set('email')} 
              required 
            />
          </div>
          <div>
            <label className="text-[10px] font-black capitalize text-gray-400 tracking-widest mb-1 block ml-1">
              {t('auth.register.password_label')}
            </label>
            <input 
              type="password" 
              className="w-full p-4 bg-gray-50 border border-gray-100 rounded-2xl text-sm outline-none focus:ring-2 focus:ring-blue-500 transition-all" 
              value={form.password} 
              onChange={set('password')} 
              required 
            />
          </div>
          <button 
            type="submit" 
            disabled={loading} 
            className="w-full py-4 bg-blue-600 hover:bg-blue-700 text-white rounded-2xl font-black capitalize text-xs tracking-[0.2em] shadow-lg transition-all active:scale-95 disabled:opacity-50 mt-2"
          >
            {loading ? t('auth.register.loading_btn') : t('auth.register.submit_btn')}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-gray-500 font-medium">
          {t('auth.register.already_account')}{' '}
          <Link to="/login" className="text-blue-600 hover:underline font-black text-xs tracking-tighter">
            {t('auth.register.login_link')}
          </Link>
        </p>
      </div>
    </div>
  )
}