import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AppContext';
import './FloatingNotification.css';

const FloatingNotification = () => {
  const { user, isAuthenticated } = useAuth();
  const [notifications, setNotifications] = useState([]);
  const [isOpen, setIsOpen] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);
  const [isLoading, setIsLoading] = useState(false);

  // Fetch notifications from API
  const fetchNotifications = useCallback(async () => {
    if (!isAuthenticated || !user) return;

    setIsLoading(true);
    try {
      const response = await fetch(`/api/v1/users/${user.id}/notifications?limit=10`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        const data = await response.json();
        if (data.status === 'success' && data.data) {
          setNotifications(data.data);
          const unread = data.data.filter(n => !n.read_at).length;
          setUnreadCount(unread);
        }
      }
    } catch (error) {
      console.error('Failed to fetch notifications:', error);
    } finally {
      setIsLoading(false);
    }
  }, [isAuthenticated, user]);

  // Mark notification as read
  const markAsRead = async (notificationId) => {
    try {
      const response = await fetch(`/api/v1/notifications/${notificationId}/read`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        setNotifications(prev => prev.map(n => 
          n.id === notificationId ? { ...n, read_at: new Date().toISOString() } : n
        ));
        setUnreadCount(prev => Math.max(0, prev - 1));
      }
    } catch (error) {
      console.error('Failed to mark notification as read:', error);
    }
  };

  // Mark notification as clicked (navigate to relevant page)
  const handleNotificationClick = (notification) => {
    markAsRead(notification.id);
    
    // Navigate based on notification payload
    if (notification.payload) {
      const payload = typeof notification.payload === 'string' 
        ? JSON.parse(notification.payload) 
        : notification.payload;
      
      if (payload.enrollment_id) {
        // Navigate to enrollment detail page
        window.location.href = `/enrollments/${payload.enrollment_id}`;
      } else if (payload.trade_id) {
        // Navigate to trade detail page
        window.location.href = `/trades/${payload.trade_id}`;
      }
    }
    
    setIsOpen(false);
  };

  // Format time ago
  const formatTimeAgo = (dateString) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString();
  };

  // Get notification icon based on type
  const getNotificationIcon = (type) => {
    switch (type) {
      case 'expiry_final':
        return '⚠️';
      case 'expiry_warning':
        return '⏰';
      case 'enrollment_reminder':
        return '📚';
      case 'course_update':
        return '🔄';
      default:
        return '🔔';
    }
  };

  // Get notification color based on type
  const getNotificationColor = (type) => {
    switch (type) {
      case 'expiry_final':
        return '#EF4444'; // red
      case 'expiry_warning':
        return '#F59E0B'; // amber
      case 'enrollment_reminder':
        return '#3B82F6'; // blue
      case 'course_update':
        return '#8B5CF6'; // purple
      default:
        return '#64748B'; // slate
    }
  };

  // Fetch on mount and when user changes
  useEffect(() => {
    if (isAuthenticated && user) {
      fetchNotifications();
      
      // Poll for new notifications every 30 seconds
      const interval = setInterval(fetchNotifications, 30000);
      return () => clearInterval(interval);
    }
  }, [fetchNotifications, isAuthenticated, user]);

  if (!isAuthenticated || !user) {
    return null;
  }

  return (
    <div className="floating-notification">
      {/* Notification Bell Button */}
      <button
        className="notification-bell"
        onClick={() => setIsOpen(!isOpen)}
        aria-label={`Notifications${unreadCount > 0 ? `, ${unreadCount} unread` : ''}`}
        aria-expanded={isOpen}
      >
        <span className="bell-icon" role="img" aria-hidden="true">🔔</span>
        {unreadCount > 0 && (
          <span className="notification-badge" aria-live="polite">
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>

      {/* Notification Dropdown */}
      {isOpen && (
        <div className="notification-dropdown" role="menu">
          <div className="notification-header">
            <h3>Notifications</h3>
            {unreadCount > 0 && (
              <button 
                className="mark-all-read"
                onClick={() => {
                  notifications.forEach(n => !n.read_at && markAsRead(n.id));
                }}
              >
                Mark all as read
              </button>
            )}
          </div>
          
          <div className="notification-list" role="list">
            {isLoading ? (
              <div className="notification-loading">Loading...</div>
            ) : notifications.length === 0 ? (
              <div className="notification-empty">
                <span className="empty-icon" role="img" aria-hidden="true">🔔</span>
                <p>No notifications yet</p>
                <span className="empty-hint">You'll see enrollment reminders here</span>
              </div>
            ) : (
              notifications.map((notification) => (
                <div
                  key={notification.id}
                  className={`notification-item ${!notification.read_at ? 'unread' : ''}`}
                  role="menuitem"
                  onClick={() => handleNotificationClick(notification)}
                >
                  <div 
                    className="notification-icon"
                    style={{ backgroundColor: getNotificationColor(notification.notification_type) }}
                    role="img"
                    aria-label={getNotificationIcon(notification.notification_type)}
                  >
                    {getNotificationIcon(notification.notification_type)}
                  </div>
                  <div className="notification-content">
                    <div className="notification-title-row">
                      <h4 className="notification-title">{notification.title}</h4>
                      <span className="notification-time">{formatTimeAgo(notification.sent_at)}</span>
                    </div>
                    <p className="notification-message">{notification.message}</p>
                    {notification.payload && (
                      <div className="notification-meta">
                        {(() => {
                          const payload = typeof notification.payload === 'string' 
                            ? JSON.parse(notification.payload) 
                            : notification.payload;
                          if (payload.days_until_expiry !== undefined) {
                            return (
                              <span className="expiry-badge" style={{ 
                                backgroundColor: payload.days_until_expiry <= 1 ? '#EF4444' : '#F59E0B' 
                              }}>
                                Expires in {payload.days_until_expiry} day{payload.days_until_expiry !== 1 ? 's' : ''}
                              </span>
                            );
                          }
                          return null;
                        })()}
                      </div>
                    )}
                  </div>
                  {!notification.read_at && (
                    <div className="unread-indicator" aria-hidden="true"></div>
                  )}
                </div>
              ))
            )}
          </div>

          <div className="notification-footer">
            <a href="/notifications" className="view-all-link">
              View all notifications
            </a>
          </div>
        </div>
      )}

      {/* Backdrop for closing on outside click */}
      {isOpen && (
        <div 
          className="notification-backdrop" 
          onClick={() => setIsOpen(false)}
          aria-hidden="true"
        />
      )}
    </div>
  );
};

export default FloatingNotification;