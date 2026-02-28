import { useState, useEffect, useCallback } from 'react'
import { getUserReviews, createReview } from '../api/reviews'
import { getHotels, getAttractions, getRestaurants } from '../api/resources'
import PageHeader from '../components/PageHeader'

const ENTITY_TYPES = [
  { value: 'hotel',      label: 'Hotel',      icon: '🏨' },
  { value: 'attraction', label: 'Attraction',  icon: '🎡' },
  { value: 'restaurant', label: 'Restaurant',  icon: '🍽️' },
]

const loaders = {
  hotel:      getHotels,
  attraction: getAttractions,
  restaurant: getRestaurants,
}

function getEntityId(type, entity) {
  if (type === 'hotel')      return entity.hotel_id
  if (type === 'attraction') return entity.attraction_id
  return entity.restaurant_id
}

function getEntityLabel(type, entity) {
  if (type === 'hotel') return `${entity.name}${entity.stars ? ` (${'⭐'.repeat(entity.stars)})` : ''}`
  return entity.name
}

function Stars({ rating }) {
  return (
    <span>
      {Array.from({ length: 5 }, (_, i) => (
        <span key={i} className={i < rating ? 'text-yellow-400' : 'text-gray-200'}>★</span>
      ))}
    </span>
  )
}

export default function Reviews() {
  const [entityType, setEntityType] = useState('hotel')
  const [entities, setEntities] = useState([])
  const [form, setForm] = useState({ entity_id: '', rating: 5, comment: '' })
  const [reviews, setReviews] = useState([])
  const [loadingReviews, setLoadingReviews] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

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
    loaders[entityType]().then(setEntities).catch(() => setEntities([]))
  }, [entityType])

  const set = (field) => (e) => setForm({ ...form, [field]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setSubmitting(true)
    try {
      await createReview({
        entity_type: entityType,
        entity_id: parseInt(form.entity_id),
        rating: parseInt(form.rating),
        comment: form.comment,
      })
      setForm({ entity_id: '', rating: 5, comment: '' })
      fetchReviews()
    } catch (err) {
      setError(err.response?.data || 'Failed to submit review.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="space-y-8">
      {/* ── Write a review ── */}
      <div>
        <PageHeader icon="⭐" title="Write a Review" subtitle="Share your experience" />

        <div className="card max-w-lg">
          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{error}</div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">
            {/* Category toggle */}
            <div>
              <label className="label">Category</label>
              <div className="flex gap-2">
                {ENTITY_TYPES.map((t) => (
                  <button
                    key={t.value}
                    type="button"
                    onClick={() => setEntityType(t.value)}
                    className={`flex-1 flex items-center justify-center gap-1.5 py-2 rounded-lg border text-sm font-medium transition-colors ${
                      entityType === t.value
                        ? 'border-brand-500 bg-brand-50 text-brand-700'
                        : 'border-gray-200 hover:border-gray-300 text-gray-600'
                    }`}
                  >
                    {t.icon} {t.label}
                  </button>
                ))}
              </div>
            </div>

            {/* Entity select */}
            <div>
              <label className="label">
                Select {ENTITY_TYPES.find((t) => t.value === entityType)?.label}
              </label>
              <select className="input" value={form.entity_id} onChange={set('entity_id')} required>
                <option value="">Choose…</option>
                {entities.map((e) => (
                  <option key={getEntityId(entityType, e)} value={getEntityId(entityType, e)}>
                    {getEntityLabel(entityType, e)}
                  </option>
                ))}
              </select>
            </div>

            {/* Star rating */}
            <div>
              <label className="label">Rating</label>
              <div className="flex gap-1 mt-1">
                {[1, 2, 3, 4, 5].map((n) => (
                  <button
                    key={n}
                    type="button"
                    onClick={() => setForm({ ...form, rating: n })}
                    className={`text-3xl transition-all hover:scale-110 ${form.rating >= n ? 'opacity-100' : 'opacity-25'}`}
                  >
                    ⭐
                  </button>
                ))}
                <span className="ml-2 self-center text-sm text-gray-500">{form.rating}/5</span>
              </div>
            </div>

            {/* Comment */}
            <div>
              <label className="label">Comment</label>
              <textarea
                className="input"
                rows={3}
                value={form.comment}
                onChange={set('comment')}
                placeholder="Share your experience…"
                required
              />
            </div>

            <button type="submit" disabled={submitting} className="btn-primary w-full justify-center">
              {submitting ? 'Submitting…' : 'Submit Review'}
            </button>
          </form>
        </div>
      </div>

      {/* ── My Reviews ── */}
      <div>
        <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
          <span>📋</span> My Reviews
          {!loadingReviews && (
            <span className="text-sm font-normal text-gray-400">({reviews.length})</span>
          )}
        </h2>

        {loadingReviews && (
          <div className="flex items-center gap-2 text-gray-400">
            <div className="w-4 h-4 border-2 border-gray-300 border-t-brand-600 rounded-full animate-spin" />
            Loading reviews…
          </div>
        )}

        {!loadingReviews && reviews.length === 0 && (
          <div className="card text-center py-10 text-gray-400">
            <p className="text-2xl mb-2">📝</p>
            No reviews submitted yet.
          </div>
        )}

        {!loadingReviews && reviews.length > 0 && (
          <div className="space-y-3 max-w-2xl">
            {reviews.map((rev) => {
              const typeInfo = ENTITY_TYPES.find((t) => t.value === rev.entity_type)
              return (
                <div key={rev.review_id} className="card py-4">
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex items-center gap-2 flex-shrink-0">
                      <span className="text-xl">{typeInfo?.icon ?? '📌'}</span>
                      <span className="text-xs font-medium text-gray-500 capitalize bg-gray-100 px-2 py-0.5 rounded-full">
                        {rev.entity_type} #{rev.entity_id}
                      </span>
                    </div>
                    <div className="text-right flex-shrink-0">
                      <Stars rating={rev.rating} />
                      <p className="text-xs text-gray-400 mt-0.5">
                        {new Date(rev.review_date).toLocaleDateString('en-US', {
                          year: 'numeric', month: 'short', day: 'numeric',
                        })}
                      </p>
                    </div>
                  </div>
                  {rev.comment && (
                    <p className="mt-3 text-sm text-gray-700 leading-relaxed">{rev.comment}</p>
                  )}
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}
