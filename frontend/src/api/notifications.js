import client from './client';

export const getNotifications = () => 
    client.get('/api/notifications').then(r => r.data);

export const markNotificationAsRead = (id) => 
    client.post(`/api/notifications/${id}/read`).then(r => r.data);

export const getUnreadCount = () => 
    client.get('/api/notifications/unread-count').then(r => r.data);

export const deleteNotification = (id) => 
    client.delete(`/api/notifications/${id}`).then(r => r.data);