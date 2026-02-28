import { Outlet, NavLink, Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const navItems = [
  { to: '/dashboard',    label: 'Dashboard',    icon: '🏠' },
  { to: '/trips',        label: 'My Trips',     icon: '🗂️' },
  { label: '— Resources —', header: true },
  { to: '/countries',    label: 'Countries',    icon: '🌍' },
  { to: '/cities',       label: 'Cities',       icon: '🏙️' },
  { to: '/attractions',  label: 'Attractions',  icon: '🎡' },
  { to: '/hotels',       label: 'Hotels',       icon: '🏨' },
  { to: '/restaurants',  label: 'Restaurants',  icon: '🍽️' },
  { to: '/flights',      label: 'Flights',      icon: '🛫' },
  { label: '— Account —', header: true },
  { to: '/preferences',  label: 'Preferences',  icon: '⚙️' },
  { to: '/reviews',      label: 'Write Review', icon: '⭐' },
]

export default function Layout() {
  const { logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="flex h-screen overflow-hidden">
      {/* Sidebar */}
      <aside className="w-60 flex-shrink-0 flex flex-col bg-gray-900 text-white overflow-y-auto">
        {/* Logo */}
        <div className="px-5 py-5 border-b border-gray-700">
          <div className="flex items-center gap-2">
            <span className="text-2xl">🗺️</span>
            <div>
              <p className="font-bold text-sm leading-tight">Travel Planning</p>
              <p className="text-xs text-gray-400">System</p>
            </div>
          </div>
        </div>

        {/* Nav Links */}
        <nav className="flex-1 px-3 py-4 space-y-0.5">
          {navItems.map((item, i) =>
            item.header ? (
              <p key={i} className="px-3 pt-4 pb-1 text-xs font-semibold uppercase tracking-wider text-gray-500">
                {item.label}
              </p>
            ) : (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.to === '/trips'}
                className={({ isActive }) =>
                  `flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
                    isActive
                      ? 'bg-brand-600 text-white font-medium'
                      : 'text-gray-300 hover:bg-gray-800 hover:text-white'
                  }`
                }
              >
                <span className="text-base">{item.icon}</span>
                {item.label}
              </NavLink>
            ),
          )}
        </nav>

        {/* Logout */}
        <div className="p-3 border-t border-gray-700">
          <button
            onClick={handleLogout}
            className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-300 hover:bg-red-800 hover:text-white transition-colors"
          >
            <span>🚪</span> Logout
          </button>
        </div>
      </aside>

      {/* Right panel */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top bar */}
        <header className="flex-shrink-0 flex items-center justify-end gap-3 px-6 py-3 bg-white/60 backdrop-blur-sm border-b border-white/80 shadow-sm">
          <Link
            to="/trips/create"
            className="inline-flex items-center gap-2 px-5 py-2 bg-brand-600 text-white text-sm font-semibold rounded-lg hover:bg-brand-700 transition-colors shadow-sm"
          >
            ✈️ Plan a Trip
          </Link>
        </header>

        {/* Main content */}
        <main className="flex-1 overflow-y-auto">
          <div className="max-w-7xl mx-auto px-6 py-8">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  )
}
