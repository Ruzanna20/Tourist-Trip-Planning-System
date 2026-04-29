import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { getUserReviews, createReview, deleteReview } from '../api/reviews'
import { getVisitedEntities } from '../api/resources' 
import PageHeader from '../components/PageHeader'
import Swal from 'sweetalert2'

function Stars({ rating }) {
  return (
    <div className="flex gap-0.5">
      {Array.from({ length: 5 }, (_, i) => (
        <span key={i} className={i < rating ? 'text-yellow-400' : 'text-gray-200'}>★</span>
      ))}
    </div>
  )
}

const StarRatingInput = ({ rating, setRating }) => {
  return (
    <div className="flex gap-2 text-4xl">
      {[1, 2, 3, 4, 5].map((star) => (
        <button
          key={star}
          type="button"
          onClick={() => setRating(star)}
          className={`transition-all duration-150 transform hover:scale-110 active:scale-95 ${
            star <= rating ? 'text-yellow-400' : 'text-gray-200'
          }`}
        >
          ★
        </button>
      ))}
    </div>
  );
};

export default function Reviews() {
  const { t } = useTranslation();
  const [entityType, setEntityType] = useState('hotel')
  const [entities, setEntities] = useState([])
  const [form, setForm] = useState({ entity_id: '', rating: 5, comment: '' })
  const [reviews, setReviews] = useState([])
  const [loadingReviews, setLoadingReviews] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  const ENTITY_TYPES = [
    { value: 'hotel',      label: t('nav.hotels'),      icon: '🏨' },
    { value: 'attraction', label: t('nav.attractions'), icon: '🎡' },
    { value: 'restaurant', label: t('nav.restaurants'), icon: '🍽️' },
  ]

  const fetchReviews = useCallback(() => {
    setLoadingReviews(true)
    getUserReviews()
      .then((data) => setReviews(Array.isArray(data) ? data : []))
      .catch(() => {})
      .finally(() => setLoadingReviews(false))
  }, [])

  useEffect(() => {
    fetchReviews()
  }, [fetchReviews])

  useEffect(() => {
    setForm((f) => ({ ...f, entity_id: '' }))
    setError('')
    getVisitedEntities(entityType)
      .then((data) => setEntities(Array.isArray(data) ? data : []))
      .catch(() => setEntities([]))
  }, [entityType])

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!form.entity_id || form.entity_id === "0" || Number(form.entity_id) === 0) {
      setError(t('reviews.error_select_place')); 
      return;
    }
    setError('')
    setSubmitting(true)
    try {
      await createReview({
        entity_type: entityType,
        entity_id: Number(form.entity_id), 
        rating: form.rating,
        comment: form.comment,
      })

      Swal.fire({
      title: t('reviews.success_msg'),
      icon: 'success',
      confirmButtonText: 'Լավ',
      confirmButtonColor: '#2563eb', 
      borderRadius: '1rem',
    })

      setForm({ entity_id: '', rating: 5, comment: '' })
      fetchReviews()
    } catch (err) {
      const msg = err.response?.data?.error || t('reviews.error_submit')

      Swal.fire({
      title: 'Սխալ',
      text: msg,
      icon: 'error',
      confirmButtonText: 'Փակել'
    })
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async (id) => {
    Swal.fire({
      title: t('reviews.confirm_delete'),
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#ef4444', 
      cancelButtonColor: '#6b7280', 
      confirmButtonText: 'Այո, ջնջել',
      cancelButtonText: 'Չեղարկել'
    }).then(async (result) => {
      if (result.isConfirmed) {
        try {
          await deleteReview(id)
          fetchReviews()
          Swal.fire('Ջնջված է', '', 'success')
        } catch (err) {
          Swal.fire('Սխալ', t('reviews.error_delete'), 'error')
        }
      }
    })
  }

  return (
    <div className="max-w-4xl mx-auto space-y-12 pb-12">
      <PageHeader icon="⭐" title={t('nav.write_review')} subtitle={t('reviews.subtitle')} />
      
      <div className="card max-w-xl mx-auto shadow-xl border-gray-100 p-8">
        <form onSubmit={handleSubmit} className="space-y-8">
          
          <div className="space-y-3">
            <label className="label text-gray-700 font-bold">{t('reviews.category')}</label>
            <div className="grid grid-cols-3 gap-4">
              {ENTITY_TYPES.map((type) => (
                <button
                  key={type.value}
                  type="button"
                  onClick={() => setEntityType(type.value)}
                  className={`flex flex-col items-center justify-center gap-2 py-5 rounded-2xl border-2 transition-all duration-200 ${
                    entityType === type.value 
                    ? 'border-blue-500 bg-blue-50 text-blue-700 shadow-sm ring-1 ring-blue-500' 
                    : 'border-gray-100 bg-white text-gray-400 hover:border-gray-200 hover:bg-gray-50'
                  }`}
                >
                  <span className="text-3xl">{type.icon}</span>
                  <span className="text-xs font-bold">{type.label}</span>
                </button>
              ))}
            </div>
          </div>

          <div className="space-y-3">
            <label className="label text-gray-700 font-bold">{t('reviews.select_place')}</label>
            <select 
              className="input h-14 text-base" 
              value={form.entity_id} 
              onChange={(e) => setForm({...form, entity_id: e.target.value})} 
              required
            >
              <option value="">{t('reviews.choose_option')}</option>
              {entities.map((e) => {

                const actualId = e.id || e.hotel_id || e.attraction_id || e.restaurant_id;
                
                return (
                  <option key={actualId || e.name} value={actualId}>
                    {e.name}
                  </option>
                );
              })}
            </select>
          </div>

          <div className="space-y-3">
            <label className="label text-gray-700 font-bold">{t('reviews.rating')}</label>
            <div className="bg-gray-50/80 p-6 rounded-2xl border border-gray-100 flex flex-col items-center gap-3">
              <StarRatingInput 
                rating={form.rating} 
                setRating={(val) => setForm({...form, rating: val})} 
              />
              <span className="text-lg font-black text-blue-600">{form.rating} / 5</span>
            </div>
          </div>

          <div className="space-y-3">
            <label className="label text-gray-700 font-bold">{t('reviews.comment')}</label>
            <textarea 
              className="input min-h-[140px] py-4 text-base" 
              placeholder={t('reviews.comment_placeholder')} 
              value={form.comment} 
              onChange={(e) => setForm({...form, comment: e.target.value})} 
            />
          </div>

          <button 
            type="submit" 
            disabled={submitting} 
            className="btn-primary w-full py-4 text-lg font-bold shadow-lg hover:shadow-blue-200 transition-all active:scale-[0.98]"
          >
            {submitting ? t('common.loading') : t('reviews.submit')}
          </button>
          
          {error && <p className="text-red-500 text-sm text-center font-semibold bg-red-50 py-2 rounded-lg">{error}</p>}
        </form>
      </div>

      <div className="space-y-6">
        <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-3">
          <span>📋</span> {t('reviews.my_reviews')}
          {!loadingReviews && (
            <span className="text-sm font-medium bg-gray-100 px-3 py-1 rounded-full text-gray-500">
              {reviews.length}
            </span>
          )}
        </h2>

        {loadingReviews ? (
          <div className="flex items-center justify-center py-20 text-gray-400 gap-3">
            <div className="w-6 h-6 border-3 border-gray-200 border-t-blue-600 rounded-full animate-spin" />
            <span className="font-medium">{t('common.loading')}</span>
          </div>
        ) : reviews.length === 0 ? (
          <div className="card text-center py-16 bg-gray-50/50 border-dashed border-2 border-gray-200">
            <p className="text-5xl mb-4">📝</p>
            <p className="text-gray-500 font-medium">{t('reviews.no_reviews')}</p>
          </div>
        ) : (
          <div className="grid gap-6 sm:grid-cols-2">
            {reviews.map((rev) => {
              const typeInfo = ENTITY_TYPES.find((t) => t.value === rev.entity_type)
              return (
                <div key={rev.review_id} className="card p-6 relative group hover:shadow-xl transition-all border-gray-100 min-h-[160px]">
                  <button 
                    onClick={() => handleDelete(rev.review_id)}
                    className="absolute top-4 right-4 text-gray-300 hover:text-red-500 transition-colors p-2 rounded-lg hover:bg-red-50"
                  >
                    🗑️
                  </button>

                  <div className="flex flex-col gap-4">
                    <div className="flex items-center gap-3">
                      <div className="flex-shrink-0 w-12 h-12 bg-white shadow-sm border border-gray-50 rounded-xl flex items-center justify-center text-2xl">
                        {typeInfo?.icon ?? '📌'}
                      </div>
                      <div className="flex-1 min-w-0 pr-8"> 
                        <h3 className="font-bold text-gray-900 text-lg break-words whitespace-normal leading-tight">
                          {rev.entity_name || `${rev.entity_type} #${rev.entity_id}`}
                        </h3>
                        <span className="text-[10px] font-bold text-gray-400 uppercase tracking-widest block mt-1">
                          {new Date(rev.created_at).toLocaleDateString()}
                        </span>
                      </div>
                    </div>
                    
                    
                    <div className="flex items-center gap-2 bg-blue-50/50 self-start px-3 py-1 rounded-lg">
                      <Stars rating={rev.rating} />
                      <span className="text-xs font-bold text-blue-700">{rev.rating}/5</span>
                    </div>

                    {rev.comment && (
                      <div className="text-sm text-gray-600 italic bg-gray-50/80 p-4 rounded-xl border border-gray-100 leading-relaxed">
                        "{rev.comment}"
                      </div>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}