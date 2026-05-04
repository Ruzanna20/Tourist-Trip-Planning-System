import { useEffect,useRef } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { NotificationProvider, useNotifications } from './context/NotificationContext';
import toast, { Toaster } from 'react-hot-toast';
import { useTranslation } from 'react-i18next';

import ProtectedRoute from './components/ProtectedRoute';
import Layout from './components/Layout';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import MyTrips from './pages/MyTrips';
import CreateTrip from './pages/trips/CreateTrip';
import Itinerary from './pages/trips/Itinerary';
import Countries from './pages/Countries';
import Cities from './pages/Cities';
import Attractions from './pages/Attractions';
import Hotels from './pages/Hotels';
import Restaurants from './pages/Restaurants';
import Flights from './pages/Flights';
import Preferences from './pages/Preferences';
import Reviews from './pages/Reviews';
import Landing from './pages/Landing'; 

function SSEListener() {
  const { user } = useAuth();
  const { setNotifications, setUnreadCount } = useNotifications(); 
  const { t } = useTranslation();
  const eventSourceRef = useRef(null);

  useEffect(() => {
    const userID = user?.user_id || user?.id || user?.ID || user?.UserID;
    
    if (!userID) {
      console.log("SSE: No userID found, skipping connection");
      return;
    }

    if (eventSourceRef.current) {
      return;
    }

    console.log(`SSE: Attempting to connect for user ${userID}`);
    
    const eventSource = new EventSource(`http://localhost:8080/sse?userID=${userID}`);
    eventSourceRef.current = eventSource;

    eventSource.onopen = () => {
      console.log("SSE: Connection established");
    };

    eventSource.onmessage = (event) => {
      try {
        const rawData = JSON.parse(event.data);
        let finalData = rawData;
        
        if (typeof rawData.message === 'string' && rawData.message.startsWith('{')) {
            const parsedInner = JSON.parse(rawData.message);  
            finalData = { ...rawData, ...parsedInner };
        }
        
        if (finalData.type === 'TRIP_READY') {
          const translatedMessage = t('nav.trip_ready');

          toast.success(translatedMessage, {
            duration: 8000,
            position: 'top-right',
            icon: '✈️',
            style: {
              background: '#0f172a',
              color: '#fff',
              border: '2px solid #3b82f6',
              borderRadius: '12px',
              fontWeight: '500'
            },
          });

          const newNotification = {
            id: finalData.id, 
            message: translatedMessage,
            trip_id: finalData.trip_id,
            trip_title: finalData.trip_title,
            is_read: false,
            created_at: new Date().toLocaleString()
          };

          setNotifications(prev => [newNotification, ...prev]);
          setUnreadCount(prev => prev + 1);
        }
      } catch (err) {
        console.error("SSE Message Error:", err);
      }
    };

    eventSource.onerror = (error) => {
      console.error("SSE: Connection error or heartbeat lost", error);
    };

    return () => {
      if (eventSourceRef.current) {
        console.log("SSE: Closing connection...");
        eventSourceRef.current.close();
        eventSourceRef.current = null;
      }
    };
  }, [user, t, setNotifications, setUnreadCount]);

  return null;
}

function AppRoutes() {
  const { user } = useAuth();
  console.log("Current User:", user);


  return (
    <Routes>
      <Route path="/" element={user ? <Navigate to="/dashboard" replace /> : <Landing />} />
      <Route path="/login" element={user ? <Navigate to="/dashboard" replace /> : <Landing />} />
      <Route path="/register" element={user ? <Navigate to="/dashboard" replace /> : <Landing />} />
      <Route
        element={
          <ProtectedRoute>
            <Layout />
          </ProtectedRoute>
        }
      >
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/trips" element={<MyTrips />} />
        <Route path="/trips/create" element={<CreateTrip />} />
        <Route path="/trips/:id/itinerary" element={<Itinerary />} />
        <Route path="/trips/:id/options" element={<CreateTrip />} />
        <Route path="/countries" element={<Countries />} />
        <Route path="/cities" element={<Cities />} />
        <Route path="/attractions" element={<Attractions />} />
        <Route path="/hotels" element={<Hotels />} />
        <Route path="/restaurants" element={<Restaurants />} />
        <Route path="/flights" element={<Flights />} />
        <Route path="/preferences" element={<Preferences />} />
        <Route path="/reviews" element={<Reviews />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <NotificationProvider>
          <Toaster />
          <SSEListener />
          <AppRoutes />
        </NotificationProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}