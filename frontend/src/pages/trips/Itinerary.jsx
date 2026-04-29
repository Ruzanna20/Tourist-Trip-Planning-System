import { useState, useEffect, useCallback, memo } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom' 
import { useTranslation } from 'react-i18next'
import { getTripItinerary, getItineraryActivities, swapActivity, completeTrip } from '../../api/trips'
import PageHeader from '../../components/PageHeader'

const ACTIVITY_META = {
  hotel:      { icon: '🏨', color: 'bg-orange-50 border-orange-200 text-orange-700', modalHeader: 'bg-orange-500' },
  attraction: { icon: '🎡', color: 'bg-purple-50 border-purple-200 text-purple-700', modalHeader: 'bg-purple-600' },
  restaurant: { icon: '🍽️', color: 'bg-yellow-50 border-yellow-200 text-yellow-700', modalHeader: 'bg-yellow-400' },
  flight:     { icon: '🛫', color: 'bg-blue-50 border-blue-200 text-blue-700', modalHeader: 'bg-blue-600' },
};

const formatDynamicDate = (dateStr, lang) => {
  const d = new Date(dateStr);

  if (lang === 'hy') {
    const days = ['ԿԻՐԱԿԻ', 'ԵՐԿՈՒՇԱԲԹԻ', 'ԵՐԵՔՇԱԲԹԻ', 'ՉՈՐԵՔՇԱԲԹԻ', 'ՀԻՆԳՇԱԲԹԻ', 'ՈՒՐԲԱԹ', 'ՇԱԲԱԹ'];
    const months = ['ԱՊՐԻԼ', 'ՄԱՅԻՍ', 'ՀՈՒՆԻՍ', 'ՀՈՒԼԻՍ', 'ՕԳՈՍՏՈՍ', 'ՍԵՊՏԵՄԲԵՐ', 'ՀՈԿՏԵՄԲԵՐ', 'ՆՈՅԵՄԲԵՐ', 'ԴԵԿՏԵՄԲԵՐ', 'ՀՈՒՆՎԱՐ', 'ՓԵՏՐՎԱՐ', 'ՄԱՐՏ'];
    const armMonthsCorrect = ['ՀՈՒՆՎԱՐ', 'ՓԵՏՐՎԱՐ', 'ՄԱՐՏ', 'ԱՊՐԻԼ', 'ՄԱՅԻՍ', 'ՀՈՒՆԻՍ', 'ՀՈՒԼԻՍ', 'ՕԳՈՍՏՈՍ', 'ՍԵՊՏԵՄԲԵՐ', 'ՀՈԿՏԵՄԲԵՐ', 'ՆՈՅԵՄԲԵՐ', 'ԴԵԿՏԵմբեր'];
    
    return `${days[d.getDay()]}, ${d.getDate()} ${armMonthsCorrect[d.getMonth()]}`;
  }

  return new Intl.DateTimeFormat('en-US', {
    weekday: 'long',
    day: 'numeric',
    month: 'long'
  }).format(d).toUpperCase();
};

function DetailsModal({ item, type, onClose, onSwap }) {
  const { t } = useTranslation();
  const [isSwapping, setIsSwapping] = useState(false);

  const Label = ({ children }) => (
    <label className="text-[10px] font-black uppercase text-gray-400 tracking-widest mb-2 block">{children}</label>
  );

  const searchStr = "alternative is ";
  const index = item.notes?.indexOf(searchStr);
  const hasAlternative = index !== -1 && index !== undefined;
  const alternativeName = hasAlternative ? item.notes.substring(index + searchStr.length).replace(".", "").trim() : "";

  const handleSwap = async () => {
    if (!alternativeName || alternativeName.toLowerCase() === 'place') return;
    setIsSwapping(true);
    try {
      if (onSwap) {
        await onSwap(alternativeName);
        onClose();
      }
    } catch (error) {
      console.error("Swap failed", error);
    } finally {
      setIsSwapping(false);
    }
  };

  const renderTranslatedNotes = (note) => {
    if (!note) return null;
    if (note.includes("arrive at the airport")) return t('trips.itinerary.notes.flight');
    if (note.includes("Accommodation check-in")) {
      return (item.order_number && item.order_number < 2) ? t('trips.itinerary.notes.hotel_welcome') : t('trips.itinerary.notes.hotel_return');
    }
    if (note.includes("Enjoy your meal")) return t('trips.itinerary.notes.restaurant');
    if (note.includes("Welcome. Enjoy a relaxing walk")) return t('trips.itinerary.notes.welcome_walk');
    if (note.includes("Last day. Perfect time")) return t('trips.itinerary.notes.last_day');
    if (note.includes("the best nearby alternative is")) {
      return <p>{t('trips.itinerary.notes.alternative_prefix')} <span className="font-bold text-slate-900">{alternativeName}</span></p>;
    }
    return note;
  };

  const flightInfo = type === 'flight' ? getFlightDetails(item.address, item.description) : null;
  const headerColor = ACTIVITY_META[type]?.modalHeader || 'bg-gray-800';
  const translatedNote = renderTranslatedNotes(item.notes);

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[150] p-4 text-left">
      <div className="bg-white rounded-[32px] max-w-md w-full overflow-hidden shadow-2xl relative transition-all animate-in zoom-in duration-200">
        <div className={`min-h-[120px] ${headerColor} flex items-end p-8 relative`}>
          <button onClick={onClose} className="absolute top-6 right-6 z-50 text-white bg-white/20 hover:bg-white/40 p-2 rounded-full transition-all">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
          <h2 className="text-xl font-black text-white uppercase tracking-tight break-words w-full leading-tight pr-12 drop-shadow-md">
            {ACTIVITY_META[type]?.icon} {type === 'flight' ? `${item.name} Flight` : item.name}
          </h2>
        </div>
        <div className="p-8 space-y-6 max-h-[70vh] overflow-y-auto text-left">
          {type === 'flight' && flightInfo && (
            <div className="space-y-6">
              <div className="grid grid-cols-2 gap-6">
                <div><Label>{t('flights.table.origin')}</Label><p className="font-bold text-gray-800">🛫 {flightInfo.origin}</p></div>
                <div><Label>{t('flights.table.destination')}</Label><p className="font-bold text-gray-800">🛬 {flightInfo.destination}</p></div>
              </div>
              <div className="grid grid-cols-2 gap-6">
                <div><Label>{t('flights.table.duration')}</Label><p className="font-medium text-gray-700">⏱️ {flightInfo.duration}</p></div>
                <div><Label>{t('flights.table.price')}</Label><p className="font-bold text-green-600 text-lg">💰 ${item.price}</p></div>
              </div>
              <div><Label>{t('flights.table.source')}</Label><p className="text-xs text-gray-500 font-mono italic break-all bg-gray-50 p-2 rounded-lg">{flightInfo.source}</p></div>
            </div>
          )}
          {type === 'hotel' && (
            <div className="space-y-6">
              <div><Label>{t('hotels.table.address')}</Label><p className="text-sm font-medium text-gray-700 break-words font-mono">📍 {item.address}</p></div>
              <div className="grid grid-cols-2 gap-6">
                <div><Label>{t('hotels.table.stars')}</Label><p className="text-amber-500 font-bold text-lg">{item.stars && item.stars > 0 ? "★".repeat(Math.min(5, item.stars)) : "---"}</p></div>
                <div><Label>{t('hotels.table.rating')}</Label><p className="font-bold text-blue-600 text-lg">📊 {item.rating} / 10</p></div>
              </div>
              <div><Label>{t('hotels.table.price')}</Label><p className="font-bold text-green-600 text-xl">💰 ${item.price} <span className="text-xs font-normal text-gray-400 italic">/ {t('flights.units.night')}</span></p></div>
              <div><Label>{t('hotels.table.description')}</Label><div className="p-4 bg-gray-50 rounded-2xl text-sm text-gray-600 italic border border-gray-100">{item.description || t('common.no_description')}</div></div>
            </div>
          )}

          {type === 'attraction' && (
            <div className="space-y-6">
              {(!item.attraction_id || item.attraction_id === 0) ? (
                <div></div>
              ) : (
                <>
                  <div><Label>{t('attractions.table.category')}</Label><p className="text-sm font-bold text-slate-800 italic uppercase">🎡 {t(`attractions.categories.${item.address?.toLowerCase().replace(/\s+/g, '_')}`) || item.address}</p></div>
                  <div><Label>{t('hotels.table.rating')}</Label><p className="font-bold text-blue-600 text-lg">📊 {item.rating ? `${item.rating} / 5` : 'N/A'}</p></div>
                </>
              )}
            </div>
          )}
          

          {type === 'restaurant' && (
            <div className="space-y-6">
              <div className="grid grid-cols-2 gap-6">
                <div><Label>{t('restaurants.table.cuisine')}</Label><p className="font-bold text-orange-700 text-base">🍳 {item.address}</p></div>
                <div><Label>{t('restaurants.table.rating')}</Label><p className="font-bold text-blue-600 text-lg">📊 {item.rating ? `${item.rating} / 5` : 'N/A'}</p></div>
              </div>
              <div><Label>{t('restaurants.table.price')}</Label><p className="font-bold text-green-600 text-xl tracking-widest">💰 {item.description?.replace('Price: ', '')}</p></div>
            </div>
          )}
          {translatedNote && (
            <div className="mt-4 pt-4 border-t border-gray-100">
              <Label>{t('trips.itinerary.fields.status')}</Label>
              <div className="p-5 bg-blue-50/50 rounded-[24px] text-sm text-blue-700 font-bold border border-blue-100 leading-relaxed shadow-sm italic">{translatedNote}</div>
            </div>
          )}
  
          {type === 'attraction' && hasAlternative && alternativeName.toLowerCase() !== 'place' && (
            <button onClick={handleSwap} disabled={isSwapping} className="w-full py-4 bg-gray-900 text-white rounded-2xl text-xs font-black uppercase tracking-[0.2em] hover:bg-black transition-all active:scale-95 shadow-xl">
                {isSwapping ? "..." : `🔄 ${t('trips.itinerary.swap_button')}`}
            </button>
          )}
          
          {item.website && type !== 'flight' && (
            <div className="pt-2"><a href={item.website.startsWith('http') ? item.website : `https://${item.website}`} target="_blank" rel="noopener noreferrer" className="flex items-center justify-center gap-2 w-full py-4 bg-gray-50 text-gray-900 rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-gray-200 transition-all border border-gray-100">🌐 {t('trips.itinerary.visit_website')} ↗</a></div>
          )}
        </div>
      </div>
    </div>
  );
}

const ActivityCard = memo(({ act, isFlexible, displayTime, onUpdate }) => {
  const { t } = useTranslation();
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleSwapActivity = async (newName) => {
    try { await swapActivity(act.activity_id, newName); if (onUpdate) onUpdate(); return true; } 
    catch (error) { alert(error.message); return false; }
  };

  const details = {
    name: act.entity_name || t('trips.itinerary.leisure'),
    address: act.entity_detail || "",
    description: act.entity_extra || "",
    notes: act.notes || "",
    website: act.entity_website?.Valid ? act.entity_website.String : (typeof act.entity_website === 'string' ? act.entity_website : ""),
    rating: act.entity_rating || null,
    stars: act.hotel_stars?.Valid ? act.hotel_stars.Int64 : (typeof act.hotel_stars === 'number' ? act.hotel_stars : 0),
    attraction_id: act.attraction_id?.Valid ? act.attraction_id.Int64 : (typeof act.attraction_id === 'number' ? act.attraction_id : 0),
    price: act.activity_type === 'flight' ? (act.entity_rating || act.price) : (act.entity_price || act.price || '---'),
    order_number: act.order_number
  };

  const formatTime = (ts) => (ts && ts.includes('T') ? ts.split('T')[1].substring(0, 5) : ts);

  return (
    <>
      <div className={`relative flex items-center gap-4 p-4 h-[100px] rounded-[24px] border-2 transition-all hover:shadow-lg ${ACTIVITY_META[act.activity_type]?.color || 'bg-gray-50'}`}>
        <div className="flex flex-col items-center justify-center min-w-[100px] border-r border-current/20 pr-4 text-center shrink-0">
          <span className="text-[12px] font-black text-gray-800 uppercase tracking-tighter">
            {isFlexible ? t('trips.itinerary.flexible_time') : displayTime.split(' - ').map(time => formatTime(time)).join('-')}
          </span>
        </div>
        <div className="flex-1 min-w-0 flex items-center justify-between gap-3 overflow-hidden text-left">
          <div className="flex items-center gap-3 min-w-0 overflow-hidden">
            <span className="text-2xl shrink-0">{ACTIVITY_META[act.activity_type]?.icon}</span>
            <h4 className="font-bold text-gray-900 line-clamp-2 text-sm leading-tight uppercase italic min-w-0">{details.name}</h4>
          </div>
          <button onClick={() => setIsModalOpen(true)} className="bg-white/90 shadow-sm px-4 py-2 rounded-xl text-[10px] border border-gray-200 hover:text-blue-600 font-black uppercase tracking-tight transition-all shrink-0 active:scale-95 whitespace-nowrap">
            {t('trips.itinerary.details_btn')}
          </button>
        </div>
      </div>
      {isModalOpen && <DetailsModal item={details} type={act.activity_type} onClose={() => setIsModalOpen(false)} onSwap={handleSwapActivity} />}
    </>
  );
});

export default function Itinerary() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { t, i18n } = useTranslation();
  const [days, setDays] = useState([]);
  const [activities, setActivities] = useState({});
  const [loading, setLoading] = useState(true);
  const [completing, setCompleting] = useState(false);
  const [isCompleted, setIsCompleted] = useState(false);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getTripItinerary(id);
      const list = Array.isArray(data) ? data : [];
      setDays(list);

      if (list.length > 0 && list[0].trip_status?.toLowerCase() === 'completed') {
        setIsCompleted(true);
      }

      const activitiesPromises = list.map(day => 
        getItineraryActivities(day.Itinerary_id)
          .then(acts => ({ id: day.Itinerary_id, acts: Array.isArray(acts) ? acts : [] }))
          .catch(() => ({ id: day.Itinerary_id, acts: [] }))
      );
      const results = await Promise.all(activitiesPromises);
      const newActivitiesMap = {};
      results.forEach(res => { newActivitiesMap[res.id] = res.acts; });
      setActivities(newActivitiesMap);
    } catch (e) { 
      console.error(e); 
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleComplete = async () => {
    const result = await Swal.fire({
      title: t('trips.itinerary.confirm_complete'),
      icon: 'question',
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: 'Այո',
      cancelButtonText: 'Ոչ',
      background: '#fff',
      borderRadius: '24px'
    });

    if (!result.isConfirmed) return;

    setCompleting(true);
    try {
      await completeTrip(id);
      
      Swal.fire({
        title: 'Հաջողված է',
        text: t('trips.itinerary.completed_success'),
        icon: 'success',
        timer: 2000,
        showConfirmButton: false,
        borderRadius: '24px'
      });

      setIsCompleted(true);
      loadData();
    } catch (err) {
      Swal.fire({
        title: 'Սխալ',
        text: "Failed to complete trip.",
        icon: 'error',
        confirmButtonColor: '#3085d6',
        borderRadius: '24px'
      });
    } finally {
      setCompleting(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto space-y-10 pb-20 px-4 text-left font-sans">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <PageHeader 
          icon="📋" 
          title={t('trips.itinerary.title')} 
          subtitle={t('trips.itinerary.subtitle')}
          action={
            <div className="flex gap-2">
              {!isCompleted ? (
                <button 
                  onClick={handleComplete}
                  disabled={completing}
                  className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors disabled:opacity-50"
                >
                  {completing ? 'Wait...' : `✅ ${t('trips.itinerary.complete_btn')}`}
                </button>
              ) : (
                <Link to="/reviews" className="bg-amber-500 hover:bg-amber-600 text-white px-3 py-2 rounded-lg text-sm font-medium">
                  ⭐ {t('trips.itinerary.write_review')}
                </Link>
              )}
              <Link to="/trips" className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 py-2 rounded-lg text-sm font-medium transition-colors">
                ← {t('trips.itinerary.back_btn')}
              </Link>
            </div>
          }
        />
      </div>

      {loading ? (
        <div className="text-center py-20 animate-pulse text-gray-400 font-black uppercase tracking-widest text-sm">Բեռնում...</div>
      ) : (
        <div className="space-y-12">
          {days.map((day, dIdx) => (
            <div key={day.Itinerary_id} className="space-y-6 text-left">
              <div className="border-b-4 border-slate-900 pb-3">
                <h3 className="text-3xl font-black text-gray-900 uppercase leading-none mb-2">{t('trips.itinerary.day')} {day.day_number}</h3>
                <p className="text-xs text-blue-600 font-black uppercase tracking-[0.3em] pl-1">{formatDynamicDate(day.date, i18n.language)}</p>
              </div>
              <div className="grid gap-4">
                {(activities[day.Itinerary_id] || []).length === 0 ? (
                  <div className="p-8 rounded-3xl border-2 border-dashed border-gray-100 text-gray-300 text-center font-bold italic tracking-wide">{t('trips.itinerary.free_day')}</div>
                ) : (
                  activities[day.Itinerary_id].map((act, aIdx) => (
                    <ActivityCard key={act.activity_id} act={act} isFlexible={dIdx === 0 || dIdx === days.length - 1} onUpdate={loadData}
                      displayTime={["10:00 - 12:00", "13:00 - 15:00", "16:00 - 18:00", "20:00 - 21:00"][aIdx] || "Activity"} 
                    />
                  ))
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

const getFlightDetails = (detailStr, extraStr) => {
  if (!detailStr) return { origin: 'N/A', destination: 'N/A', duration: 'N/A', source: 'N/A' };
  const [origin, destination] = detailStr.split(' ✈ ');
  const parts = (extraStr || "").split(' | ');
  const duration = parts[0]?.replace('Duration: ', '') || 'N/A';
  const source = parts[1]?.replace('Source: ', '') || 'N/A';
  return { origin, destination, duration, source };
};