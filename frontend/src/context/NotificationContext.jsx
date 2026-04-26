import { createContext, useContext, useState, useEffect } from 'react';
import axios from 'axios';
import { useAuth } from './AuthContext';

const NotificationContext = createContext();

export function NotificationProvider({ children }) {
    const [notifications, setNotifications] = useState([]);
    const [unreadCount, setUnreadCount] = useState(0);
    const { isAuthenticated } = useAuth();

    const fetchNotifications = async () => {
        try {
            const res = await axios.get('/api/notifications');
            setNotifications(res.data || []);
            const unread = res.data.filter(n => !n.is_read).length;
            setUnreadCount(unread);
        } catch (err) {
            console.error("Error fetching notifications", err);
        }
    };

    useEffect(() => {
        if (isAuthenticated) fetchNotifications();
    }, [isAuthenticated]);

    const markAsRead = async (id) => {
    try {
        // Նախ ուղարկում ենք API հարցումը
        await axios.post(`/api/notifications/${id}/read`);
        
        // Հաջողության դեպքում թարմացնում ենք state-ը
        setNotifications(prev => 
            prev.map(n => n.id === id ? { ...n, is_read: true } : n)
        );
        setUnreadCount(prev => Math.max(0, prev - 1));
        
        return true; // Վերադարձնում ենք true, որ իմանանք՝ ավարտվեց
    } catch (err) {
        console.error("Error marking as read", err);
        return false;
    }
};

    return (
        <NotificationContext.Provider value={{ notifications, unreadCount, setNotifications, setUnreadCount, markAsRead }}>
            {children}
        </NotificationContext.Provider>
    );
}

export const useNotifications = () => useContext(NotificationContext);