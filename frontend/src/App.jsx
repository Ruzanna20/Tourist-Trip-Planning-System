import { useEffect } from 'react';
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

function WebSocketListener() {
  const { user } = useAuth();
  const { setNotifications, setUnreadCount } = useNotifications(); 
  const { t } = useTranslation();

  useEffect(() => {
    if (!user) return;

    const userID = user.user_id || user.id || user.ID || user.UserID;
    if (!userID) return;

    const socket = new WebSocket(`ws://localhost:8080/ws?userID=${userID}`);

    socket.onmessage = (event) => {
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
        console.error("WS Message Error:", err);
      }
    };

    return () => {
      if (socket.readyState === 1) socket.close();
    };
  }, [user, setNotifications, setUnreadCount, t]);

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
          <WebSocketListener />
          <AppRoutes />
        </NotificationProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}