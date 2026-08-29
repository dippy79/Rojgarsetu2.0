import React, { useState, useEffect } from 'react';
import { Bell } from 'lucide-react';
import api from '../lib/api';

export default function NotificationBell() {
  const [notifications, setNotifications] = useState([]);
  const [open, setOpen] = useState(false);

  const fetchNotifications = async () => {
    try {
      const res = await api.get('/api/v1/notifications');
      setNotifications(Array.isArray(res.data) ? res.data : []);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchNotifications();
    const interval = setInterval(fetchNotifications, 30000);
    return () => clearInterval(interval);
  }, []);

  const unreadCount = notifications.filter(n => !n.is_read).length;

  return (
    <div className="relative">
      <button onClick={() => setOpen(!open)} className="relative p-2 text-slate-600 hover:bg-slate-100 rounded-full">
        <Bell className="w-5 h-5" />
        {unreadCount > 0 && (
          <span className="absolute top-1 right-1 w-4 h-4 bg-rose-500 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
            {unreadCount}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 mt-2 w-80 bg-white border border-slate-200 shadow-xl rounded-2xl p-4 z-50">
          <h4 className="font-bold text-slate-800 mb-3 text-sm">Notifications</h4>
          <div className="space-y-2 max-h-64 overflow-y-auto">
            {notifications.length === 0 ? (
              <p className="text-xs text-slate-400 py-2">No notifications yet.</p>
            ) : (
              notifications.slice(0, 5).map(n => (
                <div key={n.id} className={`p-2.5 rounded-xl text-xs border ${n.is_read ? 'bg-slate-50 border-slate-100' : 'bg-indigo-50/50 border-indigo-100'}`}>
                  <p className="font-semibold text-slate-800">{n.title}</p>
                  <p className="text-slate-600">{n.body}</p>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}
