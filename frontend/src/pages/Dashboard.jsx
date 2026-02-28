import { Link } from 'react-router-dom'

const cards = [
  { icon: '🌍', title: 'Countries',      desc: 'Browse available countries',     to: '/countries',   color: 'bg-emerald-50/80 border-emerald-200 hover:border-emerald-300 hover:bg-emerald-50' },
  { icon: '🏙️', title: 'Cities',         desc: 'Explore destination cities',     to: '/cities',      color: 'bg-sky-50/80 border-sky-200 hover:border-sky-300 hover:bg-sky-50' },
  { icon: '🏨', title: 'Hotels',         desc: 'Browse accommodations',          to: '/hotels',      color: 'bg-amber-50/80 border-amber-200 hover:border-amber-300 hover:bg-amber-50' },
  { icon: '🎡', title: 'Attractions',    desc: 'Discover local attractions',     to: '/attractions', color: 'bg-purple-50/80 border-purple-200 hover:border-purple-300 hover:bg-purple-50' },
  { icon: '🍽️', title: 'Restaurants',   desc: 'Find dining options',            to: '/restaurants', color: 'bg-orange-50/80 border-orange-200 hover:border-orange-300 hover:bg-orange-50' },
  { icon: '🛫', title: 'Flights',        desc: 'Browse available flights',       to: '/flights',     color: 'bg-blue-50/80 border-blue-200 hover:border-blue-300 hover:bg-blue-50' },
]

export default function Dashboard() {
  return (
    <div>
      {/* ── Hero ── */}
      <div className="relative rounded-2xl overflow-hidden mb-8 shadow-lg">
        {/* Gradient background */}
        <div className="absolute inset-0 bg-gradient-to-br from-blue-900 via-blue-700 to-teal-500" />

        {/* Decorative blobs */}
        <div className="absolute -top-16 -right-16 w-72 h-72 rounded-full bg-white/5" />
        <div className="absolute top-8 right-32 w-40 h-40 rounded-full bg-teal-400/10" />
        <div className="absolute -bottom-12 -left-12 w-56 h-56 rounded-full bg-blue-400/10" />
        <div className="absolute bottom-0 left-1/2 w-96 h-96 rounded-full bg-white/5 -translate-x-1/2 translate-y-1/2" />

        {/* Content */}
        <div className="relative px-10 py-14 text-white">
          <div className="flex items-start gap-4 mb-6">
            <div>
              <p className="text-blue-200 text-sm font-medium tracking-widest uppercase mb-1">
              </p>
              <h1 className="text-4xl font-bold leading-tight">
                Travel Planning
              </h1>
              <p className="text-blue-100 text-lg mt-2 max-w-md">
                Plan your perfect trip — discover flights, hotels, and attractions all in one place.
              </p>
            </div>

            {/* Decorative globe */}
            <div className="ml-auto text-7xl opacity-20 select-none hidden sm:block">
            </div>
          </div>

          <div className="flex flex-wrap gap-3">
          </div>
        </div>
      </div>

      {/* ── Resource grid ── */}
      <p className="text-xs font-semibold uppercase tracking-wider text-gray-400 mb-3">
        Explore
      </p>
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
        {cards.map((c) => (
          <Link
            key={c.to}
            to={c.to}
            className={`block rounded-xl border-2 p-5 transition-all hover:shadow-md ${c.color}`}
          >
            <div className="text-3xl mb-3">{c.icon}</div>
            <h3 className="font-semibold text-gray-900 text-sm">{c.title}</h3>
            <p className="text-xs text-gray-500 mt-1">{c.desc}</p>
          </Link>
        ))}
      </div>
    </div>
  )
}
