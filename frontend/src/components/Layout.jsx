import { Outlet, NavLink, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useTranslation } from 'react-i18next';
import { useNotifications } from '../context/NotificationContext';
import React, { useState, useRef, useEffect } from 'react';

export default function Layout() {
  const { user, logout } = useAuth()
  const { notifications, unreadCount, markAsRead, deleteNotification } = useNotifications()  
  const navigate = useNavigate()
  const location = useLocation()
  const { t, i18n } = useTranslation()
  const [showNotifications, setShowNotifications] = useState(false); 
  const notificationRef = useRef(null);
  const isDashboard = location.pathname === '/dashboard';
  const [isNavOpen, setIsNavOpen] = useState(isDashboard); 
  const [isScrolled, setIsScrolled] = useState(false);

  useEffect(() => {
    function handleClickOutside(event) {
      if (showNotifications && 
          notificationRef.current && 
          !notificationRef.current.contains(event.target)) {
        setShowNotifications(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [showNotifications]);

  useEffect(() => {
    const handleScroll = () => setIsScrolled(window.scrollY > 20);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  useEffect(() => {
    setIsNavOpen(isDashboard);
  }, [location.pathname, isDashboard]);

  const handleLogout = () => {
    e.stopPropagation();
    logout();
    navigate('/login');
  }

  const resourceItems = [
    { to: '/countries',   label: t('nav.countries'),    icon: '🌍' },
    { to: '/cities',      label: t('nav.cities'),       icon: '🏙️' },
    { to: '/attractions', label: t('nav.attractions'),  icon: '🎡' },
    { to: '/hotels',      label: t('nav.hotels'),       icon: '🏨' },
    { to: '/restaurants', label: t('nav.restaurants'),  icon: '🍽️' },
    { to: '/flights',     label: t('nav.flights'),      icon: '🛫' },
  ]

  const accountItems = [
    { to: '/reviews', label: t('nav.write_review'), icon: '⭐' }
  ];
  
  const Dropdown = ({ label, icon, items, isAccount = false }) => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    return (
      <div className="relative group">
        <button 
          type="button"
          className="flex items-center gap-1.5 px-3 py-2 text-sm font-black text-white/90 hover:text-white transition-all capitalize tracking-normal focus:outline-none"
        >
          <span>{icon}</span> {label} 
          <span className="text-[10px] opacity-40 ml-1 transition-transform group-hover:rotate-180">▼</span>
        </button>
        
        <div className={`absolute ${isAccount ? 'right-0' : 'left-0'} mt-0 min-w-[240px] pt-2 invisible opacity-0 group-hover:visible group-hover:opacity-100 transition-all duration-300 z-[160] transform translate-y-2 group-hover:translate-y-0 pointer-events-none group-hover:pointer-events-auto`}>
          <div className="bg-slate-900 border border-white/10 shadow-[0_20px_50px_rgba(0,0,0,0.6)] rounded-[28px] overflow-hidden p-2 space-y-1 text-left">
            
            {isAccount && user && (
              <div className="px-5 py-4 mb-2 border-b border-white/5 bg-white/5 rounded-t-xl text-left pointer-events-none">
                <p className="text-sm font-black text-white capitalize">{user.first_name} {user.last_name}</p>
                <p className="text-[11px] font-bold text-white/50 lowercase">{user.email}</p>
              </div>
            )}

            {items.map((item) => (
              <button
                key={item.to}
                type="button"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation(); 
                  navigate(item.to);
                }}
                className="flex items-center gap-3 w-full px-4 py-3 rounded-xl text-sm font-bold text-white/70 hover:bg-blue-600 hover:text-white transition-all text-left relative z-[170] cursor-pointer"
              >
                <span className="text-lg">{item.icon}</span> {item.label}
              </button>
            ))}

            {isAccount && (
              <button 
                type="button"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  logout();
                  navigate('/login');
                }} 
                className="flex items-center gap-3 w-full px-4 py-3 mt-1 rounded-xl text-sm font-bold text-red-400 hover:bg-red-500/20 transition-all border-t border-white/5 text-left relative z-[170] cursor-pointer"
              >
                <span className="text-lg">🚪</span> {t('nav.logout')}
              </button>
            )}
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="flex flex-col min-h-screen bg-[#f8f9fa] font-sans">
      
      {!isDashboard && !isNavOpen && (
      <div className="fixed top-6 left-6 z-[120]"> 
        <button 
          onClick={() => setIsNavOpen(true)}
          className="bg-slate-950/90 backdrop-blur-md text-white p-3 rounded-xl shadow-2xl border border-white/10 flex flex-col gap-1.5 hover:bg-blue-600 transition-all group focus:outline-none"
          aria-label="Մենյու"
        >
          <span className="w-6 h-0.5 bg-white rounded-full transition-all group-hover:bg-white"></span>
          <span className="w-6 h-0.5 bg-white rounded-full transition-all group-hover:bg-white"></span>
          <span className="w-6 h-0.5 bg-white rounded-full transition-all group-hover:bg-white"></span>
        </button>
      </div>
    )}

      <header className={`fixed top-0 w-full z-[110] transition-all duration-500 
        ${isNavOpen ? 'translate-y-0' : '-translate-y-full'} 
        ${isDashboard && !isScrolled ? 'bg-transparent' : 'bg-slate-950/95 backdrop-blur-xl shadow-2xl'}`}>
        
        <div className="max-w-[1400px] mx-auto px-6 h-20 flex items-center justify-between">
          <div 
            className="flex items-center gap-3 cursor-pointer shrink-0 z-[120]" 
            onClick={(e) => {
              e.stopPropagation(); 
              navigate('/dashboard');
            }}
          >
            <span className="text-2xl">🗺️</span>
          </div>

            <nav className="hidden lg:flex items-center bg-white/5 backdrop-blur-md border border-white/10 rounded-full px-2 py-1">
              <NavLink to="/dashboard" className={({ isActive }) => `px-5 py-2 rounded-full text-xs font-black capitalize transition-all ${isActive ? 'bg-blue-600 text-white' : 'text-white/80 hover:text-white hover:bg-white/5'}`}>
                {t('nav.dashboard')}
              </NavLink>
              <NavLink to="/trips" className={({ isActive }) => `px-5 py-2 rounded-full text-xs font-black capitalize transition-all ${isActive ? 'bg-blue-600 text-white' : 'text-white/80 hover:text-white hover:bg-white/5'}`}>
                {t('nav.my_trips')}
              </NavLink>
              <NavLink to="/preferences" className={({ isActive }) => `px-5 py-2 rounded-full text-xs font-black capitalize transition-all ${isActive ? 'bg-blue-600 text-white' : 'text-white/80 hover:text-white hover:bg-white/5'}`}>
                {t('nav.preferences')}
              </NavLink>
              <Dropdown label={t('nav.resources')} items={resourceItems} />

              {!isDashboard && (
                <button onClick={() => setIsNavOpen(false)} className="ml-2 w-8 h-8 flex items-center justify-center rounded-full hover:bg-white/10 text-white/50 hover:text-white transition-all">
                  ✕
                </button>
              )}
            </nav>

          <div className="flex items-center gap-4">
            <div className="relative" ref={notificationRef}>
              <button 
                onClick={() => setShowNotifications(!showNotifications)}
                className="p-2.5 rounded-full bg-white/5 hover:bg-white/10 text-white transition-all relative">
                <span className="text-xl">🔔</span>
                {unreadCount > 0 && (
                  <span className="absolute top-0 right-0 bg-red-500 text-white text-[10px] h-4 w-4 flex items-center justify-center rounded-full border-2 border-slate-900 animate-pulse">
                    {unreadCount}
                  </span>
                )}
              </button>

              {showNotifications && (
                <div className="absolute right-0 mt-4 w-[450px] max-w-[90vw] bg-white border border-slate-100 shadow-2xl rounded-[28px] overflow-hidden z-20 animate-fade-in">
                  <div className="p-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
                    <h3 className="font-bold text-xs text-slate-500">{t('nav.notifications')}</h3>
                    <span className="text-[10px] bg-blue-600 text-white px-2 py-0.5 rounded-full font-bold">{unreadCount} {t('nav.new')}</span>
                  </div>
                  <div className="max-h-[400px] overflow-y-auto custom-scrollbar bg-white">
                    {notifications.length === 0 ? (
                      <div className="p-12 text-center text-slate-400 text-sm">{t('nav.empty')}</div>
                    ) : (
                      notifications.map((n) => (
                        <div key={n.id} className={`group relative border-b border-slate-50 hover:bg-slate-50 transition-all ${!n.is_read ? 'bg-blue-50/40' : 'bg-white'}`}>
                          <div onClick={async (e) => {
                                e.stopPropagation();
                                await markAsRead(n.id); 
                                setShowNotifications(false);
                                if(n.trip_id) navigate(`/trips/${n.trip_id}/options`); 
                              }}
                              className="p-4 cursor-pointer"
                          >
                            <div className="text-[10px] font-black text-blue-500 mb-1">📍 {n.trip_title || t('nav.trip_update')}</div>
                            <p className={`text-sm leading-relaxed ${!n.is_read ? 'font-bold text-slate-900' : 'text-slate-500'}`}>
                              {t(`nav.${n.type?.toLowerCase()}`, n.message)}
                            </p>
                            <div className="flex justify-between items-center mt-2">
                              <p className="text-[9px] text-slate-400 font-bold">{new Date(n.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</p>
                            </div>
                          </div>
                          <button onClick={(e) => { e.stopPropagation(); deleteNotification(n.id); }}
                                  className="absolute right-3 top-4 p-2 opacity-0 group-hover:opacity-100 hover:bg-red-50 rounded-full text-slate-300 hover:text-red-500 transition-all">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M6 18L18 6M6 6l12 12" /></svg>
                          </button>
                        </div>
                      ))
                    )}
                  </div>
                </div>
              )}
            </div>

            <Dropdown 
              label={`${user?.first_name || user?.FirstName || t('nav.account')}`} 
              items={accountItems} 
              isAccount={true} 
            />
          </div>
        </div>
      </header>

      <main className={`flex-grow transition-all duration-500 
      ${isDashboard ? 'pt-0 bg-transparent' : 'pt-24 bg-[#f4f7fa]'}`}>
        <div className={`${isDashboard ? 'w-full' : 'max-w-10xl mx-auto px-6'}`}>
          <Outlet />
        </div>
      </main>

      <footer className="bg-[#0a0a1a] text-white py-16 px-6 mt-auto border-t border-white/5">
        <div className="max-w-7xl mx-auto flex flex-col items-center text-center">
          <div className="flex items-center gap-2 mb-4">
            <span className="text-3xl">🗺️</span>
            <p className="font-black text-2xl capitalize text-white">{t('dashboard.hero.title')}</p>
          </div>
          <p className="text-gray-500 text-[10px] font-bold uppercase tracking-widest mb-8">© 2026 TravelPlan. {t('common.all_rights_reserved')}</p>
          <div className="relative group">
            <div className="absolute inset-y-0 left-4 flex items-center pointer-events-none text-xs">🌐</div>
            <select 
              value={i18n.language} 
              onChange={(e) => i18n.changeLanguage(e.target.value)}
              className="bg-[#16162a] text-white border border-white/10 rounded-2xl pl-10 pr-12 py-3 text-sm font-bold outline-none cursor-pointer hover:border-blue-500 transition-all appearance-none min-w-[220px]"
            >
              <option value="hy">Հայերեն (ARM)</option>
              <option value="en">English (ENG)</option>
            </select>
            <div className="absolute inset-y-0 right-4 flex items-center pointer-events-none opacity-40 text-[8px]">▼</div>
          </div>
        </div>
      </footer>
    </div>
  )
}