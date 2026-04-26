import { createContext, useContext, useState, useEffect } from 'react';
import { useAuth } from './AuthContext';
import * as api from '../api/notifications'; 

const NotificationContext = createContext();

export function NotificationProvider({ children }) {
    const [notifications, setNotifications] = useState([]);
    const [unreadCount, setUnreadCount] = useState(0);
    const { isAuthenticated } = useAuth();

    const fetchNotifications = async () => {
        try {
            const data = await api.getNotifications();
            
            if (data && Array.isArray(data)) {
                setNotifications(data);
                setUnreadCount(data.filter(n => !n.is_read).length);
            } else {
                setNotifications([]);
                setUnreadCount(0);
            }
        } catch (err) {
            console.error("Error fetching notifications", err);
            setNotifications([]);
            setUnreadCount(0);
        }
    };

    useEffect(() => {
        const loadData = async () => {
            const token = localStorage.getItem('token');
            if (isAuthenticated && token) {
                await fetchNotifications();
            }
        };

        loadData();
    }, [isAuthenticated]);
    
    const markAsRead = async (id) => {
        try {
            await api.markNotificationAsRead(id);
            setNotifications(prev => 
                prev.map(n => n.id === id ? { ...n, is_read: true } : n)
            );
            setUnreadCount(prev => Math.max(0, prev - 1));
            return true;
        } catch (err) {
            return false;
        }
    };

    const deleteNotification = async (id) => {
        try {
            await api.deleteNotification(id); 
            setNotifications(prev => prev.filter(n => n.id !== id));
            setUnreadCount(prev => {
                const wasUnread = notifications.find(n => n.id === id && !n.is_read);
                return wasUnread ? Math.max(0, prev - 1) : prev;
            });
            return true;
        } catch (err) {
            console.error("Failed to delete", err);
            return false;
        }
    };

    return (
        <NotificationContext.Provider value={{ 
            notifications, 
            unreadCount, 
            setNotifications, 
            setUnreadCount, 
            markAsRead, 
            deleteNotification 
        }}>
            {children}
        </NotificationContext.Provider>
    );
}

export const useNotifications = () => useContext(NotificationContext);