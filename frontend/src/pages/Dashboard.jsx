import { Link, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useState, useRef } from 'react';
import { useEffect } from "react";
import heroBg from '../assets/background.jpg';
import filghtsBg from '../assets/flight.jpg';
import countriesBg from '../assets/countries.jpg';
import citiesBg from '../assets/city.jpg';
import hotelsBg from '../assets/hotel.jpg';
import attractionsBg from '../assets/attraction.jpg';
import restaurantsBg from '../assets/restaurant.jpg';

export default function Dashboard() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const prepSectionRef = useRef(null);
  const [showPrep, setShowPrep] = useState(false); 

  const cards = [
    { id: 1, icon: '🌍', title: t('nav.countries'),   to: '/countries',   img: countriesBg, rotate: '-rotate-2' },
    { id: 2, icon: '🏙️', title: t('nav.cities'),      to: '/cities',      img: citiesBg, rotate: 'rotate-3' },
    { id: 3, icon: '🏨', title: t('nav.hotels'),      to: '/hotels',      img: hotelsBg, rotate: 'rotate-1' },
    { id: 4, icon: '🎡', title: t('nav.attractions'), to: '/attractions', img: attractionsBg, rotate: '-rotate-3' },
    { id: 5, icon: '🍽️', title: t('nav.restaurants'), to: '/restaurants', img: restaurantsBg, rotate: 'rotate-2' },
    { id: 6, icon: '🛫', title: t('nav.flights'),      to: '/flights',     img: filghtsBg, rotate: '-rotate-1' },
  ];

  const handlePlanClick = () => {
  const willShow = !showPrep;
  setShowPrep(willShow);
  
  if (willShow) {
    setTimeout(() => {
      prepSectionRef.current?.scrollIntoView({ 
        behavior: 'smooth', 
        block: 'start' 
      });
    }, 450); 
  }
};

  return (
    <div className="pb-24 text-left font-sans bg-slate-50 min-h-screen">
      
      <div className="relative w-screen left-1/2 -translate-x-1/2 mb-20 overflow-hidden">
        <div className="relative flex flex-col md:flex-row items-stretch min-h-[600px]">
          
          <div className="relative z-10 bg-slate-950 text-white p-10 md:p-20 md:w-[50%] flex flex-col justify-center" 
               style={{ clipPath: window.innerWidth > 768 ? 'polygon(0 0, 100% 0, 90% 100%, 0% 100%)' : 'none' }}>
            <div className="max-w-xl">
              <h1 className="text-5xl md:text-7xl font-black leading-none mb-8 capitalize italic tracking-tighter">
                {t('dashboard.hero.title')}
              </h1>
              
              <p className="text-lg md:text-xl text-gray-400 font-medium mb-12 opacity-90 leading-relaxed">
                {t('dashboard.hero.description')}
              </p>
              
              <button 
                onClick={handlePlanClick}
                className="bg-blue-600 text-white px-8 py-3 rounded-xl font-bold hover:bg-blue-700 transition-all shadow-md active:scale-95"
              >
              {t('dashboard.hero.plan_button')}
                <span className={`text-xl transition-transform duration-300 ${showPrep ? 'rotate-180' : 'translate-y-1'}`}></span>
              </button>
            </div>
          </div>

          <div className="relative h-[300px] md:h-auto md:w-[55%] -ml-[5%] z-0">
            <div className="absolute inset-0 bg-cover bg-center" style={{ backgroundImage: `url(${heroBg})` }} />
            <div className="absolute inset-0 bg-gradient-to-r from-slate-950 via-transparent to-transparent md:block hidden" />
          </div>
        </div>

        <div 
          ref={prepSectionRef}
          className={`transition-all duration-700 ease-in-out overflow-hidden scroll-mt-40 mx-auto max-w-5xl px-4
            ${showPrep ? 'max-h-[1000px] opacity-100 mt-20 mb-32' : 'max-h-0 opacity-0 mt-0 mb-0'}`}
        >
          <div className="max-w-4xl mx-auto bg-white rounded-[40px] shadow-2xl border-t-[6px] border-blue-600 overflow-hidden">
            <div className="p-10 md:p-16 text-center space-y-8">
              <span className="text-blue-600 font-black text-xs uppercase tracking-widest mb-4 block">
                {t('dashboard.prep.subtitle')}
              </span>
              <h2 className="text-4xl font-black text-slate-900 uppercase italic mb-6">
                {t('dashboard.prep.title')}
              </h2>
              <p className="text-slate-500 font-medium mb-10">
                {t('dashboard.prep.description')}
              </p>

              <div className="flex flex-col md:flex-row gap-6 justify-center items-center">
                <button 
                  onClick={() => navigate('/preferences')}
                  className="bg-blue-600 text-white px-8 py-3 rounded-xl font-bold hover:bg-blue-700 transition-all shadow-md active:scale-95"
                >
                  {t('dashboard.prep.customize_button')}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <section className="max-w-7xl mx-auto px-4 md:px-8 mt-20">
        <div className="flex items-center gap-6 mb-20">
          <div className="h-[2px] w-12 bg-slate-950" />
          <h2 className="text-4xl md:text-6xl font-black capitalize italic tracking-tighter text-slate-900">
            {t('dashboard.resources.title')}
          </h2>
          <div className="h-[1px] flex-1 bg-slate-200" />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-y-16 gap-x-10">
          {cards.map((c) => (
            <Link
              key={c.to}
              to={c.to}
              className={`group relative block transition-all duration-300 hover:-translate-y-4 h-full ${c.rotate}`}
            >
              <div className="absolute inset-0 bg-slate-950 translate-x-2 translate-y-2 rounded-[2.5rem] z-0 transition-transform group-hover:translate-x-1 group-hover:translate-y-1" />
              <div className="relative z-10 bg-white p-6 rounded-[2.5rem] border-4 border-slate-950 shadow-sm h-full flex flex-col overflow-hidden">
                <div className="relative h-48 rounded-3xl overflow-hidden mb-6 border-2 border-slate-950">
                  <img src={c.img} alt="" className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110" />
                  <div className="absolute top-4 left-4 w-10 h-10 bg-white rounded-xl flex items-center justify-center text-xl shadow-lg border-2 border-slate-950">
                    {c.icon}
                  </div>
                </div>
                <h3 className="font-black text-slate-950 text-2xl xl:text-3xl uppercase italic tracking-tighter mb-6">{c.title}</h3>
                <div className="inline-flex items-center gap-2 self-start text-slate-950 bg-white px-5 py-2 rounded-full font-black text-[10px] uppercase border-2 border-slate-950 transition-all group-hover:bg-slate-950 group-hover:text-white">
                  {t('common.explore')} <span>→</span>
                </div>
              </div>
            </Link>
          ))}
        </div>
      </section>
    </div>
  );
}