import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../context/AuthContext'
import { login, verifyEmail } from '../api/auth'
import Swal from 'sweetalert2'

export default function Login() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const { t } = useTranslation()
  const [form, setForm] = useState({ username: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(form.username, form.password)
      navigate('/dashboard')
    } catch (err) {
      if (err.response?.status === 403 && err.response?.data?.status === 'unverified') {
        const emailForVerify = err.response.data.email || form.username;
        
        Swal.fire({
          title: 'Հաշիվը հաստատված չէ',
          text: `Կոդն ուղարկվել է ${form.email} հասցեին`,
          input: 'text',
          inputPlaceholder: '______',
          inputAttributes: { maxlength: 6, autofocus: true },
          confirmButtonText: 'Հաստատել',
          showLoaderOnConfirm: true,
          preConfirm: async (code) => {
            try {
              const res = await verifyEmail(emailForVerify, code);
              return res.data;
            } catch (err) {
              if (!err.response) {
                Swal.showValidationMessage("Կապի սխալ սերվերի հետ կամ CORS արգելափակում");
              } else {
                Swal.showValidationMessage(err.response.data?.message || "Սխալ կոդ");
              }
              return false;
            }
          },
          allowOutsideClick: () => !Swal.isLoading()
        }).then((result) => {
          if (result.isConfirmed) {
            Swal.fire({
            title: 'Հաջողված է',
            text: 'Ձեր էլ. հասցեն հաստատվեց։ Այժմ կարող եք մուտք գործել։',
            icon: 'success',
            confirmButtonText: 'Լավ'
          });
          }
        });
      } else {
        setError(err.response?.data || t('auth.login.error_invalid'))
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="w-full font-sans text-left">
      <div className="text-center mb-6">
          <span className="text-4xl">🗺️</span>
          <h1 className="text-2xl font-black text-gray-900 mt-2">
            {t('auth.login.title')}
          </h1>
          <p className="text-sm text-gray-500 font-medium">
            {t('auth.login.subtitle')}
          </p>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-2xl text-sm text-red-700 font-bold">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="text-[10px] font-black capitalize text-gray-400 tracking-widest mb-2 block ml-1">
              {t('auth.login.username_label')}
            </label>
            <input
              type="text"
              className="w-full p-4 bg-gray-50 border border-gray-100 rounded-2xl text-sm focus:ring-2 focus:ring-blue-500 outline-none transition-all"
              value={form.username}
              onChange={(e) => setForm({ ...form, username: e.target.value })}
              placeholder={t('auth.login.username_placeholder')}
              required
            />
          </div>
          <div>
            <label className="text-[10px] font-black capitalize text-gray-400 tracking-widest mb-2 block ml-1">
              {t('auth.login.password_label')}
            </label>
            <input
              type="password"
              className="w-full p-4 bg-gray-50 border border-gray-100 rounded-2xl text-sm focus:ring-2 focus:ring-blue-500 outline-none transition-all"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              placeholder={t('auth.login.password_placeholder')}
              required
            />
          </div>
          <button 
            type="submit" 
            disabled={loading} 
            className="w-full py-4 bg-blue-600 hover:bg-blue-700 text-white rounded-2xl font-black capitalize text-xs tracking-[0.2em] shadow-lg transition-all active:scale-95 disabled:opacity-50"
          >
            {loading ? t('auth.login.loading_btn') : t('auth.login.submit_btn')}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-gray-500 font-medium">
          {t('auth.login.no_account')}{' '}
          <Link to="/register" className="text-blue-600 hover:underline font-black capitalize text-xs tracking-tighter">
            {t('auth.login.register_link')}
          </Link>
        </p>
      </div>
  )
}