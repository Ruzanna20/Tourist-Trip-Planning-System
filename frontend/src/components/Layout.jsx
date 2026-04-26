import { Outlet, NavLink, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useTranslation } from 'react-i18next';
import { useNotifications } from '../context/NotificationContext';
import { useState, useRef, useEffect } from 'react';

export default function Layout() {
  const { logout } = useAuth()
  const { notifications, unreadCount, markAsRead, deleteNotification } = useNotifications()  
  const navigate = useNavigate()
  const location = useLocation()
  const { t, i18n } = useTranslation()
  const [showNotifications, setShowNotifications] = useState(false); 
  const notificationRef = useRef(null);

  useEffect(() => {
    function handleClickOutside(event) {
      if (showNotifications && 
          notificationRef.current && 
          !notificationRef.current.contains(event.target)) {
        setShowNotifications(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [showNotifications]);

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const resourceItems = [
    { to: '/countries',   label: t('nav.countries'),     icon: '🌍' },
    { to: '/cities',      label: t('nav.cities'),        icon: '🏙️' },
    { to: '/attractions', label: t('nav.attractions'),   icon: '🎡' },
    { to: '/hotels',      label: t('nav.hotels'),        icon: '🏨' },
    { to: '/restaurants', label: t('nav.restaurants'),   icon: '🍽️' },
    { to: '/flights',     label: t('nav.flights'),       icon: '🛫' },
  ]

  const accountItems = [
    { to: '/preferences', label: t('nav.preferences'),  icon: '⚙️' },
    { to: '/reviews',     label: t('nav.write_review'), icon: '⭐' },
  ]

  const Dropdown = ({ label, icon, items }) => (
    <div className="relative group flex-shrink-0">
      <button className="flex items-center gap-1.5 px-3 py-2 text-sm font-black text-slate-700 hover:text-blue-600 transition-colors uppercase tracking-tight">
        <span>{icon}</span> {label} <span className="text-[8px] opacity-40 ml-0.5">▼</span>
      </button>
      
      <div className="absolute left-0 mt-0 min-w-max pt-2 invisible opacity-0 group-hover:visible group-hover:opacity-100 transition-all duration-200 z-[100]">
        <div className="bg-white border border-slate-100 shadow-2xl rounded-[24px] overflow-hidden p-2 space-y-1 text-left">
          {items.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) =>
                `flex items-center gap-3 px-4 py-3 pr-8 rounded-xl text-sm font-bold transition-all whitespace-nowrap ${
                  isActive ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50'
                }`
              }
            >
              <span className="text-lg">{item.icon}</span> {item.label}
            </NavLink>
          ))}
        </div>
      </div>
    </div>
  )

  return (
    <div className="flex flex-col min-h-screen bg-[#ececii] font-sans">
      <header className="sticky top-0 z-[110] bg-white/80 backdrop-blur-md border-b border-slate-100 px-6 py-3 flex items-center justify-between shadow-sm">
        
        <div className="flex items-center gap-2 flex-shrink-0">
          <span className="text-2xl" onClick={() => navigate('/dashboard')} style={{cursor:'pointer'}}>🗺️</span>
          <p className="font-black text-lg uppercase tracking-tighter text-slate-900 italic">
            {t('dashboard.hero.title')}
          </p>
        </div>

        <nav className="hidden lg:flex items-center gap-1 xl:gap-2 flex-grow justify-center mx-4">
          <NavLink to="/dashboard" className={({ isActive }) => `flex items-center gap-2 px-5 py-2.5 rounded-full text-xs xl:text-sm font-black uppercase tracking-tight transition-all ${isActive ? 'bg-blue-600 text-white shadow-lg shadow-blue-200' : 'text-slate-600 hover:bg-slate-100'}`}>
            🏠 {t('nav.dashboard')}
          </NavLink>
          <NavLink to="/trips" className={({ isActive }) => `flex items-center gap-2 px-5 py-2.5 rounded-full text-xs xl:text-sm font-black uppercase tracking-tight transition-all ${isActive ? 'bg-blue-600 text-white shadow-lg shadow-blue-200' : 'text-slate-600 hover:bg-slate-100'}`}>
            🗂️ {t('nav.my_trips')}
          </NavLink>
          <Dropdown label={t('nav.resources')} icon="📚" items={resourceItems} />
          <Dropdown label={t('nav.account')} icon="👤" items={accountItems} />
        </nav>

        <div className="flex items-center gap-3 flex-shrink-0">
          <div className="relative" ref={notificationRef}>
            <button 
              onClick={() => setShowNotifications(!showNotifications)}
              className="relative p-2 rounded-full hover:bg-slate-50 transition-all"
            >
              <span className="text-xl">🔔</span>
              {unreadCount > 0 && (
                <span className="absolute top-1.5 right-1.5 bg-red-500 text-white text-[10px] h-4 w-4 flex items-center justify-center rounded-full border-2 border-white animate-pulse">
                  {unreadCount}
                </span>
              )}
            </button>

            {showNotifications && (
              <div className="absolute right-0 mt-3 w-80 bg-white border border-slate-100 shadow-2xl rounded-[24px] overflow-hidden z-20 transition-all">
                <div className="p-4 border-b border-slate-50 flex justify-between items-center">
                  <h3 className="font-black text-xs uppercase tracking-widest text-slate-400">{t('nav.notifications')}</h3>
                  <span className="text-[10px] bg-blue-50 text-blue-600 px-2 py-0.5 rounded-full font-bold">
                    {unreadCount} {t('nav.new')}
                  </span>
                </div>
                
                <div className="max-h-[400px] overflow-y-auto custom-scrollbar">
                  {notifications.length === 0 ? (
                    <div className="p-8 text-center text-slate-400 text-sm">
                       {t('nav.empty')}
                    </div>
                  ) : (
                    notifications.map((n) => (
                      <div 
                        key={n.id} 
                        className="group relative border-b border-slate-50 hover:bg-slate-50 transition-all"
                      >
                        <div 
                          onClick={async (e) => {
                            e.stopPropagation();
                            await markAsRead(n.id); 
                            setShowNotifications(false);
                            if(n.trip_id) navigate(`/trips/${n.trip_id}/options`); 
                          }}
                          className={`p-4 cursor-pointer ${!n.is_read ? 'bg-blue-50/40' : 'bg-white'}`}
                        >
                          <div className="text-[10px] font-black text-blue-500 uppercase mb-1">
                            📍 {n.trip_title || t('nav.trip_update')}
                          </div>
                          <p className={`text-sm leading-tight ${!n.is_read ? 'font-bold text-slate-900' : 'text-slate-500'}`}>
                            {t(`nav.${n.type?.toLowerCase()}`, n.message)}
                          </p>
                          <p className="text-[9px] text-slate-400 font-bold uppercase mt-2">
                            {new Date(n.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                          </p>
                        </div>

                        <button
                          onClick={(e) => {
                            e.stopPropagation(); 
                            deleteNotification(n.id);
                          }}
                          className="absolute right-3 top-1/2 -translate-y-1/2 p-2 opacity-0 group-hover:opacity-100 hover:bg-red-50 rounded-full text-slate-300 hover:text-red-500 transition-all"
                          title={t('nav.delete')}
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M6 18L18 6M6 6l12 12" />
                          </svg>
                        </button>
                      </div>
                    ))
                  )}
                </div>
              </div>
            )}
          </div>

          <button onClick={handleLogout} className="flex items-center gap-2 px-3 py-2 text-sm font-black text-slate-700 hover:text-red-600 transition-colors uppercase tracking-tight group">
            <span className="text-xl group-hover:rotate-12 transition-transform">🚪</span>
            <span className="hidden sm:inline text-[10px] font-black uppercase tracking-widest">{t('nav.logout')}</span>
          </button>
        </div>
      </header>

      <main className="flex-grow p-4 md:p-8">
        <div className="max-w-7xl mx-auto">
          <Outlet />
        </div>
      </main>

      <footer className="bg-[#0a0a1a] text-white py-12 px-6 mt-auto">
        <div className="max-w-7xl mx-auto flex flex-col items-center">
          <div className="flex items-center gap-2 mb-4">
            <span className="text-3xl">🗺️</span>
            <p className="font-black text-2xl uppercase tracking-tighter italic text-white">
               {t('dashboard.hero.title')}
            </p>
          </div>
          
          <p className="text-gray-500 text-[10px] font-bold uppercase tracking-widest mb-8">
            © 2026 TravelPlan. {t('common.all_rights_reserved') || "All rights reserved."}
          </p>

          <div className="relative mb-10 group">
            <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none">
              <span className="text-xs">🌐</span>
            </div>
            <select 
              value={i18n.language} 
              onChange={(e) => i18n.changeLanguage(e.target.value)}
              className="bg-[#16162a] text-white border border-gray-800 rounded-2xl pl-10 pr-10 py-3 text-sm font-bold outline-none cursor-pointer hover:border-gray-600 transition-all appearance-none min-w-[200px] text-center"
            >
              <option value="hy">Հայերեն (ARM)</option>
              <option value="en">English (ENG)</option>
            </select>
            <div className="absolute inset-y-0 right-4 flex items-center pointer-events-none opacity-40 text-[8px]">
              ▼
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}