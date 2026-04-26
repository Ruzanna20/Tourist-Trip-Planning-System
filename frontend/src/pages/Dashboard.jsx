import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next';
import heroBg from '../assets/background.jpg';
import filghtsBg from '../assets/flight.jpg';
import countriesBg from '../assets/countries.jpg';
import citiesBg from '../assets/city.jpg';
import hotelsBg from '../assets/hotel.jpg';
import attractionsBg from '../assets/attraction.jpg';
import restaurantsBg from '../assets/restaurant.jpg';

export default function Dashboard() {
  const { t, i18n } = useTranslation();

  const cards = [
    { id: 1, icon: '🌍', title: t('nav.countries'),   to: '/countries',   img: countriesBg, rotate: '-rotate-2' },
    { id: 2, icon: '🏙️', title: t('nav.cities'),      to: '/cities',      img: citiesBg, rotate: 'rotate-3' },
    { id: 3, icon: '🏨', title: t('nav.hotels'),      to: '/hotels',      img: hotelsBg, rotate: 'rotate-1' },
    { id: 4, icon: '🎡', title: t('nav.attractions'), to: '/attractions', img: attractionsBg, rotate: '-rotate-3' },
    { id: 5, icon: '🍽️', title: t('nav.restaurants'), to: '/restaurants', img: restaurantsBg, rotate: 'rotate-2' },
    { id: 6, icon: '🛫', title: t('nav.flights'),      to: '/flights',     img: filghtsBg, rotate: '-rotate-1' },
  ];

  return (
    <div className="pb-24 text-left font-sans bg-[#ececii]">
      
      <div className="relative flex flex-col md:flex-row items-stretch min-h-[500px] mb-24 lg:mb-32 overflow-hidden md:overflow-visible">
        
        <div className="relative z-10 bg-slate-950 text-white p-10 md:p-16 md:w-[55%] flex flex-col justify-center" 
             style={{ clipPath: window.innerWidth > 768 ? 'polygon(0 0, 100% 0, 85% 100%, 0% 100%)' : 'none' }}>
          <div className="max-w-md">
            <h1 className="text-4xl md:text-6xl font-black text-white leading-[1.1] mb-6 uppercase italic tracking-tighter">
              {t('dashboard.hero.title')}
            </h1>
            
            <p className="text-base md:text-lg text-gray-300 font-medium leading-relaxed mb-10 opacity-90">
              {t('dashboard.hero.description')}
            </p>
            
            <Link 
              to="/trips/create" 
              className="inline-flex items-center gap-4 bg-white text-slate-950 px-8 py-4 rounded-full font-black uppercase text-[10px] tracking-widest hover:bg-blue-600 hover:text-white transition-all shadow-2xl active:scale-95 group/btn"
            >
              {t('nav.plan_trip')} 
              <span className="text-xl transition-transform group-hover/btn:translate-x-2">→</span>
            </Link>
          </div>
        </div>

        <div className="relative h-[300px] md:h-auto md:absolute md:inset-y-0 md:right-0 md:w-[60%] z-0 shadow-2xl">
          <div 
            className="absolute inset-0 bg-cover bg-center transition-transform duration-1000" 
            style={{ backgroundImage: `url(${heroBg})` }}
          />
          <div className="absolute inset-0 bg-black/20" />
        </div>
      </div>

      <section className="px-4 md:px-8">
        <div className="flex items-center gap-6 mb-20">
          <div className="h-[2px] w-12 bg-slate-950" />
          <h2 className="text-4xl md:text-6xl font-black uppercase italic tracking-tighter text-slate-900 whitespace-nowrap">
            {t('dashboard.resources.title')}
          </h2>
          <div className="h-[1px] flex-1 bg-slate-200" />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-y-16 gap-x-10 items-stretch">
          {cards.map((c) => (
            <Link
              key={c.to}
              to={c.to}
              className={`group relative block transition-all duration-300 hover:-translate-y-4 h-full ${c.rotate}`}
            >
              <div className="absolute inset-0 bg-slate-950 translate-x-2 translate-y-2 rounded-[2.5rem] z-0 transition-transform group-hover:translate-x-1 group-hover:translate-y-1" />
              
              <div className="relative z-10 bg-white p-6 rounded-[2.5rem] border-4 border-slate-950 shadow-sm overflow-hidden h-full flex flex-col transition-all">
                
                <div className="relative h-44 rounded-3xl overflow-hidden mb-6 border-2 border-slate-950 flex-shrink-0">
                  <img src={c.img} alt="" className="absolute inset-0 w-full h-full object-cover transition-transform duration-700 group-hover:scale-110" />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
                  {/* Icon Badge */}
                  <div className="absolute top-4 left-4 w-10 h-10 bg-white rounded-xl flex items-center justify-center text-xl shadow-lg border-2 border-slate-950">
                    {c.icon}
                  </div>
                </div>

                <div className="flex-grow flex flex-col justify-between">
                  <h3 className="font-black text-slate-950 text-2xl xl:text-3xl uppercase italic tracking-tighter mb-6 leading-tight">
                    {c.title}
                  </h3>
                  
                  <div className="inline-flex items-center gap-2 self-start text-slate-950 bg-white px-5 py-2 rounded-full font-black text-[10px] uppercase tracking-widest shadow-md border-2 border-slate-950 transition-all group-hover:bg-slate-950 group-hover:text-white">
                    {t('common.explore') || 'Explore'} <span>→</span>
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>
      </section>
    </div>
  )
}